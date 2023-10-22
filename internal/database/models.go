package database

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type IHasBlizzID interface {
	GetBlizzID() uint
}

type HasBlizzID struct {
	BlizzID uint `gorm:"uniqueIndex"`
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
	ID      uint   `gorm:"primaryKey"`
	BlizzID uint   `gorm:"index; not null"`
	Slug    string `gorm:"index; not null"`
	Name    string
	Active  bool     `gorm:"not null"`
	Realms  []*Realm `diff:"ignore"`
}

type ConnectedRealm struct {
	GormModel
	HasBlizzID `gorm:"embedded"`
	RegionID   uint
	Region     *Region `diff:"ignore"`
	Queue      bool
	Online     bool
	Population string
	Realms     []*Realm `diff:"ignore"`
}

type Realm struct {
	GormModel
	HasBlizzID       `gorm:"embedded"`
	RegionID         uint            `gorm:"index"`
	Region           *Region         `diff:"ignore"`
	Timezone         string          `gorm:"index"`
	ConnectedRealmID uint            `gorm:"index"`
	ConnectedRealm   *ConnectedRealm `diff:"ignore"`
	Name             string
	Slug             string `gorm:"index"`
	Tournament       bool
	Locale           string
	Type             string
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
