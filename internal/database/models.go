package database

import (
	"time"

	"gorm.io/gorm"
)

type Region struct {
	ID        uint
	BlizzID   uint   `gorm:"uniqueIndex"`
	Name      string `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Realms    []Realm
}

type Timezone struct {
	ID        uint
	Timezone  string `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Realms    []Realm
}

type Realm struct {
	ID               uint
	BlizzID          uint `gorm:"uniqueIndex"`
	RegionID         uint `gorm:"index"`
	Region           Region
	TimezoneID       uint `gorm:"index"`
	Timezone         Timezone
	ConnectedRealmID uint
	ConnectedRealm   ConnectedRealm
	Name             string
	Slug             string `gorm:"index"`
	Tournament       bool
	Locale           string
	Type             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

type ConnectedRealm struct {
	ID        uint
	BlizzID   uint `gorm:"uniqueIndex"`
	Queue     bool
	Status    bool
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Realms    []Realm
}
