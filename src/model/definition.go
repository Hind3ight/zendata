package model

type Definition struct {
	Title string `yaml:"title"`
	Desc string `yaml:"desc"`
	Author string `yaml:"author"`
	Version string`yaml:"version"`

	Fields  []Field `yaml:"fields,flow"`
}

type Field struct {
	Name string `yaml:"name"`
	Datatype string `yaml:"datatype"`
	Step string `yaml:"step"`
	Range string `yaml:"range"`
	Prefix string `yaml:"prefix"`
	Postfix string `yaml:"postfix"`

	Precision int
}
