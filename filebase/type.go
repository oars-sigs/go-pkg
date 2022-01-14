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
}
