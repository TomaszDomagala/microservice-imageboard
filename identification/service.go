package identification

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Service interface {
	Identify(ip string) (string, error)
}

type service struct {
	db *gorm.DB
}

const idLifeTime = time.Minute

func (s *service) Identify(ip string) (string, error) {
	t := time.Now().Add(-idLifeTime)
	var identification Identification

	res := s.db.Where("ip = ? AND updated_at > ?", ip, t).Take(&identification)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return "", res.Error
	}

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		identification = Identification{
			Ip: ip,
		}
		res = s.db.Create(&identification)
		if res.Error != nil {
			return "", res.Error
		}
	}

	return strconv.Itoa(int(identification.ID)), nil
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}
