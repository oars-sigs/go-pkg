package crud

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"pkg.oars.vip/go-pkg/former"
	"pkg.oars.vip/go-pkg/perr"
)

// ResourceFlowInfo 资源相关流程表
type ResourceFlowInfo struct {
	CommonModel

	FromId   string `json:"fromId" gorm:"column:from_id;size:100;comment:资源ID"`
	FromType string `json:"fromType" gorm:"column:from_type;size:100;comment:资源类型"`
	Remark   string `json:"remark" gorm:"column:remark;size:200;comment:流程模型Remark"`
	FlowId   string `json:"flowId" gorm:"column:flow_id;size:100;comment:关联的流程ID"`
	Model    string `json:"model" gorm:"column:model;size:100;comment:关联的流程Model"`
}

// TableName ResourceFlowInfo's table name
func (*ResourceFlowInfo) TableName() string {
	return "base_resource_flowinfo"
}

type FormerFlowData struct {
	Data        map[string]any      `json:"data"`
	Update      any                 `json:"update"`
	ActionUsers []former.ActionUser `json:"actionUsers"`
}

func (c *BaseInfoController) CreateFormer(g *gin.Context) {
	modelMark := g.Param("mark")
	id := g.Param("id")
	resource := g.Param("resource")
	if modelMark == "" {
		c.Error(g, perr.New("参数错误"))
		return
	}
	if c.opt.Former == nil {
		c.Error(g, perr.New("未配置former"))
		return
	}
	var flowData FormerFlowData
	err := g.ShouldBindJSON(&flowData)
	if err != nil {
		c.Error(g, perr.New("参数错误"))
		return
	}
	if flowData.Data == nil {
		flowData.Data = map[string]any{}
	}
	m, err := c.GetBaseInfo(resource, g, GetKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if flowData.Update != nil {
		err = mapstructure.Decode(flowData.Update, m)
		if err != nil {
			c.Error(g, err)
			return
		}
	}

	flowData.Data["id"] = id
	uid := c.GetUid(g)
	db := c.Tx.GetDB()
	var res any
	db.Transaction(func(tx *gorm.DB) error {
		res, err = CreateFormer(tx, c.opt.Former, id, uid, resource, modelMark, flowData.Data)
		if err != nil {
			return err
		}
		if flowData.Update != nil {
			err = tx.Model(m).Where("id = ?", id).Updates(m).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	c.OK(g, res)
}

func CreateFormer(db *gorm.DB, formercli *former.Client, id, uid, resource, modelMark string, flowData map[string]any) (*ResourceFlowInfo, error) {
	m, err := formercli.GetModel(uid, &former.BusData{
		ModelMark:  modelMark,
		Data:       flowData,
		ResourceId: id,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	res, err := formercli.Create(uid, &former.BusData{
		ModelId:    m.Flow.TableID,
		Data:       flowData,
		Actions:    m.Actions,
		ResourceId: id,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	flowInfo := &ResourceFlowInfo{
		FromId:   id,
		FromType: resource,
		FlowId:   res.ID,
		Model:    modelMark,
		Remark:   m.Name,
	}
	flowInfo.CreatedBy = uid
	flowInfo.GenID()
	err = db.Create(flowInfo).Error
	if err != nil {
		return nil, err
	}
	return flowInfo, nil
}

func (c *BaseInfoController) FlowHook(h *former.Hook) error {
	if hook, ok := c.formerHooks[h.Model.Mark]; ok {
		return hook.FlowHook(h)
	}
	return nil
}
