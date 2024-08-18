package crud

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"pkg.oars.vip/go-pkg/constant"
	"pkg.oars.vip/go-pkg/idaas"
	"pkg.oars.vip/go-pkg/perr"
	"pkg.oars.vip/go-pkg/server/base"
)

const (
	ListKind   = "list"
	CreateKind = "create"
	UpdateKind = "update"
	GetKind    = "get"
	DeleteKind = "delete"
)

type ResourceModel interface {
	GetResourceModel() interface{}
	ListResourceModel() interface{}
	ResourceName() string
}

type CommonModelList interface {
	ListORM(db *gorm.DB, g *gin.Context, resourceNames *idaas.ResourceNames) (*gorm.DB, interface{}, error)
}
type CommonModelGet interface {
	GetORM(db *gorm.DB, id string, g *gin.Context) (*gorm.DB, interface{}, error)
}
type CommonModelCreate interface {
	CreateORM(data any, db *gorm.DB, g *gin.Context) error
}
type CommonModelUpdate interface {
	UpdateORM(data any, db *gorm.DB, id string, g *gin.Context) error
}
type CommonModelDelete interface {
	DeleteORM(db *gorm.DB, id string, g *gin.Context) error
}
type CommonModelChangeCallback interface {
	ChangeCallback(data any)
}

type UpdateResourceModel interface {
	UpdateResourceModel() interface{}
}
type CreateResourceModel interface {
	CreateResourceModel() interface{}
}
type DeleteResourceModel interface {
	DeleteResourceModel() interface{}
}

type CommonModel struct {
	Id         string         `json:"id" gorm:"column:id;size:40"`
	Created    int64          `json:"created" gorm:"column:created;autoCreateTime:milli;comment:创建时间戳"`
	Updated    int64          `json:"updated" gorm:"column:updated;autoUpdateTime:milli;comment:更新时间戳"`
	CreatedBy  string         `json:"createdBy" gorm:"column:created_by;size:255;comment:创建用户ID"`
	SearchText string         `json:"searchText" gorm:"column:search_text;size:1024;comment:搜索字段"`
	DeletedAt  gorm.DeletedAt `json:"deleteAt" gorm:"index"`
	AppId      string         `json:"appId" gorm:"column:app_id;comment:应用ID"`
}
type CommonSearchModel struct {
	Id string `json:"id" gorm:"column:id;size:40"`
}

func (m *CommonModel) GenID() {
	m.Id = uuid.NewString()
}

func (m *CommonModel) SetID(id string) {
	m.Id = id
}
func (m *CommonModel) GetId() string {
	return m.Id
}
func (m *CommonModel) Bus() string {
	return "ID"
}
func (m *CommonModel) SetAppId(appId string) {
	m.AppId = appId
}

func (m *CommonModel) GenCreate(c any, g *gin.Context) error {
	return nil
}
func (m *CommonModel) SetCreatedBy(uid string) {
	m.CreatedBy = uid
}

type CommonModelInf interface {
	GenID()
	SetID(id string)
	GetId() string
	Bus() string
	SetCreatedBy(uid string)
	SetAppId(appId string)
}

type StoreTransaction interface {
	GetDB() *gorm.DB
}

type BaseInfoController struct {
	*base.BaseController
	Tx            StoreTransaction
	resources     map[string]ResourceModel
	services      map[string]any
	idaas         *idaas.Client
	resourceGroup string
	opt           *Option
}

type Option struct {
	ResourceGroup   string
	OperationLogSrv OperationLogService
}

func NewCrud(b *base.BaseController, tx StoreTransaction, idaas *idaas.Client, opt *Option) *BaseInfoController {
	return &BaseInfoController{
		BaseController: b,
		Tx:             tx,
		idaas:          idaas,
		resources:      make(map[string]ResourceModel),
		services:       make(map[string]any),
		resourceGroup:  opt.ResourceGroup,
		opt:            opt,
	}
}

func (c *BaseInfoController) RegResourceModel(m ResourceModel) {
	c.resources[m.ResourceName()] = m
}
func (c *BaseInfoController) RegService(m ResourceModel, s any) {
	c.services[m.ResourceName()] = s
}

func (c *BaseInfoController) GetBaseInfo(resource string, g *gin.Context, kind string) (interface{}, error) {
	md, ok := c.resources[resource]
	if !ok {
		return nil, errors.New("资源不存在")
	}
	m := md.GetResourceModel()
	if kind == ListKind {
		m = md.ListResourceModel()
	}
	if kind == UpdateKind {
		if um, ok := md.(UpdateResourceModel); ok {
			m = um.UpdateResourceModel()
		}
	}
	if kind == CreateKind {
		if um, ok := md.(CreateResourceModel); ok {
			m = um.CreateResourceModel()
		}
	}
	if kind == DeleteKind {
		if um, ok := md.(DeleteResourceModel); ok {
			m = um.DeleteResourceModel()
		}
	}
	if g != nil {
		err := g.ShouldBindJSON(m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}
func (c *BaseInfoController) GetService(resource string) any {
	s, ok := c.services[resource]
	if ok {
		return s
	}
	return nil
}

func (c *BaseInfoController) Create(g *gin.Context) {
	resource := g.Param("resource")
	m, err := c.GetBaseInfo(resource, g, CreateKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if c.resourceGroup != "" {
		ok, err := c.idaas.PermissionEnforce(idaas.EnforceParam{
			Group:        c.resourceGroup,
			Resource:     resource,
			ResourceName: "#",
			Action:       constant.CreateAction,
			UserId:       c.GetUid(g),
		})
		if err != nil {
			logrus.Error(err)
			c.Error(g, err)
			return
		}
		if !ok {
			c.Error(g, perr.ErrForbidden)
			return
		}
	}
	m.(CommonModelInf).GenID()
	m.(CommonModelInf).SetCreatedBy(c.GetUid(g))
	m.(CommonModelInf).SetAppId(c.GetAppId(g))
	BuildCreateGen(m)
	if l, ok := c.GetService(resource).(CommonModelCreate); ok {
		err = l.CreateORM(m, c.Tx.GetDB(), g)
	} else {
		err = c.Tx.GetDB().Create(m).Error
	}
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if l, ok := c.GetService(resource).(CommonModelChangeCallback); ok {
		l.ChangeCallback(m)
	}

	if c.opt.OperationLogSrv != nil {
		log := &OperationLog{
			Resource:     resource,
			ResourceName: m.(CommonModelInf).GetId(),
			Action:       CreateKind,
		}
		c.genOperationLog(log, nil, m)
	}

	c.OK(g, m)
}

func (c *BaseInfoController) Update(g *gin.Context) {
	resource := g.Param("resource")
	id := g.Param("id")
	m, err := c.GetBaseInfo(resource, g, UpdateKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if c.resourceGroup != "" {
		ok, err := c.idaas.PermissionEnforce(idaas.EnforceParam{
			Group:        c.resourceGroup,
			Resource:     resource,
			ResourceName: id,
			Action:       constant.UpdateAction,
			UserId:       c.GetUid(g),
		})
		if err != nil {
			logrus.Error(err)
			c.Error(g, err)
			return
		}
		if !ok {
			c.Error(g, perr.ErrForbidden)
			return
		}
	}
	BuildUpdateGen(m)
	if l, ok := c.GetService(resource).(CommonModelUpdate); ok {
		err = l.UpdateORM(m, c.Tx.GetDB(), id, g)
	} else {
		err = c.Tx.GetDB().Model(m).Where("id=?", id).Updates(m).Error
	}
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if l, ok := c.GetService(resource).(CommonModelChangeCallback); ok {
		l.ChangeCallback(m)
	}
	if c.opt.OperationLogSrv != nil {
		log := &OperationLog{
			Resource:     resource,
			ResourceName: m.(CommonModelInf).GetId(),
			Action:       UpdateKind,
		}
		c.genOperationLog(log, nil, m)
	}
	c.OK(g, m)
}

func (c *BaseInfoController) Delete(g *gin.Context) {
	id := g.Param("id")
	resource := g.Param("resource")
	m, err := c.GetBaseInfo(resource, nil, GetKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if c.resourceGroup != "" {
		ok, err := c.idaas.PermissionEnforce(idaas.EnforceParam{
			Group:        c.resourceGroup,
			Resource:     resource,
			ResourceName: id,
			Action:       constant.DeleteAction,
			UserId:       c.GetUid(g),
		})
		if err != nil {
			logrus.Error(err)
			c.Error(g, err)
			return
		}
		if !ok {
			c.Error(g, perr.ErrForbidden)
			return
		}
	}
	if l, ok := c.GetService(resource).(CommonModelDelete); ok {
		err = l.DeleteORM(c.Tx.GetDB(), id, g)
	} else {
		err = c.Tx.GetDB().Where("id=?", id).Delete(m).Error
	}
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if l, ok := c.GetService(resource).(CommonModelChangeCallback); ok {
		m.(CommonModelInf).SetID(id)
		l.ChangeCallback(m)
	}
	if c.opt.OperationLogSrv != nil {
		log := &OperationLog{
			Resource:     resource,
			ResourceName: m.(CommonModelInf).GetId(),
			Action:       DeleteKind,
		}
		c.genOperationLog(log, nil, m)
	}
	c.OK(g, nil)
}

func (c *BaseInfoController) Get(g *gin.Context) {
	id := g.Param("id")
	resource := g.Param("resource")
	res, err := c.GetBaseInfo(resource, nil, GetKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if c.resourceGroup != "" {
		ok, err := c.idaas.PermissionEnforce(idaas.EnforceParam{
			Group:        c.resourceGroup,
			Resource:     resource,
			ResourceName: id,
			Action:       GetKind,
			UserId:       c.GetUid(g),
		})
		if err != nil {
			logrus.Error(err)
			c.Error(g, err)
			return
		}
		if !ok {
			c.Error(g, perr.ErrForbidden)
			return
		}
	}

	db := c.Tx.GetDB()
	if l, ok := c.GetService(resource).(CommonModelGet); ok {
		db, res, err = l.GetORM(c.Tx.GetDB(), id, g)
		if err != nil {
			c.Error(g, err)
			return
		}
	} else {
		if v, ok := BuildGetORM(res, db); ok {
			db = v
			db = db.Where("m.id=?", id)
		} else {
			db = db.Where("id=?", id)
		}
	}
	err = db.First(res).Error
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if c.opt.OperationLogSrv != nil {
		log := &OperationLog{
			Resource:     resource,
			ResourceName: id,
			Action:       GetKind,
		}
		c.genOperationLog(log, nil, res)
	}
	c.OK(g, res)
}

func (c *BaseInfoController) List(g *gin.Context) {
	resource := g.Param("resource")
	q, err := c.GetBaseInfo(resource, nil, GetKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	resources := &idaas.ResourceNames{
		All: true,
	}
	if c.resourceGroup != "" {
		resources, err = c.idaas.PermissionResources(idaas.EnforceParam{
			Group:        c.resourceGroup,
			Resource:     resource,
			ResourceName: "*",
			Action:       constant.SelectAction,
			UserId:       c.GetUid(g),
		})
		if err != nil {
			logrus.Error(err)
			c.Error(g, err)
			return
		}
	}

	err = g.ShouldBindQuery(q)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}

	res, err := c.GetBaseInfo(resource, nil, ListKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}

	resType, err := c.GetBaseInfo(resource, nil, GetKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}

	pageNum := g.Query("pageNum")
	page, err := c.PageQuery(g)
	if err != nil {
		c.Error(g, err)
		return
	}
	db := c.Tx.GetDB()
	if l, ok := c.GetService(resource).(CommonModelList); ok {
		db, res, err = l.ListORM(c.Tx.GetDB(), g, resources)
		if err != nil {
			c.Error(g, err)
			return
		}
	} else {
		borm, ok := BuildListORM(res, db)
		if ok {
			db = borm
		} else {
			db = db.Model(resType)
		}

		if !resources.All {
			if len(resources.ResourceNames) == 0 {
				if pageNum != "" {
					c.PageOK(g, res, 0, page.PageNum, page.PageSize)
					return
				}
				c.OK(g, res)
				return
			}
			if ok {
				db = db.Where("m.id in (?)", resources.ResourceNames)
			} else {
				db = db.Where("id in (?)", resources.ResourceNames)
			}
		}
	}

	var total int64
	if pageNum != "" {
		err = db.Where(q).Count(&total).Error
		if err != nil {
			c.Error(g, err)
			return
		}

		start := (page.PageNum - 1) * page.PageSize
		db = db.Limit(page.PageSize).Offset(start)
	}

	err = db.Order("created desc").Find(&res, q).Error
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}

	if pageNum != "" {
		c.PageOK(g, res, int(total), page.PageNum, page.PageSize)
		return
	}

	c.OK(g, res)
}

func (c *BaseInfoController) Put(g *gin.Context) {
	resource := g.Param("resource")
	q, err := c.GetBaseInfo(resource, nil, GetKind)
	if err != nil {
		c.Error(g, err)
		return
	}
	err = g.ShouldBindQuery(q)
	if err != nil {
		c.Error(g, err)
		return
	}
	res, err := c.GetBaseInfo(resource, nil, GetKind)
	if err != nil {
		c.Error(g, err)
		return
	}
	err = c.Tx.GetDB().Order("created desc").First(&res, q).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			m, err := c.GetBaseInfo(resource, g, CreateKind)
			if err != nil {
				logrus.Error(err)
				c.Error(g, err)
				return
			}
			if c.resourceGroup != "" {
				ok, err := c.idaas.PermissionEnforce(idaas.EnforceParam{
					Group:        c.resourceGroup,
					Resource:     resource,
					ResourceName: "#",
					Action:       constant.CreateAction,
					UserId:       c.GetUid(g),
				})
				if err != nil {
					logrus.Error(err)
					c.Error(g, err)
					return
				}
				if !ok {
					c.Error(g, perr.ErrForbidden)
					return
				}
			}
			m.(CommonModelInf).GenID()
			err = c.Tx.GetDB().Create(m).Error
			if err != nil {
				logrus.Error(err)
				c.Error(g, err)
				return
			}
			if l, ok := c.GetService(resource).(CommonModelChangeCallback); ok {
				l.ChangeCallback(m)
			}
			c.OK(g, m)
			return
		}
	}

	m, err := c.GetBaseInfo(resource, g, UpdateKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if c.resourceGroup != "" {
		ok, err := c.idaas.PermissionEnforce(idaas.EnforceParam{
			Group:        c.resourceGroup,
			Resource:     resource,
			ResourceName: res.(CommonModelInf).GetId(),
			Action:       constant.UpdateAction,
			UserId:       c.GetUid(g),
		})
		if err != nil {
			logrus.Error(err)
			c.Error(g, err)
			return
		}
		if !ok {
			c.Error(g, perr.ErrForbidden)
			return
		}
	}

	err = c.Tx.GetDB().
		Model(m).Where("id=?", res.(CommonModelInf).GetId()).Updates(m).Error
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if l, ok := c.GetService(resource).(CommonModelChangeCallback); ok {
		l.ChangeCallback(m)
	}
	c.OK(g, m)
}
