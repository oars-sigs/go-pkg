package crud

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"pkg.oars.vip/go-pkg/constant"
	"pkg.oars.vip/go-pkg/former"
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
	ImportKind = "import"

	MsgCtx = "msgCtx"
)

type ResourceModel interface {
	GetResourceModel() interface{}
	ListResourceModel() interface{}
	ResourceName() string
}

type CommonModelList interface {
	ListORM(db *gorm.DB, g *gin.Context, resourceNames *idaas.ResourceNames) (*gorm.DB, interface{}, error)
}

type CommonModelListGen interface {
	ListGen(any) any
}

type CommonModelGet interface {
	GetORM(db *gorm.DB, id string, g *gin.Context) (*gorm.DB, interface{}, error)
}

type CommonModelGetGen interface {
	GetGen(any) any
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
	ChangeCallback(data any, action string)
}

type UpdateResourceModel interface {
	UpdateResourceModel() interface{}
}
type UpdateResourceModelWithCtx interface {
	UpdateResourceModel(g *gin.Context) interface{}
}
type CreateResourceModel interface {
	CreateResourceModel() interface{}
}
type CreateResourceModelWithCtx interface {
	CreateResourceModel(g *gin.Context) interface{}
}
type DeleteResourceModel interface {
	DeleteResourceModel() interface{}
}

type ImportResourceModel interface {
	ImportResourceModel() interface{}
}

type GetResourceName interface {
	GetResourceName() (string, string)
}

type DisablePermission interface {
	DisablePermission() bool
}

type FlowHookSvc interface {
	FlowHook(h *former.Hook) error
}

type CommonModel struct {
	Id         string         `json:"id" gorm:"column:id;type:varchar(40);size:40"`
	Created    int64          `json:"created" gorm:"column:created;autoCreateTime:milli;comment:创建时间戳"`
	Updated    int64          `json:"updated" gorm:"column:updated;autoUpdateTime:milli;comment:更新时间戳"`
	CreatedBy  string         `json:"createdBy" gorm:"column:created_by;type:varchar(255);size:255;comment:创建用户ID"`
	SearchText string         `json:"searchText" gorm:"column:search_text;type:varchar(1024);size:1024;comment:搜索字段"`
	DeletedAt  gorm.DeletedAt `json:"deleteAt" gorm:"index"`
	AppId      string         `json:"appId" gorm:"column:app_id;type:varchar(255);size:255;comment:应用ID"`
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

type SimpleModel struct {
	Id        string `json:"id" gorm:"column:id;type:varchar(40);size:40"`
	Created   int64  `json:"created" gorm:"column:created;autoCreateTime:milli;comment:创建时间戳"`
	Updated   int64  `json:"updated" gorm:"column:updated;autoUpdateTime:milli;comment:更新时间戳"`
	CreatedBy string `json:"createdBy" gorm:"column:created_by;type:varchar(255);size:255;comment:创建用户ID"`
}

func (m *SimpleModel) GenID() {
	m.Id = uuid.NewString()
}

func (m *SimpleModel) SetID(id string) {
	m.Id = id
}
func (m *SimpleModel) GetId() string {
	return m.Id
}
func (m *SimpleModel) Bus() string {
	return "ID"
}
func (m *SimpleModel) SetAppId(appId string) {
}

func (m *SimpleModel) GenCreate(c any, g *gin.Context) error {
	return nil
}
func (m *SimpleModel) SetCreatedBy(uid string) {
	m.CreatedBy = uid
}

type StoreTransaction interface {
	GetDB() *gorm.DB
}

type BaseInfoController struct {
	*base.BaseController
	Tx            StoreTransaction
	resources     map[string]resourceModelReg
	services      map[string]any
	formerHooks   map[string]FlowHookSvc
	idaas         *idaas.Client
	resourceGroup string
	opt           *Option
}

type resourceModelReg struct {
	m      ResourceModel
	option *ResourceModelOption
}

type ResourceModelOption struct {
	Actions []string
}

type Option struct {
	ResourceGroup   string
	OperationLogSrv OperationLogService
	Mgr             any
	Former          *former.Client
}

type Page struct {
	base.Page
	List interface{} `json:"rows"`
}

func GenPage(result interface{}, total, pageNum, pageSize int) Page {
	var res Page
	res.List = result
	res.Total = total
	res.PageNum = pageNum
	res.PageSize = pageSize
	return res
}

func NewCrud(b *base.BaseController, tx StoreTransaction, idaas *idaas.Client, opt *Option) *BaseInfoController {
	if opt == nil {
		opt = &Option{}
	}
	if opt.Former != nil {
		tx.GetDB().Set("gorm:table_options", "COMMENT='资源流程关联表'").AutoMigrate(&ResourceFlowInfo{})
	}
	return &BaseInfoController{
		BaseController: b,
		Tx:             tx,
		idaas:          idaas,
		resources:      make(map[string]resourceModelReg),
		services:       make(map[string]any),
		formerHooks:    make(map[string]FlowHookSvc),
		resourceGroup:  opt.ResourceGroup,
		opt:            opt,
	}
}

func (c *BaseInfoController) RegResourceModel(m ResourceModel, opt ...*ResourceModelOption) {
	r := resourceModelReg{
		m: m,
	}
	if len(opt) > 0 {
		r.option = opt[0]
	}
	c.resources[m.ResourceName()] = r
}
func (c *BaseInfoController) RegService(m ResourceModel, s any) {
	c.services[m.ResourceName()] = s
}

func (c *BaseInfoController) RegFormerHook(mark string, s FlowHookSvc) {
	c.formerHooks[mark] = s
}

func (c *BaseInfoController) GetBaseInfo(resource string, g *gin.Context, kind string) (interface{}, error) {
	md, ok := c.resources[resource]
	if !ok {
		return nil, errors.New("资源不存在")
	}
	if md.option != nil && len(md.option.Actions) > 0 {
		exist := false
		for _, act := range md.option.Actions {
			if act == kind {
				exist = true
			}
		}
		if !exist {
			return nil, perr.ErrForbidden
		}
	}
	m := md.m.GetResourceModel()
	if kind == ListKind {
		m = md.m.ListResourceModel()
	}
	if kind == UpdateKind {
		if um, ok := md.m.(UpdateResourceModel); ok {
			m = um.UpdateResourceModel()
		}
		if um, ok := md.m.(UpdateResourceModelWithCtx); ok {
			m = um.UpdateResourceModel(g)
		}
	}
	if kind == CreateKind {
		if um, ok := md.m.(CreateResourceModel); ok {
			m = um.CreateResourceModel()
		}
		if um, ok := md.m.(CreateResourceModelWithCtx); ok {
			m = um.CreateResourceModel(g)
		}
	}
	if kind == DeleteKind {
		if um, ok := md.m.(DeleteResourceModel); ok {
			m = um.DeleteResourceModel()
		}
	}
	if kind == ImportKind {
		if um, ok := md.m.(ImportResourceModel); ok {
			m = um.ImportResourceModel()
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

func (c *BaseInfoController) EnPermission(s any) bool {
	if c.resourceGroup == "" {
		return false
	}
	dis, ok := s.(DisablePermission)
	if ok {
		return !dis.DisablePermission()
	}
	return true
}

func (c *BaseInfoController) Create(g *gin.Context) {
	resource := g.Param("resource")
	m, err := c.GetBaseInfo(resource, g, CreateKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if c.EnPermission(m) {
		if v, ok := m.(GetResourceName); ok {
			prResource, prResourceName := v.GetResourceName()
			if prResource != resource {
				ok, err := c.idaas.GetClient(g).PermissionEnforce(idaas.EnforceParam{
					Group:        c.resourceGroup,
					Resource:     prResource,
					ResourceName: prResourceName,
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
		} else {
			ok, err := c.idaas.GetClient(g).PermissionEnforce(idaas.EnforceParam{
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
		l.ChangeCallback(m, CreateKind)
	}

	if c.opt.OperationLogSrv != nil {
		log := &OperationLog{
			Resource:     resource,
			ResourceName: m.(CommonModelInf).GetId(),
			Action:       CreateKind,
			CommonModel: CommonModel{
				CreatedBy: c.GetUid(g),
				AppId:     c.GetAppId(g),
			},
		}
		log.GenID()
		c.genOperationLog(log, nil, m)
	}

	c.OK(g, m, g.GetString(MsgCtx))
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
	if c.EnPermission(m) {
		if v, ok := m.(GetResourceName); ok {
			prResource, prResourceName := v.GetResourceName()
			if prResource != resource {
				ok, err := c.idaas.GetClient(g).PermissionEnforce(idaas.EnforceParam{
					Group:        c.resourceGroup,
					Resource:     prResource,
					ResourceName: prResourceName,
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
		} else {
			ok, err := c.idaas.GetClient(g).PermissionEnforce(idaas.EnforceParam{
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
	}
	oldRes, err := c.GetBaseInfo(resource, nil, UpdateKind)
	if err != nil {
		logrus.Error(err)
	}
	err = c.Tx.GetDB().Model(oldRes).Where("id=?", id).First(oldRes).Error
	if err != nil {
		logrus.Error(err)
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
		l.ChangeCallback(m, UpdateKind)
	}
	if c.opt.OperationLogSrv != nil {
		log := &OperationLog{
			Resource:     resource,
			ResourceName: m.(CommonModelInf).GetId(),
			Action:       UpdateKind,
			CommonModel: CommonModel{
				CreatedBy: c.GetUid(g),
				AppId:     c.GetAppId(g),
			},
		}
		log.GenID()
		c.genOperationLog(log, oldRes, m)
	}
	c.OK(g, m, g.GetString(MsgCtx))
}

func (c *BaseInfoController) getRes(resource, id string) (any, error) {
	q, err := c.GetBaseInfo(resource, nil, GetKind)
	if err != nil {
		return nil, err
	}
	err = c.Tx.GetDB().Where("id=?", id).First(&q).Error
	if err != nil {
		return nil, err
	}
	return q, nil
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
	if c.EnPermission(m) {
		q, _ := c.getRes(resource, id)
		if v, ok := q.(GetResourceName); ok {
			prResource, prResourceName := v.GetResourceName()
			if prResource != resource {
				ok, err := c.idaas.GetClient(g).PermissionEnforce(idaas.EnforceParam{
					Group:        c.resourceGroup,
					Resource:     prResource,
					ResourceName: prResourceName,
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
		} else {
			ok, err := c.idaas.GetClient(g).PermissionEnforce(idaas.EnforceParam{
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
		l.ChangeCallback(m, DeleteKind)
	}
	if c.opt.OperationLogSrv != nil {
		log := &OperationLog{
			Resource:     resource,
			ResourceName: m.(CommonModelInf).GetId(),
			Action:       DeleteKind,
			CommonModel: CommonModel{
				CreatedBy: c.GetUid(g),
				AppId:     c.GetAppId(g),
			},
		}
		log.GenID()
		c.genOperationLog(log, nil, m)
	}
	c.OK(g, nil, g.GetString(MsgCtx))
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
	vres := res
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
	if db != nil {
		err = db.First(res).Error
		if err != nil {
			logrus.Error(err)
			c.Error(g, err)
			return
		}
	}
	if c.EnPermission(vres) {
		if v, ok := res.(GetResourceName); ok {
			prResource, prResourceName := v.GetResourceName()
			if prResource != resource {
				ok, err := c.idaas.GetClient(g).PermissionEnforce(idaas.EnforceParam{
					Group:        c.resourceGroup,
					Resource:     prResource,
					ResourceName: prResourceName,
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
		} else {
			ok, err := c.idaas.GetClient(g).PermissionEnforce(idaas.EnforceParam{
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
	}
	if c.opt.OperationLogSrv != nil {
		log := &OperationLog{
			Resource:     resource,
			ResourceName: id,
			Action:       GetKind,
			CommonModel: CommonModel{
				CreatedBy: c.GetUid(g),
				AppId:     c.GetAppId(g),
			},
		}
		log.GenID()
		c.genOperationLog(log, nil, res)
	}
	if l, ok := c.GetService(resource).(CommonModelGetGen); ok {
		res = l.GetGen(res)
	}
	c.OK(g, res, g.GetString(MsgCtx))
}

func (c *BaseInfoController) List(g *gin.Context) {
	resource := g.Param("resource")
	q, err := c.GetBaseInfo(resource, nil, GetKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	err = g.ShouldBindQuery(q)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	resources := &idaas.ResourceNames{
		All: true,
	}

	if c.EnPermission(q) {
		if v, ok := q.(GetResourceName); ok {
			prResource, prResourceName := v.GetResourceName()
			if prResource != resource {
				ok, err := c.idaas.GetClient(g).PermissionEnforce(idaas.EnforceParam{
					Group:        c.resourceGroup,
					Resource:     prResource,
					ResourceName: prResourceName,
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
		} else {
			resources, err = c.idaas.GetClient(g).PermissionResources(idaas.EnforceParam{
				Group:        c.resourceGroup,
				Resource:     resource,
				ResourceName: "*",
				Action:       "*",
				UserId:       c.GetUid(g),
			})
			if err != nil {
				logrus.Error(err)
				c.Error(g, err)
				return
			}
			if resources == nil {
				resources = &idaas.ResourceNames{
					All: false,
				}
			}
		}
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
		if db == nil {
			c.OK(g, res)
			return
		}
	} else {
		borm, ok := BuildListORM(res, db, &BuildORMOption{
			Search:       g.Query("search"),
			SearchText:   g.Query("searchText"),
			SeniorSearch: g.Query("seniorSearch"),
			SortField:    g.Query("sortField"),
			Order:        g.Query("order"),
		})
		if ok {
			db = borm
		} else {
			db = borm.Model(resType)
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
		dw := "deleted_at is NULL"
		if len(db.Statement.Table) != 0 {
			dw = db.Statement.Table + "." + dw
		}
		err = db.Where(dw).Where(q).Count(&total).Error
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
	if l, ok := c.GetService(resource).(CommonModelListGen); ok {
		res = l.ListGen(res)
	}
	if pageNum != "" {
		c.PageOK(g, res, int(total), page.PageNum, page.PageSize)
		return
	}

	c.OK(g, res, g.GetString(MsgCtx))
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
			if c.EnPermission(res) {
				ok, err := c.idaas.GetClient(g).PermissionEnforce(idaas.EnforceParam{
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
				l.ChangeCallback(m, CreateKind)
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
	if c.EnPermission(m) {
		ok, err := c.idaas.GetClient(g).PermissionEnforce(idaas.EnforceParam{
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
		l.ChangeCallback(m, UpdateKind)
	}
	c.OK(g, m, g.GetString(MsgCtx))
}

func (c *BaseInfoController) Export(g *gin.Context) {
	resource := g.Param("resource")

	svc := c.GetService(resource)
	sheets := &ExportListOption{
		Sheets: make(map[string]CommonModelExport),
	}
	if listexport, ok := svc.(CommonModelExportList); ok {
		sheets = listexport.ListExportORM(g)
		fmt.Println(sheets)
	} else {
		export, ok := svc.(CommonModelExport)
		if !ok {
			c.Error(g, perr.New("不支持导出"))
			return
		}
		sheets.Sheets["Sheet1"] = export
	}
	xlsx := excelize.NewFile()
	defer xlsx.Close()
	isfirst := true
	for sheetName, export := range sheets.Sheets {
		xlsx.NewSheet(sheetName)
		if isfirst && sheetName != "Sheet1" {
			xlsx.DeleteSheet("Sheet1")
			isfirst = false
		}
		db, data, opt, err := export.ExportORM(c.Tx.GetDB(), c, g)
		if err != nil {
			logrus.Error(err)
			c.Error(g, err)
			return
		}
		if sheets.Name == "" {
			sheets.Name = opt.Name
		}
		if db != nil {
			err = db.Find(&data).Error
			if err != nil {
				logrus.Error(err)
				c.Error(g, err)
				return
			}
		}
		sliceValue := reflect.ValueOf(data)
		var headers [][]HeaderItem
		var allHeaders []HeaderItem
		style, _ := xlsx.NewStyle(&excelize.Style{
			Border: []excelize.Border{
				{
					Style: 1,
					Type:  "left",
					Color: "0000000",
				},
				{
					Style: 1,
					Type:  "right",
					Color: "0000000",
				},
				{
					Style: 1,
					Type:  "top",
					Color: "0000000",
				},
				{
					Style: 1,
					Type:  "bottom",
					Color: "0000000",
				},
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})
		type posCell struct {
			v string
			p int
		}
		mergeCell := make(map[int][]posCell)
		for i := 0; i < sliceValue.Len(); i++ {
			if len(opt.Headers) == 0 {
				opt.Headers = make([]HeaderItem, 0)
				GetExportHeaders(sliceValue.Index(i).Type(), &opt.Headers)
				if len(opt.Extends) > 0 {
					opt.Headers = append(opt.Headers, opt.Extends...)
				}
				sort.Slice(opt.Headers, func(i, j int) bool {
					return opt.Headers[i].Order < opt.Headers[j].Order
				})

				if len(opt.Ignores) > 0 {
					hs := make([]HeaderItem, 0)
					for _, h := range opt.Headers {
						ignore := false
						for _, s := range opt.Ignores {
							if s == h.Filed {
								ignore = true
								break
							}
						}
						if !ignore {
							hs = append(hs, h)
						}
					}
					opt.Headers = hs
				}
			}
			if len(headers) == 0 {
				headers = GenHeaders(opt.Headers)
			}
			headerRow := len(headers)
			if len(allHeaders) == 0 {
				var vallHeaders []HeaderItem
				for _, hs := range headers {
					for _, h := range hs {
						if h.W == 0 {
							vallHeaders = append(vallHeaders, h)
						}
					}
				}
				allHeaders = vallHeaders
				indexspan := make(map[int]bool)
				allHeaders = make([]HeaderItem, len(vallHeaders))
				for _, hs := range headers {
					w := 0
					for _, h := range hs {
						w = GetNextIndex(indexspan, w)
						if h.W == 0 {
							allHeaders[w] = h
							indexspan[w] = true
							w++
						} else {
							w += h.W
						}
					}
				}
			}

			for j, h := range allHeaders {
				var data any
				if v, ok := sliceValue.Index(i).Interface().(map[string]interface{}); ok {
					data = v[h.Filed]

					if data != nil && reflect.TypeOf(data).Kind() == reflect.Ptr {
						vdata := reflect.ValueOf(data)
						if !vdata.IsNil() {
							data = vdata.Elem().Interface()
						} else {
							data = ""
						}
					}
				} else {
					v := sliceValue.Index(i).FieldByName(h.Filed)
					if !v.IsValid() {
						continue
					}
					if v.Kind() == reflect.Ptr {
						if !v.IsNil() {
							data = v.Elem().Interface()
						} else {
							data = ""
						}
					} else {
						data = v.Interface()
					}
				}

				if fn, ok := h.Fns["format"]; ok {
					switch fn {
					case "toDate":
						data = time.UnixMilli(data.(int64)).Format("2006-01-02")
					case "toDatetime":
						data = time.UnixMilli(data.(int64)).Format("2006-01-02 15:04:05")
					}
				}
				if h.AutoMerge {
					if len(mergeCell[j]) == 0 {
						mergeCell[j] = append(mergeCell[j], posCell{
							p: i + headerRow,
							v: fmt.Sprint(data),
						})
					} else if mergeCell[j][len(mergeCell[j])-1].v != fmt.Sprint(data) {
						mergeCell[j] = append(mergeCell[j], posCell{
							p: i + headerRow,
							v: fmt.Sprint(data),
						})
					}
				}
				xlsx.SetCellValue(sheetName, GetEcelAxis(i+headerRow+1, j+1), data)
				xlsx.SetCellStyle(sheetName, GetEcelAxis(i+headerRow+1, j+1), GetEcelAxis(i+headerRow+1, j+1), style)
			}
		}

		//表头
		indexspan := make(map[int]bool)
		for i, row := range headers {
			j := 0
			cindexspan := make(map[int]bool)
			for _, col := range row {
				j = GetNextIndex(indexspan, j)
				if col.D == 0 && col.L < len(headers) {
					xlsx.MergeCell(sheetName, GetEcelAxis(i+1, j+1),
						GetEcelAxis(i+1+(len(headers)-col.L), j+1))
					cindexspan[j] = true
					xlsx.SetCellStyle(sheetName, GetEcelAxis(i+1, j+1),
						GetEcelAxis(i+1+(len(headers)-col.L), j+1), style)
				}
				xlsx.SetCellValue(sheetName, GetEcelAxis(i+1, j+1), col.Title)
				if col.W > 0 {
					xlsx.MergeCell(sheetName, GetEcelAxis(i+1, j+1), GetEcelAxis(i+1, j+col.W))
					xlsx.SetCellStyle(sheetName, GetEcelAxis(i+1, j+1), GetEcelAxis(i+1, j+col.W), style)
					j += col.W
				} else {
					xlsx.SetCellStyle(sheetName, GetEcelAxis(i+1, j+1), GetEcelAxis(i+1, j+1), style)
					j++
				}

				if col.Width > 0 {
					xlsx.SetColWidth(sheetName, GetColumnIndex(i+1),
						GetColumnIndex(i+1), float64(col.Width))
				}
				for j, m := range mergeCell {
					for i, p := range m {
						if i == len(m)-1 {
							if sliceValue.Len()+len(headers)-p.p > 2 {
								xlsx.MergeCell(sheetName, GetEcelAxis(p.p+1, j+1), GetEcelAxis(sliceValue.Len()+len(headers), j+1))
							}
							continue
						}
						if m[i+1].p-p.p > 1 {
							xlsx.MergeCell(sheetName, GetEcelAxis(p.p+1, j+1), GetEcelAxis(m[i+1].p, j+1))
						}
					}
				}
			}
			for k, v := range cindexspan {
				indexspan[k] = v
			}
		}
	}

	g.Writer.Header().Set("Content-Type", "application/octet-stream")
	disposition := fmt.Sprintf("attachment; filename=\"%s_%s.xlsx\"", sheets.Name, time.Now().Format("20060102"))
	g.Writer.Header().Set("Content-Disposition", disposition)
	_ = xlsx.Write(g.Writer)

}

func (c *BaseInfoController) Import(g *gin.Context) {
	resource := g.Param("resource")
	m, err := c.GetBaseInfo(resource, nil, ImportKind)
	if err != nil {
		c.Error(g, err)
		return
	}
	sheet := g.PostForm("sheet")
	offsetRow, _ := strconv.Atoi(g.PostForm("offsetRow"))

	hs := []ImportColumn{}
	GetImportHeaders(reflect.TypeOf(m).Elem().Elem(), &hs)
	err = ParseImport(g, &ImportOption{Column: hs, Sheet: sheet, OffsetRow: offsetRow}, &m)
	if err != nil {
		c.Error(g, err)
		return
	}

	d := reflect.ValueOf(m)
	for i := 0; i < d.Len(); i++ {
		if v, ok := d.Index(i).Interface().(CommonModelInf); ok {
			v.GenID()
		}
	}
	db := c.Tx.GetDB()
	err = db.CreateInBatches(m, 100).Error
	if err != nil {
		c.Error(g, err)
		return
	}
	c.OK(g, nil)
}
