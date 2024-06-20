package base

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	GenResourceModel(m interface{}, kind int, g *gin.Context) error
}

type CommonModel struct {
	Id      string `json:"id" gorm:"column:id;size:40"`
	Created int64  `json:"created" gorm:"column:created;autoCreateTime:milli;comment:创建时间戳"`
	Updated int64  `json:"updated" gorm:"column:updated;autoUpdateTime:milli;comment:更新时间戳"`
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

type CommonModelInf interface {
	GenID()
	SetID(id string)
	GetId() string
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
		if kind == CreateKind || kind == UpdateKind {
			err = md.GenResourceModel(m, kind, g)
			if err != nil {
				return nil, err
			}
		}
	}
	return m, nil
}

func (c *BaseInfoController) Create(g *gin.Context) {
	resource := g.Param("resource")
	m, err := c.GetBaseInfo(resource, g, CreateKind)
	if err != nil {
		c.Error(g, err)
		return
	}
	m.(CommonModelInf).GenID()
	err = c.Tx.GetDB().Create(m).Error
	if err != nil {
		c.Error(g, err)
		return
	}
	if fn, ok := m.(CallbackFn); ok {
		fn.Cb(c.Mgr)
	}
	c.OK(g, m)
}

func (c *BaseInfoController) Update(g *gin.Context) {
	resource := g.Param("resource")
	id := g.Param("id")
	m, err := c.GetBaseInfo(resource, g, UpdateKind)
	if err != nil {
		c.Error(g, err)
		return
	}
	err = c.Tx.GetDB().Debug().Model(m).Where("id=?", id).Updates(m).Error
	if err != nil {
		c.Error(g, err)
		return
	}
	if fn, ok := m.(CallbackFn); ok {
		fn.Cb(c.Mgr)
	}
	c.OK(g, m)
}

func (c *BaseInfoController) Delete(g *gin.Context) {
	id := g.Param("id")
	resource := g.Param("resource")
	m, err := c.GetBaseInfo(resource, nil, GetKind)
	if err != nil {
		c.Error(g, err)
		return
	}
	err = c.Tx.GetDB().Where("id=?", id).Delete(m).Error
	if err != nil {
		c.Error(g, err)
		return
	}

	if fn, ok := m.(CallbackFn); ok {
		fn.Cb(c.Mgr)
	}
	c.OK(g, nil)
}

func (c *BaseInfoController) Get(g *gin.Context) {
	id := g.Param("id")
	resource := g.Param("resource")
	res, err := c.GetBaseInfo(resource, nil, GetKind)
	if err != nil {
		c.Error(g, err)
		return
	}
	err = c.Tx.GetDB().Where("id=?", id).First(res).Error
	if err != nil {
		c.Error(g, err)
		return
	}
	c.OK(g, res)
}

func (c *BaseInfoController) List(g *gin.Context) {
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
	res, err := c.GetBaseInfo(resource, nil, ListKind)
	if err != nil {
		c.Error(g, err)
		return
	}
	err = c.Tx.GetDB().Debug().Order("created desc").Find(&res, q).Error
	if err != nil {
		c.Error(g, err)
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
				c.Error(g, err)
				return
			}
			m.(CommonModelInf).GenID()
			err = c.Tx.GetDB().Create(m).Error
			if err != nil {
				c.Error(g, err)
				return
			}
			if fn, ok := m.(CallbackFn); ok {
				fn.Cb(c.Mgr)
			}
			c.OK(g, m)
			return
		}
	}

	m, err := c.GetBaseInfo(resource, g, UpdateKind)
	if err != nil {
		c.Error(g, err)
		return
	}
	err = c.Tx.GetDB().
		Model(m).Where("id=?", res.(CommonModelInf).GetId()).Updates(m).Error
	if err != nil {
		c.Error(g, err)
		return
	}
	if fn, ok := m.(CallbackFn); ok {
		fn.Cb(c.Mgr)
	}
	c.OK(g, m)
}
