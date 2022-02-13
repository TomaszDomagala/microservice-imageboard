package identification

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Service interface {
	Identify(ip string) (uint, error)
}

type service struct {
	db *gorm.DB
}

const idLifeTime = time.Minute

func (s *service) Identify(ip string) (uint, error) {
	t := time.Now().Add(-idLifeTime)
	var identification Identification

	res := s.db.Where("ip = ? AND updated_at > ?", ip, t).Take(&identification)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return 0, res.Error
	}

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		identification = Identification{
			Ip: ip,
		}
		res = s.db.Create(&identification)
		if res.Error != nil {
			return 0, res.Error
		}
	}

	return identification.ID, nil
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}
