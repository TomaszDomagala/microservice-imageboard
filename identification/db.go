package identification

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const postgresUser = "postgres"
const postgresPassword = "postgres"
const postgresDb = "identification"
const host = "identification_db"
const port = 5432

func getPostgresConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, postgresUser, postgresPassword, postgresDb)
}

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(getPostgresConnectionString()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	err = db.AutoMigrate(&Identification{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}
	return db, nil
}

type Identification struct {
	gorm.Model
	Ip string `gorm:"type:varchar(255);index"`
}
