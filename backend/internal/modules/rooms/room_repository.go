package rooms

import "gorm.io/gorm"

type RoomRepositoryInterface interface {
	GetAll() ([]Room, error)
	GetById(id string) (*Room, error)
	Create(room *Room) error
	Update(room *Room, id string) error
	Delete(id string) error
	GetLastRoom() (*Room, error)
	FilerByCategory(categoryId string) ([]Room, error)
	FilterByType(typeId string) ([]Room, error)
}

type roomRepo struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepositoryInterface {
	return &roomRepo{db: db}
}

func (r *roomRepo) GetAll() ([]Room, error) {
	var rooms []Room

	err := r.db.Preload("Category").Preload("Type").Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepo) GetById(id string) (*Room, error) {
	var room Room
	err := r.db.Preload("Category").Preload("Type").Where("id = ?", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepo) Create(room *Room) error {
	return r.db.Omit("Category", "Type").Create(room).Error
}

func (r *roomRepo) Update(room *Room, id string) error {
	return r.db.Model(&Room{}).Where("id = ?", id).Omit("Category", "Type").Updates(room).Error
}

func (r *roomRepo) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&Room{}).Error
}

func (r *roomRepo) GetLastRoom() (*Room, error) {
	var room Room
	err := r.db.Order("created_at desc").First(&room).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *roomRepo) FilerByCategory(categoryId string) ([]Room, error) {

	var rooms []Room

	err := r.db.Preload("Type").Preload("Category").Where("category_id = ?", categoryId).Find(&rooms).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rooms, nil

}

func (r *roomRepo) FilterByType(typeId string) ([]Room, error) {

	var rooms []Room

	err := r.db.Preload("Type").Preload("Category").Where("type_id = ?", typeId).Find(&rooms).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rooms, nil

}
