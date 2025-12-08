package filebase

type FileMetadata struct {
	ID        string `gorm:"column:file_id" json:"id"`
	Name      string `gorm:"column:file_name" json:"name"`
	Kind      string `gorm:"column:file_kind" json:"kind"`
	Type      string `gorm:"column:file_type" json:"type"`
	Size      int64  `gorm:"column:file_size" json:"size"`
	Digest    string `gorm:"column:file_digest" json:"digest"`
	Namespace string `gorm:"column:file_namespace" json:"namespace"`
	Parent    string `gorm:"column:file_parent" json:"parent"`
	Driver    string `gorm:"column:file_driver" json:"driver"`
	Status    string `gorm:"column:file_status" json:"status"`
	Creator   string `gorm:"column:file_creator" json:"creator"`
	Version   string `gorm:"column:file_version" json:"version"`
	Main      int    `gorm:"column:file_main" json:"main"`
	Created   int64  `gorm:"column:file_created" json:"created"`
	Updated   int64  `gorm:"column:file_updated" json:"updated"`
	DirPath   string `gorm:"-" json:"dirPath"`
}

const (
	NamespaceRoleOwner     = "owner"
	NamespaceRoleOnlyRead  = "r"
	NamespaceRoleOnlyWrite = "w"
	NamespaceRoleReadWrite = "rw"
)

const (
	NamespaceActionMgr   = "mgr"
	NamespaceActionRead  = "r"
	NamespaceActionWrite = "w"
)

const (
	FileKind = "file"
	DirKind  = "dir"
	LinkKind = "link"

	MainVersion  = 1
	MinorVersion = 2

	NormalStatus = "1"
	DeleteStatus = "2"
	ExistStatus  = "3"

	FilePrivate  = 1
	FileInternal = 2
	FilePublic   = 3
)

// Namespace 项目
type Namespace struct {
	ID       string `gorm:"column:namespace_id" json:"id"`
	Name     string `gorm:"column:namespace_name" json:"name"`
	Class    int    `gorm:"column:namespace_class" json:"class"`
	AppId    string `gorm:"column:namespace_app_id" json:"appId"`
	IsNotice int    `gorm:"column:namespace_notice" json:"isNotice"`
	Created  int64  `gorm:"column:namespace_created" json:"created"`
	Updated  int64  `gorm:"column:namespace_updated"  json:"updated"`
}

// NamespaceMember 项目成员
type NamespaceMember struct {
	ID          string `gorm:"column:nm_id" json:"id"`
	NamespaceID string `gorm:"column:nm_namespace_id" json:"namespaceId"`
	MemberID    string `gorm:"column:nm_member_id" json:"memberId"`
	Kind        string `gorm:"column:nm_kind" json:"kind"`
	Role        string `gorm:"column:nm_role" json:"role"`
	Created     int64  `gorm:"column:nm_created" json:"created"`
	Updated     int64  `gorm:"column:nm_updated" json:"updated"`
}

type URLFile struct {
	URL        string `json:"url"`
	Name       string `json:"name"`
	Visibility int    `json:"visibility"`
}
