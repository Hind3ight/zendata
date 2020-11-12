package model

import (
	constant "github.com/easysoft/zendata/src/utils/const"
)

type Def struct {
	Model
	Title  string `gorm:"column:title" json:"title"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`
	Type   string `gorm:"column:type" json:"type"`
	Desc   string `gorm:"column:desc" json:"desc"`
	Yaml   string `gorm:"yaml" json:"yaml"`
	Folder string `gorm:"-" json:"folder" yaml:"-"`
}
func (*Def) TableName() string {
	return constant.TablePrefix + "def"
}

type Field struct {
	Model
	DefID uint `gorm:"column:defID" json:"defID"`
	ParentID uint `gorm:"column:parentID" json:"parentID"`
	Field string `gorm:"column:field" json:"field"`
	Note string `gorm:"column:note" json:"note"`

	Range string `gorm:"column:range" json:"range"`
	Exp  string `gorm:"column:exp" json:"exp"`
	Prefix string `gorm:"column:prefix" json:"prefix"`
	Postfix string `gorm:"column:postfix" json:"postfix"`
	Loop string `gorm:"column:loop" json:"loop"`
	Loopfix string `gorm:"column:loopfix" json:"loopfix"`
	Format string `gorm:"column:format" json:"format"`
	Type string `gorm:"column:type" json:"type"`
	Mode string `gorm:"column:mode" json:"mode"`
	Length int `gorm:"column:length" json:"length"`
	LeftPad string `gorm:"column:leftPad" json:"leftPad"`
	RightPad string `gorm:"column:rightPad" json:"rightPad"`
	Rand bool `gorm:"column:rand" json:"rand"`

	ConfigID	uint `gorm:"column:configID" json:"configID"`
	Config	string `gorm:"column:config" json:"config"`
	Use	string `gorm:"column:use" json:"use"`
	UseID	uint `gorm:"column:useID" json:"useID"`
	From	string `gorm:"column:fromCol" json:"fromCol"`
	Select	string `gorm:"column:selectCol" json:"selectCol"`
	Where	string `gorm:"column:whereCol" json:"whereCol"`
	Limit	int `gorm:"column:limitCol" json:"limitCol"`

	Ord int `gorm:"column:ord;default:1" json:"ord"`

	Children []*Field `gorm:"-" json:"children"`
	Froms []*Field `gorm:"-" json:"froms"`

	IsRange bool `gorm:"column:isRange;default:true" json:"isRange"`
	Sections []Section `gorm:"column:sections;ForeignKey:fieldID" json:"sections"`
}
func (*Field) TableName() string {
	return constant.TablePrefix + "field"
}

type Section struct {
	Model
	FieldID uint   `gorm:"column:fieldID" json:"fieldID"`
	Type    string `gorm:"column:type;default:scope" json:"type"`
	Value     string `gorm:"column:value" json:"value"`
	Ord     int    `gorm:"column:ord;default:1" json:"ord"`

	// for range
	Start string `gorm:"column:start" json:"start"`
	End string `gorm:"column:end" json:"end"`
	Step string `gorm:"column:step;default:1" json:"step"`
	Repeat string `gorm:"column:repeat;default:1" json:"repeat"`
	Rand bool `gorm:"column:rand;default:false" json:"rand"`

	// for arr and const
	Text string `gorm:"-" json:"-"`
}
func (*Section) TableName() string {
	return constant.TablePrefix + "section"
}
