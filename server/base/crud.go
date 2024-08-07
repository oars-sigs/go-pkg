package base

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	ListKind   = 1
	CreateKind = 2
	UpdateKind = 3
	GetKind    = 4
)

type ResourceModel interface {
	GetResourceModel() interface{}
	ListResourceModel() interface{}
}

type CommonModelList interface {
	ListORM(db *gorm.DB, c any, g *gin.Context) (*gorm.DB, interface{}, error)
}
type CommonModelGet interface {
	GetORM(db *gorm.DB, id string, c any, g *gin.Context) (*gorm.DB, interface{}, error)
}
type CommonModelCreate interface {
	CreateORM(db *gorm.DB, c any, g *gin.Context) error
}
type CommonModelUpdate interface {
	UpdateORM(db *gorm.DB, id string, c any, g *gin.Context) error
}
type CommonModelDelete interface {
	DeleteORM(db *gorm.DB, id string, c any, g *gin.Context) error
}

type CommonModel struct {
	Id        string         `json:"id" gorm:"column:id;size:40"`
	Created   int64          `json:"created" gorm:"column:created;autoCreateTime:milli;comment:创建时间戳"`
	Updated   int64          `json:"updated" gorm:"column:updated;autoUpdateTime:milli;comment:更新时间戳"`
	CreatedBy string         `json:"createdBy" gorm:"column:created_by;size:255;comment:创建用户ID"`
	DeletedAt gorm.DeletedAt `json:"deleteAt" gorm:"index"`
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
	Cb(mgr any)
	Bus() string
	GenCreate(c any, g *gin.Context) error
	SetCreatedBy(uid string)
}

type CallbackFn interface {
	Cb(mgr any)
}

type StoreTransaction interface {
	GetDB() *gorm.DB
}

type BaseInfoController struct {
	*BaseController
	Tx        StoreTransaction
	Mgr       any
	Resources map[string]ResourceModel
}

func NewCrud(b *BaseController, tx StoreTransaction, mgr any) *BaseInfoController {
	return &BaseInfoController{
		BaseController: b,
		Tx:             tx,
		Mgr:            mgr,
		Resources:      make(map[string]ResourceModel),
	}
}

func (c *BaseInfoController) RegResourceModel(name string, m ResourceModel) {
	c.Resources[name] = m
}

func (c *BaseInfoController) GetBaseInfo(resource string, g *gin.Context, kind int) (interface{}, error) {
	md, ok := c.Resources[resource]
	if !ok {
		return nil, errors.New("资源不存在")
	}
	m := md.GetResourceModel()
	if kind == ListKind {
		m = md.ListResourceModel()
	}
	if g != nil {
		err := g.ShouldBindJSON(m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (c *BaseInfoController) Create(g *gin.Context) {
	resource := g.Param("resource")
	m, err := c.GetBaseInfo(resource, g, CreateKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	m.(CommonModelInf).GenID()
	m.(CommonModelInf).GenCreate(c.Mgr, g)
	m.(CommonModelInf).SetCreatedBy(c.GetUid(g))
	if l, ok := m.(CommonModelCreate); ok {
		err = l.CreateORM(c.Tx.GetDB(), c.Mgr, g)
	} else {
		err = c.Tx.GetDB().Create(m).Error
	}
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	m.(CommonModelInf).Cb(c.Mgr)
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
	if l, ok := m.(CommonModelUpdate); ok {
		err = l.UpdateORM(c.Tx.GetDB(), id, c.Mgr, g)
	} else {
		err = c.Tx.GetDB().Model(m).Where("id=?", id).Updates(m).Error
	}
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	m.(CommonModelInf).Cb(c.Mgr)
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
	if l, ok := m.(CommonModelDelete); ok {
		err = l.DeleteORM(c.Tx.GetDB(), id, c.Mgr, g)
	} else {
		err = c.Tx.GetDB().Where("id=?", id).Delete(m).Error
	}
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	m.(CommonModelInf).Cb(c.Mgr)
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
	db := c.Tx.GetDB()
	if l, ok := res.(CommonModelGet); ok {
		db, res, err = l.GetORM(c.Tx.GetDB(), id, c.Mgr, g)
		if err != nil {
			c.Error(g, err)
			return
		}
	} else {
		db = db.Where("id=?", id)
	}
	err = db.First(res).Error
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
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
	if l, ok := resType.(CommonModelList); ok {
		db, res, err = l.ListORM(c.Tx.GetDB(), c.Mgr, g)
		if err != nil {
			c.Error(g, err)
			return
		}
	} else {
		db = db.Model(resType)
	}

	var total int64
	if pageNum != "" {
		err = db.Count(&total).Error
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
			m.(CommonModelInf).GenID()
			err = c.Tx.GetDB().Create(m).Error
			if err != nil {
				logrus.Error(err)
				c.Error(g, err)
				return
			}
			m.(CommonModelInf).Cb(c.Mgr)
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
	err = c.Tx.GetDB().
		Model(m).Where("id=?", res.(CommonModelInf).GetId()).Updates(m).Error
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	m.(CommonModelInf).Cb(c.Mgr)
	c.OK(g, m)
}
