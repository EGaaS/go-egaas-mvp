package model

type SystemState struct {
	ID   int64 `gorm:"primary_key;not null"`
	RbID int64 `gorm:"not null"`
}

func (ss *SystemState) TableName() string {
	return "system_states"
}

func GetAllSystemStatesIDs() ([]int64, error) {
	states := new([]SystemState)
	if err := DBConn.Find(&states).Order("id").Error; err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(*states))
	for _, s := range *states {
		ids = append(ids, s.ID)
	}
	return ids, nil
}

func (ss *SystemState) Get(id int64) (bool, error) {
	return isFound(DBConn.Where("id = ?", id).First(ss))
}

func (ss *SystemState) GetCount() (int64, error) {
	count := int64(0)
	err := DBConn.Table("system_states").Count(&count).Error
	return count, err
}

func (ss *SystemState) GetAllLimitOffset(limit, offset int64) ([]SystemState, error) {
	result := new([]SystemState)
	err := DBConn.Table("system_states").Order("id desc").Limit(limit).Offset(offset).Find(&result).Error
	return *result, err
}

func (ss *SystemState) GetLast() (bool, error) {
	last := DBConn.Last(ss)
	if last.RecordNotFound() {
		return true, nil
	}
	return false, last.Error
}

func (ss *SystemState) Delete() error {
	return DBConn.Delete(ss).Error
}

func (ss *SystemState) Create() error {
	return DBConn.Create(ss).Error
}
