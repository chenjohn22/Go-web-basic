package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Client struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	Account   string    `gorm:"column:account" json:"account"`
	Password  string    `gorm:"column:password" json:"password"`
	IP        string    `gorm:"column:ip" json:"ip"`
	Status    int       `gorm:"column:status" json:"status"`
	Balance   string    `gorm:"column:balance" json:"balance"`
	Token     string    `gorm:"column:token" json:"token"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (c *Client) TableName() string {
	return "client"
}

func (c *Client) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", int(time.Now().Unix()))
	return nil
}
