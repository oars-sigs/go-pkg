package former

import (
	"encoding/json"

	"gorm.io/datatypes"
)

const (
	TaskWaitStatus    = 1
	TaskPassStatus    = 2
	TaskNOPassStatus  = 3
	TaskRecallStatus  = 4
	TaskDestroyStatus = 5
	TaskSkipStatus    = 6
	TaskRejectStatus  = 7
	TaskCheckStatus   = 8
	TaskConfirmStatus = 9

	ReleaseModelKind = "release"

	TaskEvent   = "task"
	NodeEvent   = "node"
	StatusEvent = "status"
)

var TaskStatusNames = map[int]string{
	0:                 "",
	TaskWaitStatus:    "审批中",
	TaskPassStatus:    "已通过",
	TaskNOPassStatus:  "已拒绝",
	TaskRecallStatus:  "撤回",
	TaskDestroyStatus: "作废",
	TaskSkipStatus:    "跳过",
	TaskRejectStatus:  "驳回",
	TaskCheckStatus:   "待确认",
	TaskConfirmStatus: "已确认",
}

type BusCategory struct {
	ID      string            `json:"id"  gorm:"column:id"`
	Name    string            `json:"name" gorm:"column:name"`
	Config  datatypes.JSONMap `json:"config" gorm:"column:config"`
	Created int64             `json:"created" gorm:"column:created;autoCreateTime:milli"`
	Updated int64             `json:"updated" gorm:"column:updated;autoUpdateTime:milli"`
}

func (f *BusCategory) TableName() string {
	return "bus_category"
}

type BusModel struct {
	ID       string         `json:"id" gorm:"column:id"`
	Kind     string         `json:"kind" gorm:"column:kind"`
	Class    string         `json:"class" gorm:"column:class"`
	Category string         `json:"category" gorm:"column:category"`
	Name     string         `json:"name" gorm:"column:name"`
	Mark     string         `json:"mark" gorm:"column:mark"`
	Data     datatypes.JSON `json:"data" grom:"column:data"`
	Version  int            `json:"version" grom:"column:version"`
	OrderNum int            `json:"orderNum" grom:"column:order_num"`
	Created  int64          `json:"created" gorm:"column:created;autoCreateTime:milli"`
	Updated  int64          `json:"updated" gorm:"column:updated;autoUpdateTime:milli"`
}

type BusModelData struct {
	Name   string                 `json:"name"`
	Flow   FlowConfig             `json:"flow"`
	Form   map[string]interface{} `json:"form"`
	Config map[string]interface{} `json:"config"`
}

func (f *BusModel) TableName() string {
	return "bus_model"
}

func (f *BusModel) Parse() (*BusModelData, error) {
	var data BusModelData
	err := json.Unmarshal(f.Data, &data)
	return &data, err
}

func (f *BusModel) SetData(d *BusModelData) {
	b, err := json.Marshal(d)
	if err != nil {
		return
	}
	f.Data = b
}

type BusData struct {
	ID          string                 `json:"id" gorm:"column:id"`
	ResourceId  string                 `json:"resourceId" gorm:"column:resource_id"`
	ModelId     string                 `json:"modelId" gorm:"column:model_id"`
	ReleaseId   string                 `json:"releaseId" gorm:"column:releaseId"`
	Title       string                 `json:"title" gorm:"column:title"`
	Creator     string                 `json:"creator" gorm:"column:creator"`
	Status      int                    `json:"status" gorm:"column:status"`
	Model       FlowAction             `json:"model" grom:"column:model"`
	Data        map[string]interface{} `json:"data" grom:"column:data"`
	Actions     []FlowAction           `json:"actions" grom:"column:actions"`
	Created     int64                  `json:"created" gorm:"column:created;autoCreateTime:milli"`
	Updated     int64                  `json:"updated" gorm:"column:updated;autoUpdateTime:milli"`
	IsCandidate bool                   `json:"isCandidate" gorm:"-"`
	IsCreator   bool                   `json:"isCreator" gorm:"-"`
	Task        *BusTask               `json:"task" gorm:"-"`
	Action      *FlowAction            `json:"action" gorm:"-"`
	IsSubmit    bool                   `json:"isSubmit" gorm:"-"`
	IsCanRecall bool                   `json:"isCanRecall" gorm:"-"`
	CurModel    *BusModel              `json:"curModel" gorm:"-"`
	ModelMark   string                 `json:"modelMark" gorm:"-"`
}

func (f *BusData) TableName() string {
	return "bus_data"
}

type BusTask struct {
	ID          string         `json:"id" gorm:"column:id"`
	BusID       string         `json:"busId" gorm:"column:bus_id"`
	ActionID    string         `json:"actionId" gorm:"column:action_id"`
	Name        string         `json:"name" gorm:"column:name"`
	UserId      string         `json:"userId" gorm:"column:user_id"`
	Type        int            `json:"type" gorm:"column:type"`
	ExamineMode int            `json:"examineMode" gorm:"column:examine_mode"`
	Status      int            `json:"status" gorm:"column:status"`
	Remark      string         `json:"remark" gorm:"column:remark"`
	SignFile    string         `json:"signFile" gorm:"column:sign_file"`
	Created     int64          `json:"created" gorm:"column:created;autoCreateTime:milli"`
	Updated     int64          `json:"updated" gorm:"column:updated;autoUpdateTime:milli"`
	Data        datatypes.JSON `json:"data" gorm:"column:data"`
	Proportion  float64        `json:"proportion" gorm:"column:proportion"`
	DestId      string         `json:"destId" gorm:"-"`
	NextAction  *FlowAction    `json:"nextAction" gorm:"-"`
	Mark        string         `json:"mark" gorm:"column:mark;size:255"`
}

func (f *BusTask) TableName() string {
	return "bus_task"
}

type ConfigItem struct {
	ID      string `json:"id" gorm:"column:id;primaryKey"`
	Name    string `json:"name" gorm:"column:name"`
	Group   string `json:"cfgGroup" gorm:"column:cfg_group"`
	Key     string `json:"key" gorm:"column:key;unique"`
	Value   string `json:"value" gorm:"column:value;size:1024"`
	Created int64  `json:"created" gorm:"column:created;autoCreateTime:milli"`
	Updated int64  `json:"updated" gorm:"column:updated;autoUpdateTime:milli"`
}

func (f *ConfigItem) TableName() string {
	return "bus_config"
}

const (
	LaunchFlowType    = 0
	CandidateFlowType = 1
	CopyFlowType      = 2
)

type FlowAction struct {
	ID           string       `json:"id"`
	Type         int          `json:"type"` //0 发起人 1审批 2抄送 3条件 4路由
	ExamineMode  int          `json:"examineMode"`
	Name         string       `json:"name"`
	Users        []ActionUser `json:"users"`
	DefaultUsers []ActionUser `json:"defaultUsers"`
	SelectUsers  []ActionUser `json:"selectUsers"`
	NodeConfig   *NodeConfig  `json:"nodeConfig"`
	Key          string       `json:"key"`
	KeyWeight    int          `json:"keyWeight"`
	Mark         string       `json:"mark"`
}

type ActionUser struct {
	Name       string   `json:"name"`
	ID         string   `json:"id"`
	SearchText string   `json:"searchText"`
	Task       *BusTask `json:"task"`
}
type FlowData struct {
	Flow    FlowConfig             `json:"flow"`
	Form    map[string]interface{} `json:"form"`
	Data    map[string]interface{} `json:"data"`
	Creator string                 `json:"creator"`
	Actions []FlowAction           `json:"actions"`
	Config  map[string]interface{} `json:"config"`
	Name    string                 `json:"name"`
}

type FlowConfig struct {
	TableID          string       `json:"tableId"`
	WorkFlowDef      *WorkFlowDef `json:"workFlowDef"`
	DirectorMaxLevel int          `json:"directorMaxLevel"`
	FlowPermission   []any        `json:"flowPermission"`
	NodeConfig       *NodeConfig  `json:"nodeConfig"`
}
type WorkFlowDef struct {
	Name string `json:"name"`
}

type ConditionList struct {
	ColumnID          string         `json:"columnId"`
	Type              int            `json:"type"`
	ConditionEn       string         `json:"conditionEn"`
	ConditionCn       string         `json:"conditionCn"`
	OptType           string         `json:"optType"`
	Zdy1              string         `json:"zdy1"`
	Zdy2              string         `json:"zdy2"`
	Opt1              string         `json:"opt1"`
	Opt2              string         `json:"opt2"`
	ColumnDbname      string         `json:"columnDbname"`
	ColumnType        string         `json:"columnType"`
	ShowType          string         `json:"showType"`
	ShowName          string         `json:"showName"`
	FixedDownBoxValue string         `json:"fixedDownBoxValue"`
	NodeUserList      []NodeUserList `json:"nodeUserList"`
}
type NodeUserList struct {
	TargetID   string `json:"targetId"`
	Type       int    `json:"type"`
	Name       string `json:"name"`
	SearchText string `json:"searchText"`
}
type NodeConfig struct {
	NodeName                string          `json:"nodeName"`
	Type                    int             `json:"type"`
	PriorityLevel           int             `json:"priorityLevel"`
	Settype                 int             `json:"settype"`
	SelectMode              int             `json:"selectMode"`
	SelectRange             int             `json:"selectRange"`
	DeptRange               int             `json:"deptRange"`
	DirectorLevel           int             `json:"directorLevel"`
	ExamineMode             int             `json:"examineMode"`
	NoHanderAction          int             `json:"noHanderAction"`
	ExamineEndDirectorLevel int             `json:"examineEndDirectorLevel"`
	CcSelfSelectFlag        int             `json:"ccSelfSelectFlag"`
	ConditionList           []ConditionList `json:"conditionList"`
	NodeUserList            []NodeUserList  `json:"nodeUserList"`
	DeptRangeList           []NodeUserList  `json:"deptRangeList"`
	ChildNode               *NodeConfig     `json:"childNode"`
	ConditionNodes          []NodeConfig    `json:"conditionNodes"`
	Error                   bool            `json:"error"`
	FormConfig              []interface{}   `json:"formConfig"`
	MsgTpl                  string          `json:"msgTpl"`
	Proportion              float64         `json:"proportion"`
	Mark                    string          `json:"mark"`
	Preconditions           string          `json:"preconditions"`
}

type Hook struct {
	Event   string       `json:"event"`
	Data    *BusData     `json:"data"`
	Task    *BusTask     `json:"task"`
	Actions []FlowAction `json:"actions"`
	Action  *FlowAction  `json:"action"`
	Model   *BusModel    `json:"model"`
	Status  int          `json:"status"`
}
