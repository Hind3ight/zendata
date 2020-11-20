package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

type DefService struct {
	defRepo *serverRepo.DefRepo
	fieldRepo *serverRepo.FieldRepo
	referRepo *serverRepo.ReferRepo
	resService *ResService
}

func (s *DefService) List(keywords string, page int) (list []*model.ZdDef, total int) {
	list, total, _ = s.defRepo.List(strings.TrimSpace(keywords), page)
	return
}

func (s *DefService) Get(id int) (def model.ZdDef, dirTree model.Dir) {
	def, _ = s.defRepo.Get(uint(id))

	dirTree = model.Dir{Name: fileUtils.AddSepIfNeeded(constant.ResDirUsers)}
	serverUtils.GetDirTree(&dirTree)

	return
}

func (s *DefService) Save(def *model.ZdDef) (err error) {
	def.Folder = serverUtils.DealWithPathSepRight(def.Folder)
	def.Path = vari.WorkDir + def.Folder + serverUtils.AddExt(def.Title, ".yaml")

	if def.ID == 0 {
		err = s.Create(def)
	} else {
		err = s.Update(def)
	}
	return
}

func (s *DefService) Create(def *model.ZdDef) (err error) {
	err = s.defRepo.Create(def)

	rootField, err := s.fieldRepo.CreateTreeNode(def.ID, 0, "字段", "root")
	s.referRepo.CreateDefault(rootField.ID, constant.ResTypeDef)

	s.dataToYaml(def)
	err = s.defRepo.Update(def)

	return
}

func (s *DefService) Update(def *model.ZdDef) (err error) {
	var old model.ZdDef
	old, err = s.defRepo.Get(def.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if def.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	s.dataToYaml(def)
	err = s.defRepo.Update(def)

	return
}

func (s *DefService) Remove(id int) (err error) {
	var old model.ZdDef
	old, err = s.defRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}
	fileUtils.RemoveExist(old.Path)

	err = s.defRepo.Remove(uint(id))
	return
}

func (s *DefService) UpdateYaml(defId uint) (err error) {
	var def model.ZdDef
	def, _ = s.defRepo.Get(defId)

	s.dataToYaml(&def)
	err = s.defRepo.UpdateYaml(def)

	return
}

func (s *DefService) dataToYaml(def *model.ZdDef) (str string) {
	root, err := s.fieldRepo.GetDefFieldTree(def.ID)
	if err != nil {
		return
	}

	defData := model.DefData{}
	s.defRepo.GenDef(*def, &defData)

	for _, child := range root.Fields { // ignore the root
		defField := model.DefField{}
		convertToConfModel(*child, &defField)

		defData.Fields = append(defData.Fields, defField)
	}

	bytes, err := yaml.Marshal(defData)
	def.Yaml = string(bytes)

	return
}

func (s *DefService) Sync(files []model.ResFile) (err error) {
	defList := s.defRepo.ListAll()

	defMap := map[string]*model.ZdDef{}
	for _, item := range defList {
		defMap[item.Path] = item
	}

	for _, fi := range files {
		_, found := defMap[fi.Path]
		if found {
			s.defRepo.Remove(defMap[fi.Path].ID)
		}
		if !found || fi.UpdatedAt.Unix() > defMap[fi.Path].UpdatedAt.Unix() {
			s.SyncToDB(fi)
		}
	}

	return
}
func (s *DefService) saveFieldToDB(item *model.ZdField, parentID, defID uint) {
	item.DefID = defID
	item.ParentID = parentID
	s.fieldRepo.Save(item)

	for _, child := range item.Fields {
		s.saveFieldToDB(child, item.ID, defID)
	}
}

func (s *DefService) SyncToDB(fi model.ResFile) (err error) {
	content, _ := ioutil.ReadFile(fi.Path)
	yamlContent := stringUtils.ReplaceSpecialChars(content)
	defPo := model.ZdDef{}
	err = yaml.Unmarshal(yamlContent, &defPo)
	defPo.Title = fi.Title
	defPo.Type = fi.ResType
	defPo.Desc = fi.Desc
	defPo.Path = fi.Path
	defPo.Folder = serverUtils.GetRelativePath(defPo.Path)
	defPo.Yaml = string(content)

	s.defRepo.Create(&defPo)

	rootField, _ := s.fieldRepo.CreateTreeNode(defPo.ID, 0, "字段", "root")
	s.referRepo.CreateDefault(rootField.ID, constant.ResTypeDef)
	for i, field := range defPo.Fields {
		field.Ord = i + 1
		s.saveFieldToDB(&field, rootField.ID, defPo.ID)
	}

	return
}
func (s *DefService) SyncFromDB(files []model.ResFile) (err error) {

	return
}

func NewDefService(defRepo *serverRepo.DefRepo, fieldRepo *serverRepo.FieldRepo,
	referRepo *serverRepo.ReferRepo) *DefService {
	return &DefService{defRepo: defRepo, fieldRepo: fieldRepo, referRepo: referRepo}
}
