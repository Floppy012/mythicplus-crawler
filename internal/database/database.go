package database

import (
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

	_ = db.AutoMigrate(&Region{})
	_ = db.AutoMigrate(&Timezone{})
	_ = db.AutoMigrate(&Realm{})
	_ = db.AutoMigrate(&ConnectedRealm{})

	return &Database{
		Gorm: db,
	}, nil
}

func (db *Database) UpsertRegion(regionId int, regionName string) *Region {
	region := Region{
		BlizzID: uint(regionId),
		Name:    regionName,
	}

	db.Gorm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "blizz_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&region)

	return &region
}
