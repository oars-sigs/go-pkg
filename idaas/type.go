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
