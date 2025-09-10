package crud

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"pkg.oars.vip/go-pkg/former"
	"pkg.oars.vip/go-pkg/perr"
)

// ResourceFlowInfo 资源相关流程表
type ResourceFlowInfo struct {
	CommonModel

	FromId       string         `json:"fromId" gorm:"column:from_id;size:100;comment:资源ID"`
	FromType     string         `json:"fromType" gorm:"column:from_type;size:100;comment:资源类型"`
	Remark       string         `json:"remark" gorm:"column:remark;size:200;comment:流程模型Remark"`
	FlowId       string         `json:"flowId" gorm:"column:flow_id;size:100;comment:关联的流程ID"`
	Model        string         `json:"model" gorm:"column:model;size:100;comment:关联的流程Model"`
	CompleteData datatypes.JSON `json:"completeData" gorm:"column:complete_data;type:json;comment:流程模型数据"`
}

// TableName ResourceFlowInfo's table name
func (*ResourceFlowInfo) TableName() string {
	return "base_resource_flowinfo"
}

type FormerFlowData struct {
	Data         map[string]any      `json:"data"`
	ResourceData json.RawMessage     `json:"resourceData"`
	CompleteData datatypes.JSON      `json:"completeData"`
	ActionUsers  []former.ActionUser `json:"actionUsers"`
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
	isCreate := id == ""
	actKind := UpdateKind
	if isCreate {
		actKind = CreateKind
	}
	m, err := c.GetBaseInfo(resource, nil, actKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if isCreate {
		m.(CommonModelInf).GenID()
		id = m.(CommonModelInf).GetId()
	}
	if flowData.ResourceData != nil {
		err = json.Unmarshal(flowData.ResourceData, m)
		if err != nil {
			logrus.Error(err)
			c.Error(g, err)
			return
		}
	}

	flowData.Data["id"] = id
	uid := c.GetUid(g)
	db := c.Tx.GetDB()
	var res any
	err = db.Transaction(func(tx *gorm.DB) error {
		res, err = CreateFormer(tx, c.opt.Former, id, uid, resource, modelMark, flowData.Data, flowData.CompleteData)
		if err != nil {
			return err
		}
		if flowData.ResourceData != nil {
			if isCreate {
				BuildCreateGen(m)
				if l, ok := c.GetService(resource).(CommonModelCreate); ok {
					err = l.CreateORM(m, tx, g)
				} else {
					err = tx.Create(m).Error
				}
				if err != nil {
					return err
				}
			} else {
				err = tx.Model(m).Where("id = ?", id).Updates(m).Error
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}

	c.OK(g, res)
}

func (c *BaseInfoController)Approve(g *gin.Context){
	resource := g.Param("resource")
	m, err := c.GetBaseInfo(resource, nil, GetKind)
	if err != nil {
		logrus.Error(err)
		c.Error(g, err)
		return
	}
	if svc,ok:=c.GetService(resource).(FlowHookBeforeApprove);ok{
		svc.FlowHookBeforeApprove(c.GetUid(g),m.(*former.BusTask))
	}

	uid:=c.GetUid(g)
	c.opt.Former.Approve(uid, &former.BusTask{
		
	})
}

func CreateFormer(db *gorm.DB, formercli *former.Client, id, uid, resource, modelMark string, flowData map[string]any, completeData datatypes.JSON) (*ResourceFlowInfo, error) {
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
		FromId:       id,
		FromType:     resource,
		FlowId:       res.ID,
		Model:        modelMark,
		Remark:       m.Name,
		CompleteData: completeData,
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
	db := c.Tx.GetDB()
	var flowInfo ResourceFlowInfo
	err := db.Find(&flowInfo, &ResourceFlowInfo{FlowId: h.Data.ID}).Error
	if err != nil {
		return err
	}
	if svc, ok := c.GetService(flowInfo.FromType).(FlowHookWithResourceSvc); ok {
		return svc.FlowHook(h, &flowInfo)
	}

	if h.Event != "status" {
		return nil
	}

	if flowInfo.CompleteData != nil {
		var d map[string]json.RawMessage
		err = json.Unmarshal(flowInfo.CompleteData, &d)
		if err != nil {
			return err
		}
		if v, ok := d[strconv.Itoa(h.Status)]; ok {
			m, err := c.GetBaseInfo(flowInfo.FromType, nil, UpdateKind)
			if err != nil {
				return err
			}
			err = json.Unmarshal(v, m)
			if err != nil {
				return err
			}
			err = db.Model(m).Where("id = ?", flowInfo.FromId).Updates(m).Error
			if err != nil {
				return err
			}
		}

	}
	return nil
}
