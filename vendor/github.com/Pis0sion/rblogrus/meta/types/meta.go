package types

import "time"

type ObjectMeta struct {
	// Primary key
	ID uint64 `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	// createAt
	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"column:createdAt"`
	// updateAt
	UpdatedAt time.Time `json:"updatedAt,omitempty" gorm:"column:updatedAt"`
}

type ListMeta struct {
	TotalCount int64 `json:"totalCount,omitempty"`
}
