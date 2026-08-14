package channel

import "gorm.io/gorm"

type ChannelRepository interface {
	FindAll() ([]Channel, error)
	Find(page int) ([]Channel, int, error)
	Create(channel *Channel) error
	FindById(id string) (*Channel, error)
	FindByName(name string) (*Channel, error)
	Update(channel *Channel, id string) error
	Delete(id string) error
	GetLastChannel() (*Channel, error)
}

type channelRepo struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepo{db: db}
}

func (r *channelRepo) FindAll() ([]Channel, error) {
	var channels []Channel

	err := r.db.Find(&channels).Error
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (r *channelRepo) Find(page int) ([]Channel, int, error) {
	var channels []Channel
	limit := 10
	offset := (page - 1) * limit
	var total int64

	err := r.db.Model(&Channel{}).Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	err = r.db.Limit(limit).Offset(offset).Find(&channels).Error
	if err != nil {
		return nil, 0, err
	}

	return channels, int(total), nil
}

func (r *channelRepo) FindById(id string) (*Channel, error) {
	var channel Channel

	err := r.db.First(&channel, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &channel, nil
}

func (r *channelRepo) Create(channel *Channel) error {
	err := r.db.Create(&channel).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *channelRepo) FindByName(name string) (*Channel, error) {
	var channel Channel

	err := r.db.Where("name = ?", name).First(&channel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &channel, nil
}

func (r *channelRepo) Update(channel *Channel, id string) error {
	err := r.db.Model(&Channel{}).Where("id = ?", id).Updates(channel).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *channelRepo) Delete(id string) error {
	err := r.db.Delete(&Channel{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *channelRepo) GetLastChannel() (*Channel, error) {
	var channel Channel
	err := r.db.Order("created_at desc").First(&channel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &channel, nil
}
