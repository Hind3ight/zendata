package model

import (
	"time"
)

var (
	CommonPrefix = "zd_"
	Models = []interface{}{ &Def{}, &Field{} }
)

type ReqData struct {
	Action string `json:"action"`
	Id int  `json:"id"`
	Mode string `json:"mode"`
	Data interface{} `json:"data"`

	Src int `json:"src"`
	Dist int `json:"dist"`
}

type ResData struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
	Field interface{} `json:"field"`
}

type Model struct {
	ID        uint      `gorm:"column:id;primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:createTime" json:"createTime" yaml:"-"`
	UpdatedAt time.Time `gorm:"column:updateTime" json:"updateTime" yaml:"-"`

	Disabled bool `gorm:"column:disabled;default:false" json:"disabled" yaml:"-"`
	Deleted  bool `gorm:"column:deleted;default:false" json:"deleted" yaml:"-"`
}
