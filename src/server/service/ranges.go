package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type RangesService struct {
	rangesRepo *serverRepo.RangesRepo
	resService *ResService
}

func (s *RangesService) List() (list []*model.ZdRanges) {
	ranges := s.resService.LoadRes("ranges")
	list, _ = s.rangesRepo.List()

	s.saveResToDB(ranges, list)
	list, _ = s.rangesRepo.List()

	return
}

func (s *RangesService) Get(id int) (ranges model.ZdRanges) {
	ranges, _ = s.rangesRepo.Get(uint(id))

	return
}

func (s *RangesService) Save(ranges *model.ZdRanges) (err error) {
	err = s.rangesRepo.Save(ranges)

	return
}

func (s *RangesService) Remove(id int) (err error) {
	var old model.ZdRanges
	old, err = s.rangesRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}
	fileUtils.RemoveExist(old.Path)

	err = s.rangesRepo.Remove(uint(id))
	return
}

func (s *RangesService) GetItemTree(rangesId int) (root model.ZdRangesItem) {
	items, _ := s.rangesRepo.GetItems(rangesId)

	root.ID = 0
	root.Name = "序列"
	for _, item := range items {
		item.ParentID = root.ID
		root.Children = append(root.Children, item)
	}

	return
}
func (s *RangesService) GetItem(id int) (item model.ZdRangesItem) {
	item, _ = s.rangesRepo.GetItem(uint(id))
	return
}

func (s *RangesService) CreateItem(domainId, targetId int, mode string) (item *model.ZdRangesItem, err error) {
	item = &model.ZdRangesItem{Name: "ranges_", RangesID: uint(domainId)}
	item.Ord = s.rangesRepo.GetMaxOrder(domainId)

	err = s.rangesRepo.SaveItem(item)

	return
}
func (s *RangesService) SaveItem(item *model.ZdRangesItem) (err error) {
	err = s.rangesRepo.SaveItem(item)
	return
}

func (s *RangesService) RemoveItem(id int) (err error) {
	err = s.rangesRepo.RemoveItem(uint(id))
	return
}

func (s *RangesService) saveResToDB(ranges []model.ResFile, list []*model.ZdRanges) (err error) {
	names := make([]string, 0)
	for _, item := range list {
		names = append(names, item.Path)
	}

	for _, item := range ranges {
		if !stringUtils.FindInArrBool(item.Path, names) {
			content, _ := ioutil.ReadFile(item.Path)
			yamlContent := stringUtils.ReplaceSpecialChars(content)
			ranges := model.ZdRanges{}
			err = yaml.Unmarshal(yamlContent, &ranges)
			ranges.Title = item.Title
			ranges.Name = item.Name
			ranges.Desc = item.Desc
			ranges.Path = item.Path
			ranges.Field = item.Title
			ranges.Note = item.Desc
			ranges.Yaml = string(content)

			s.rangesRepo.Save(&ranges)

			i := 1
			for k, v := range ranges.RangeMap {
				item := model.ZdRangesItem{Name: k, Value: v}
				item.RangesID = ranges.ID
				item.Ord = i
				s.rangesRepo.SaveItem(&item)
				i += 1
			}
		}
	}

	return
}

func NewRangesService(rangesRepo *serverRepo.RangesRepo) *RangesService {
	return &RangesService{rangesRepo: rangesRepo}
}
