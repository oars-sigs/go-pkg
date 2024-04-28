package idaas

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"gorm.io/datatypes"
)

type UserInfo struct {
	//用户ID
	UserId string `json:"userId" gorm:"column:user_id"`
	//用户池ID
	PoolId string `json:"poolId" gorm:"column:user_pool_id"`
	//用户名
	Username string `json:"username" gorm:"column:user_username"`
	//用户密码
	NewPassword string `json:"newPassword" gorm:"-"`
	//swagger:ignore
	Password string `json:"-" gorm:"column:user_password"`
	//用户别名
	NickName string `json:"nickName" gorm:"column:user_nick_name"`
	//性别
	Gender string `json:"gender" gorm:"column:user_gender"`
	//部门ID
	DeptId string `json:"deptId" gorm:"column:user_dept_id"`
	//部门IDS
	DeptIds []string `json:"deptIds" gorm:"-"`
	//岗位ID
	PostId string `json:"postId" gorm:"column:user_post_id"`
	//岗位名称
	PostName string `json:"postName" gorm:"column:user_post_name"`
	//角色ID、岗位类型
	RoleId string `json:"roleId" gorm:"column:user_role_id"`
	//用户组ID
	GroupId string `json:"groupId" gorm:"column:user_group_id"`
	//手机号
	Mobile string `json:"mobile" gorm:"column:user_mobile"`
	//电话
	Telephone string `json:"telephone" gorm:"column:user_telephone"`
	//邮箱
	Email string `json:"email" gorm:"column:user_email"`
	//头像
	Avatar string `json:"avatar" gorm:"column:user_avatar"`
	//状态
	State int `json:"state" gorm:"column:user_state"`
	//状态值
	Status string `json:"status" gorm:"column:user_status"`
	//调动
	Transfer string `json:"transfer" gorm:"column:user_transfer"`
	//细化状态
	DetailStatus string `json:"detailStatus" gorm:"column:user_detail_status"`
	//排序
	OrderNum int `json:"orderNum" gorm:"column:user_order_num"`
	//创建时间
	Created int64 `json:"created" gorm:"column:user_created"`
	//更新时间
	Updated int64 `json:"updated" gorm:"column:user_updated"`
	//源用户ID
	SrcId string `json:"srcId" gorm:"column:user_src_id"`
	//extend
	Extend datatypes.JSONMap `json:"extend" gorm:"column:user_extend" swaggertype:"object"`
	//入职时间
	Hiredate JSONTime `json:"hiredate" gorm:"column:user_hiredate" swaggertype:"string"`
	//生日
	Birthday JSONTime `json:"birthday" gorm:"column:user_birthday" swaggertype:"string"`
	//卡号1
	CardNum1 string `json:"cardNum1" gorm:"column:user_cardnum1" swaggertype:"string"`
	//卡号2
	CardNum2 string `json:"cardNum2" gorm:"column:user_cardnum2" swaggertype:"string"`
	//卡号3
	CardNum3 string `json:"cardNum3" gorm:"column:user_cardnum3" swaggertype:"string"`
	//独生子女
	OnlyChild *int `json:"onlyChild" gorm:"column:user_only_child"`
	//已婚
	Married *int `json:"married" gorm:"column:user_married"`
	//类型
	Kind int `json:"kind" gorm:"-"`
	//部门
	Dept *Department `json:"dept" gorm:"-"`
	//用户组
	Group *Group `json:"group" gorm:"-"`
	//用户组数组
	Groups []Group `json:"groups" gorm:"-"`
	//应用ID
	AppID string `json:"appId" gorm:"-"`
	//真实用户ID
	RealId string `json:"realId" gorm:"-"`
	//是否实名
	Verified bool `json:"verified" gorm:"-"`
	//搜索字符
	SearchText string `json:"searchText" gorm:"column:user_search_text"`
	Type       string `json:"type" gorm:"-"`
}

type Group struct {
	ID          string      `json:"id" gorm:"column:group_id"`
	PoolId      string      `json:"poolId" gorm:"column:group_pool_id"`
	DeptId      string      `json:"deptId" gorm:"column:group_dept_id"`
	ParentId    string      `json:"parentId" gorm:"column:group_parent_id"`
	Name        string      `json:"name" gorm:"column:group_name"`
	Desc        string      `json:"desc" gorm:"column:group_desc"`
	Created     int64       `json:"created" gorm:"column:group_created;autoCreateTime:milli"`
	Updated     int64       `json:"updated" gorm:"column:group_updated;autoUpdateTime:milli"`
	OrderNum    int         `json:"orderNum" gorm:"column:group_order_num"`
	SrcId       string      `json:"srcId" gorm:"column:group_src_id"`
	Path        string      `json:"path" gorm:"-"`
	Paths       []string    `json:"paths"  gorm:"-"`
	NamePath    string      `json:"namePath"  gorm:"-"`
	Children    []*Group    `json:"children"  gorm:"-"`
	Dept        *Department `json:"dept"  gorm:"-"`
	HasChildren bool        `json:"hasChildren" gorm:"-"`
	//前端界面使用
	IsAfter    bool        `json:"isAfter" gorm:"-"`
	UserCount  int         `json:"userCount" gorm:"-"`
	Users      []*UserInfo `json:"users"  gorm:"-"`
	CacheUsers []*UserInfo `json:"-"  gorm:"-"`
	Type       string      `json:"type"  gorm:"-"`
}

type Department struct {
	DeptId   string        `json:"deptId"`
	PoolId   string        `json:"poolId"`
	DeptName string        `json:"deptName"`
	ParentId string        `json:"parentId"`
	OrderNum int           `json:"orderNum"`
	Created  int64         `json:"created"`
	Updated  int64         `json:"updated"`
	Children []*Department `json:"children"`
	Path     string        `json:"path"`
	NamePath string        `json:"namePath"`
	Kind     string        `json:"kind"`
}

// ThirdUser 第三方用户
type ThirdUser struct {
	TUserId     string `json:"tuser_id" gorm:"column:tuser_id"`
	PoolId      string `json:"poolId" gorm:"column:pool_id"`
	UserId      string `json:"userId" gorm:"column:user_id"`
	OpenId      string `json:"openid" gorm:"column:openid"`
	LoginType   int    `json:"login_type" gorm:"column:login_type"`
	Visitor     int    `json:"visitor" gorm:"column:visitor"`
	UnionId     string `json:"unionid" gorm:"column:unionid"`
	SessionKey  string `json:"sessionKey" gorm:"-"`
	AccessToken string `json:"accessToken" gorm:"-"`
	Username    string `json:"username" gorm:"-"`
	Mobile      string `json:"mobile" gorm:"-"`
	Email       string `json:"email" gorm:"-"`
	AutoReg     bool   `json:"autoReg" gorm:"-"`
	LoginState  int    `json:"login_state" gorm:"-"`
}

type TokenInfo struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	NickName string `json:"nickName"`
	DeptId   string `json:"deptId"`
	DeptPath string `json:"deptPath"`
	AppId    string `json:"appId"`
	RealId   string `json:"realId"`
	Kind     int    `json:"kind"`
	Verified bool   `json:"verified"`
}

const (
	UserKindUser         = "u"
	UserKindDept         = "d"
	UserKindPosition     = "p"
	UserKindUserRole     = "ur"
	UserKindResourceRole = "rr"

	ActionKindAction            = "a"
	ActionKindUserRole          = "ur"
	ActionKindResourceRole      = "rr"
	ActionKindPermissionInherit = "pi"

	RoleKindUserRole     = "ur"
	RoleKindResourceRole = "rr"
)

type PermissionResource struct {
	Id       int                 `json:"id" gorm:"column:id"`
	Group    string              `json:"group" gorm:"column:group"`       //分组，如idaas
	Name     string              `json:"name" gorm:"column:name"`         //资源名，如文件
	Resource string              `json:"resource" gorm:"column:resource"` //资源，如file
	Creator  string              `json:"creator" gorm:"column:creator"`
	Created  int64               `json:"created" gorm:"column:created;autoUpdateTime:milli"`
	Updated  int64               `json:"updated" gorm:"column:updated;autoUpdateTime:milli"`
	Actions  []PermissionActions `json:"actions" gorm:"-"`
}

type PermissionActions struct {
	Id           int64  `json:"id" gorm:"column:id"`
	Group        string `json:"group" gorm:"column:group"`       //分组，如idaas
	Resource     string `json:"resource" gorm:"column:resource"` //资源，如file
	ResourceName string `json:"resourceName" gorm:"column:resource_name"`
	Name         string `json:"name" gorm:"column:name"`     //操作名，如创建
	Action       string `json:"action" gorm:"column:action"` //操作，如create
	Creator      string `json:"creator" gorm:"column:creator"`
	Created      int64  `json:"created" gorm:"column:created;autoUpdateTime:milli"`
	Updated      int64  `json:"updated" gorm:"column:updated;autoUpdateTime:milli"`
}

type PermissionRoles struct {
	Id           int64  `json:"id" gorm:"column:id"`
	Group        string `json:"group" gorm:"column:group"`                //分组，如idaas
	Kind         string `json:"kind" gorm:"column:kind"`                  //类型，ur: 用户角色，rr: 资源角色
	Name         string `json:"name" gorm:"column:name"`                  //角色名，如管理员
	Role         string `json:"role" gorm:"column:role"`                  //角色，如admin
	Resource     string `json:"resource" gorm:"column:resource"`          //关联资源，*代表所有
	ResourceName string `json:"resourceName" gorm:"column:resource_name"` //关联资源名，*代表所有
	Creator      string `json:"creator" gorm:"column:creator"`
	Created      int64  `json:"created" gorm:"column:created;autoUpdateTime:milli"`
	Updated      int64  `json:"updated" gorm:"column:updated;autoUpdateTime:milli"`
}

type PermissionRolebindings struct {
	Id                int64       `json:"id" gorm:"column:id"`
	Pool              string      `json:"pool" gorm:"column:pool"`
	Group             string      `json:"group" gorm:"column:group"`                           //分组，如idaas
	UserKind          string      `json:"userKind" gorm:"column:user_kind"`                    //授权对象类型，u：用户，d: 部门，p: 岗位， ur: 用户角色， rr: 资源角色
	User              string      `json:"user" gorm:"column:user"`                             //授权对象，用户id、部门id、岗位id、用户角色、资源角色
	MatchResource     string      `json:"matchResource" gorm:"column:match_resource"`          //匹配资源
	MatchResourceName string      `json:"matchResourceName" gorm:"column:match_resource_name"` //匹配资源名
	Resource          string      `json:"resource" gorm:"column:resource"`                     //资源
	ResourceName      string      `json:"resourceName" gorm:"column:resource_name"`            //资源名
	ActionKind        string      `json:"actionKind" gorm:"column:action_kind"`                //关联操作类型，a： 操作， ur: 用户角色， rr: 资源角色
	Action            string      `json:"action" gorm:"column:action"`                         //操作
	Exclude           int         `json:"exclude" gorm:"column:exclude"`                       //黑白名单，默认： 白名单，1： 黑名单
	Creator           string      `json:"creator" gorm:"column:creator"`
	Created           int64       `json:"created" gorm:"column:created;autoUpdateTime:milli"`
	Updated           int64       `json:"updated" gorm:"column:updated;autoUpdateTime:milli"`
	UserDetail        interface{} `json:"userDetail" gorm:"-"` //detail
}

type ResourceNames struct {
	All           bool                `json:"all"`
	ResourceNames []string            `json:"resourceNames"`
	Permissions   map[string][]string `json:"permissions"`
}

type PermissionData struct {
	Rules        []PermissionRolebindings `json:"rules"`
	ExcludeRules []PermissionRolebindings `json:"excludeRules"`
}

type InitPermissionData struct {
	Pool      string   `json:"-"`
	Group     string   `json:"group"`     //分组，如idaas
	Name      string   `json:"name"`      //唯一标识，如idaas_init
	Version   string   `json:"version"`   //版本，每个版本初始化一下，如需更新修改版本
	Roles     []string `json:"roles"`     //角色，使用，分割元素。{Name,Role,Kind,Resource,ResourceName}
	Rules     []string `json:"rules"`     //规则，使用，分割元素。{User,UserKind,MatchResource,MatchResourceName,Resource,ResourceName,Action,ActionKind}
	Resources []string `json:"resources"` //资源，使用，分割元素。{Resource，Name}
	Actions   []string `json:"actions"`   //操作，使用，分割元素。{Resource，Action，Name}
}

type PermissionRulePutParam struct {
	Filter *PermissionRolebindings  `json:"filter"`
	Data   []PermissionRolebindings `json:"data"`
}

type EnforceParam struct {
	Group        string `json:"group"`        //分组
	Resource     string `json:"resource"`     //资源
	ResourceName string `json:"resourceName"` //资源名
	Children     bool   `json:"children"`     //获取后续资源
	Action       string `json:"action"`       //操作
	Token        string `json:"token"`        //用户token, 和（userId、deptPath）取一个
	UserId       string `json:"userId"`       //用户id
	DeptPath     string `json:"deptPath"`     //部门path
}

// JSONTime format json time field by myself
type JSONTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t *JSONTime) UnmarshalJSON(d []byte) error {
	tb, _ := time.ParseInLocation("2006-01-02 15:04:05", strings.Trim(string(d), "\""), time.Local)
	(*t).Time = tb
	return nil
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

type VerifyCaptcha struct {
	ID     string `json:"id"`     //验证码ID
	Answer string `json:"answer"` //用户填写验证码
	Clear  bool   `json:"clear"`  //是否清除缓存
}

const (
	LoginTypeWxApplet = 1
	LoginTypePhone    = 2
	LoginTypeWxWeb    = 3
)

const (
	LoginStateWait   = 1
	LoginStateScan   = 2
	LoginStateFinish = 3
)

const (
	VisitorUserType    = 1
	NotVisitorUserType = 2
)

const (
	UserNormalState = 1
	UserDelState    = 2
)
