package database

import (
	"errors"
	"fmt"

	"mythic-plus-crawler/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	Gorm *gorm.DB
}

func Connect(config *config.Config) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%v user=%v dbname=%v port=%v password=%v sslmode=disable TimeZone=Europe/Berlin",
		config.Database.Host,
		config.Database.User,
		config.Database.Database,
		config.Database.Port,
		config.Database.Password,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("error while connecting to database: %w", err)
	}

	err = Migrate(db)

	if err != nil {
		return nil, fmt.Errorf("error while performing migrations: %w", err)
	}

	return &Database{
		Gorm: db,
	}, nil
}

func (db *Database) Exists(query interface{}, args ...interface{}) bool {
	var exists bool
	err := db.Gorm.Model(query).
		Select("count(id) > 0").
		Where(query, args).
		Limit(1).
		Find(&exists).
		Error
	if err == nil {
		return exists
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	panic(err)
}

func (db *Database) Upsert(query interface{}, update interface{}) {
	if db.Exists(query) {
		db.Gorm.Clauses(clause.Returning{}).Where(query).Updates(update)
		return
	}

	db.Gorm.Save(update)
}
