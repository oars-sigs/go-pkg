package idaas

type UserInfo struct {
	UserId      string      `json:"userId"`
	PoolId      string      `json:"poolId"`
	Username    string      `json:"username"`
	NewPassword string      `json:"newPassword"`
	NickName    string      `json:"nickName"`
	Gender      string      `json:"gender"`
	DeptId      string      `json:"deptId"`
	PostId      string      `json:"postId"`
	RoleId      string      `json:"roleId"`
	Mobile      string      `json:"mobile"`
	Telephone   string      `json:"telephone"`
	Email       string      `json:"email"`
	Avatar      string      `json:"avatar"`
	State       int         `json:"state"`
	OrderNum    int         `json:"orderNum"`
	Created     int64       `json:"created"`
	Updated     int64       `json:"updated"`
	Kind        int         `json:"kind"`
	Dept        *Department `json:"dept"`
	AppID       string      `json:"appId"`
	RealId      string      `json:"realId"`
	Verified    bool        `json:"verified"`
	SearchText  string      `json:"searchText"`
	Type        string      `json:"type"`
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
}

type ThirdUser struct {
	TUserId    string `json:"tuser_id" gorm:"column:tuser_id"`
	PoolId     string `json:"poolId" gorm:"column:pool_id"`
	UserId     string `json:"userId" gorm:"column:user_id"`
	OpenId     string `json:"openid" gorm:"column:openid"`
	LoginType  int    `json:"login_type" gorm:"column:login_type"`
	SessionKey string `json:"sessionKey" gorm:"-"`
	Username   string `json:"username" gorm:"-"`
	Mobile     string `json:"mobile" gorm:"-"`
	Email      string `json:"email" gorm:"-"`
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
