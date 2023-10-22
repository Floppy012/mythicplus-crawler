package logs

import (
	"mythic-plus-crawler/internal/database"

	"gorm.io/datatypes"
)

func Create[T any](region *database.Region, logType string, payload T, dest ...*[]database.Log[T]) database.Log[T] {
	out := database.Log[T]{
		Type:    logType,
		Region:  region,
		Payload: datatypes.NewJSONType(payload),
	}

	if len(dest) > 0 {
		*dest[0] = append(*dest[0], out)
	}

	return out
}
