package v1

import (
	metav1 "github.com/Pis0sion/rblogrus/meta/types"
	"time"
)

type Secret struct {

	metav1.ObjectMeta `json:"metadata,omitempty"`

	AppName         string    `json:"appName"        gorm:"column:appName"  validate:"omitempty"`
	InstanceID      int       `json:"instanceID"     gorm:"column:instanceID"  validate:"omitempty"`
	AppAccessKey    string    `json:"appAccessKey"   gorm:"column:appAccessKey"  validate:"omitempty"`
	AppAccessSecret string    `json:"appAccessSecret" gorm:"column:appAccessSecret"  validate:"omitempty"`
	Description     string    `json:"description"    gorm:"column:description"  validate:"omitempty"`
	ExpireTime      time.Time `json:"expireTime"     gorm:"column:expireTime"  validate:"omitempty"`
}

func (s *Secret) TableName() string {
	return "secret"
}


