package model

type ClsBase struct {
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	Author  string `yaml:"author"`
	Version string `yaml:"version"`
}

// range res
type ResRanges struct {
	ClsBase   `yaml:",inline"`
	FieldBase   `yaml:",inline"`
	Field string  `yaml:"field"`
	Ranges map[string]string  `yaml:"ranges"`
}
// instance res
type ResInsts struct {
	ClsBase   `yaml:",inline"`
	Field string        `yaml:"field"`
	Instances []ResInst `yaml:"instances,flow"`
}
type ResInst struct {
	FieldBase   `yaml:",inline"`
	Instance string  `yaml:"instance"`
	Range    string  `yaml:"range"`
	Fields  []DefField `yaml:"fields,flow"`
}

// common item
type DefData struct {
	ClsBase   `yaml:",inline"`
	Fields  []DefField `yaml:"fields,flow"`
}
type DefField struct {
	FieldBase   `yaml:",inline"`
	Field     string  `yaml:"field"`
	Range    string  `yaml:"range"`
	Fields   []DefField `yaml:"fields,flow"`

	Type     string
}

type FieldBase struct {
	Note     string  `yaml:"note"`

	From	string  `yaml:"from"`
	Select	string  `yaml:"select"`
	Where	string  `yaml:"where"`
	Use	string  `yaml:"use"`

	Prefix   string  `yaml:"prefix"`
	Postfix  string  `yaml:"postfix"`
	Loop  int  `yaml:"loop"`
	Loopfix  string  `yaml:"loopfix"`
	Format  string  `yaml:"format"`
	IsNumb  bool  `yaml:"isNumb"`
	Expect  string  `yaml:"expect"`

	Precision int
}

type FieldValue struct {
	FieldBase   `yaml:",inline"`
	Field     string  `yaml:"field"`
	Values   []interface{}
}