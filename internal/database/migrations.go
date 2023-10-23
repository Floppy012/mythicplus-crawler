package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func initSchema(tx *gorm.DB) error {
	err := tx.AutoMigrate(
		&Region{},
		&ConnectedRealm{},
		&Realm{},
		&MythicPlusAffix{},
		&MythicPlusDungeon{},
		&Log[any]{},
	)

	if err != nil {
		return err
	}

	tx.Create(&[]Region{
		{
			BlizzID: 1,
			Slug:    "us",
			Name:    "North America",
			Active:  true,
		},
		{
			BlizzID: 2,
			Slug:    "kr",
			Name:    "Korea",
			Active:  false,
		},
		{
			BlizzID: 3,
			Slug:    "eu",
			Name:    "Europe",
			Active:  true,
		},
		{
			BlizzID: 4,
			Slug:    "tw",
			Name:    "Taiwan",
			Active:  false,
		},
		{
			BlizzID: 5,
			Slug:    "cn",
			Name:    "China",
			Active:  false,
		},
	})

	return nil
}

func Migrate(db *gorm.DB) error {
	migrator := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{})

	migrator.InitSchema(initSchema)

	return migrator.Migrate()
}
