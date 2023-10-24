package database

import (
	"database/sql"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type IHasBlizzID interface {
	GetBlizzID() uint
}

type HasBlizzID struct {
	BlizzID uint `gorm:"index:,unique,composite:blizz_unique"`
}

type HasRegion struct {
	RegionID uint    `gorm:"index:,unique,composite:blizz_unique"`
	Region   *Region `diff:"ignore"`
}

func (i HasBlizzID) GetBlizzID() uint {
	return i.BlizzID
}

type IGormModel interface {
	GetID() uint
}

type GormModel struct {
	gorm.Model
}

func (i GormModel) GetID() uint {
	return i.ID
}

type Region struct {
	ID                       uint   `gorm:"primaryKey"`
	BlizzID                  uint   `gorm:"index; not null"`
	Slug                     string `gorm:"index; not null"`
	Name                     string
	Active                   bool     `gorm:"not null"`
	Realms                   []*Realm `diff:"ignore"`
	ActiveMythicPlusSeasonID *uint
	ActiveMythicPlusSeason   *MythicPlusSeason
}

type ConnectedRealm struct {
	GormModel
	HasBlizzID `gorm:"embedded"`
	HasRegion  `gorm:"embedded"`
	Queue      bool
	Online     bool
	Population string
	Realms     []*Realm `diff:"ignore"`
}

type Realm struct {
	GormModel
	HasBlizzID       `gorm:"embedded"`
	HasRegion        `gorm:"embedded"`
	Timezone         string          `gorm:"index"`
	ConnectedRealmID uint            `gorm:"index"`
	ConnectedRealm   *ConnectedRealm `diff:"ignore"`
	Name             string
	Slug             string `gorm:"index"`
	Tournament       bool
	Locale           string
	Type             string
}

type MythicPlusAffix struct {
	GormModel
	HasBlizzID  `gorm:"embedded"`
	HasRegion   `gorm:"embedded"`
	Name        string
	Description string
}

type MythicPlusDungeon struct {
	GormModel
	HasBlizzID          `gorm:"embedded"`
	HasRegion           `gorm:"embedded"`
	Name                string
	MapID               uint
	MapName             string
	ZoneSlug            string
	DungeonID           uint
	DungeonName         string
	QualifierTime       int
	QualifierDoubleTime int
	QualifierTripleTime int
}

type MythicPlusSeason struct {
	GormModel
	HasBlizzID        `gorm:"embedded"`
	HasRegion         `gorm:"embedded"`
	StartTime         time.Time
	EndTime           sql.NullTime
	MythicPlusPeriods []*MythicPlusPeriod
}

type MythicPlusPeriod struct {
	GormModel
	HasBlizzID         HasBlizzID `gorm:"embedded"`
	MythicPlusSeasonID uint       `gorm:"index:,composite:blizz_unique"`
	MythicPlusSeason   *MythicPlusSeason
	StartTime          time.Time
	EndTime            time.Time
}

type Log[T any] struct {
	ID        uint
	RegionID  uint `gorm:"index"`
	Region    *Region
	Type      string `gorm:"index"`
	Payload   datatypes.JSONType[T]
	CreatedAt time.Time
}

func (Log[T]) TableName() string {
	return "logs"
}
