package utils

import (
	"errors"
	"reflect"

	"github.com/FuzzyStatic/blizzard/v3"
	"github.com/r3labs/diff/v3"
)

func DiffDatabaseModel(modelA interface{}, modelB interface{}) (diff.Changelog, error) {
	return diff.Diff(modelA, modelB, diff.Filter(
		func(path []string, parent reflect.Type, field reflect.StructField) bool {
			if val, ok := field.Tag.Lookup("diff"); ok {
				if val == "ignore" {
					return false
				}
			}
			switch p := path[0]; p {
			case "GormModel", "ID", "CreatedAt", "UpdatedAt", "DeletedAt":
				return false
			default:
				return true
			}
		},
	))
}

func BlizzardRegionFromString(strRegion string) (blizzard.Region, error) {
	switch strRegion {
	case "us":
		return blizzard.US, nil
	case "eu":
		return blizzard.EU, nil
	case "kr":
		return blizzard.KR, nil
	case "tw":
		return blizzard.TW, nil
	case "cn":
		return blizzard.CN, nil
	default:
		return -1, errors.New("unknown string region")
	}
}
