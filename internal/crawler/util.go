package crawler

import (
	"fmt"

	"github.com/r3labs/diff/v3"
	"gorm.io/gorm/clause"
	"mythic-plus-crawler/internal/database"
	"mythic-plus-crawler/internal/logs"
	"mythic-plus-crawler/internal/utils"
)

type idAndBlizzId interface {
	database.IGormModel
	database.IHasBlizzID
}

func HandleDBUpdates[T idAndBlizzId, TAddedLog any, TUpdatedLog any, TRemovedLog any](
	db *database.Database,
	region *database.Region,
	model interface{},
	additionalSelects []string,
	elements *map[uint]*T,
	addedLogProvider func(element *T) database.Log[TAddedLog],
	updatedLogProvider func(element *T, changelog *diff.Changelog) database.Log[TUpdatedLog],
	removedLogProvider func(element *T) database.Log[TRemovedLog],
) error {
	var blizzIds []uint

	selects := []string{"id", "blizz_id"}
	selects = append(selects, additionalSelects...)

	for blizzId := range *elements {
		blizzIds = append(blizzIds, blizzId)
	}

	// Updates
	var existingModels []T
	db.Gorm.Model(model).
		Where("region_id = ?", region.ID).
		Where("blizz_id IN ?", blizzIds).
		Find(&existingModels)

	var updateLogs []database.Log[TUpdatedLog]
	for _, existing := range existingModels {
		update := *(*elements)[existing.GetBlizzID()]
		delete(*elements, existing.GetBlizzID())

		db.Gorm.Model(model).Clauses(clause.Returning{}).Where("id", existing.GetID()).Updates(update)
		// TODO why does the "clause.Returning" above not update the missing values?
		db.Gorm.Model(model).Where("id", existing.GetID()).Find(&update)

		changelog, err := utils.DiffDatabaseModel(existing, update)

		if err != nil {
			return fmt.Errorf("error while creating database model diff: %w", err)
		}

		if len(changelog) > 0 {
			updateLogs = append(updateLogs, updatedLogProvider(&existing, &changelog))
		}
	}

	if len(updateLogs) > 0 {
		db.Gorm.Create(updateLogs)
	}

	// Creates
	var creates []*T
	for _, newElement := range *elements {
		creates = append(creates, newElement)
	}

	if len(creates) > 0 {
		db.Gorm.Create(creates)
	}

	var addLogs []database.Log[TAddedLog]
	for _, created := range creates {
		addLogs = append(addLogs, addedLogProvider(created))
	}

	if len(addLogs) > 0 {
		db.Gorm.Create(addLogs)
	}

	// Removes
	var removes []T
	query := db.Gorm.Model(model).
		Select(selects).
		Where("region_id = ?", region.ID)

	if len(blizzIds) > 0 {
		query = query.Where("blizz_id NOT IN ?", blizzIds)
	}

	query.Find(&removes)

	var removeLogs []database.Log[TRemovedLog]
	for _, missingElement := range removes {
		removeLogs = append(removeLogs, removedLogProvider(&missingElement))
	}

	if len(removeLogs) > 0 {
		db.Gorm.Create(removeLogs)
	}

	if len(removes) > 0 {
		db.Gorm.Delete(removes)
	}

	return nil
}

func HandleDBUpdatesGenericLogs[T idAndBlizzId](
	db *database.Database,
	region *database.Region,
	model interface{},
	additionalSelects []string,
	elements *map[uint]*T,
	addedLogType string,
	updatedLogType string,
	removedLogType string,
) error {
	return HandleDBUpdates(db, region, model, additionalSelects, elements,
		func(element *T) database.Log[logs.GenericAdded] {
			return logs.Create(region, addedLogType, logs.GenericAdded{
				ID: (*element).GetID(),
			})
		},
		func(element *T, changelog *diff.Changelog) database.Log[logs.GenericUpdated] {
			return logs.Create(region, updatedLogType, logs.GenericUpdated{
				ID:        (*element).GetID(),
				Changelog: *changelog,
			})
		},
		func(element *T) database.Log[logs.GenericRemoved] {
			return logs.Create(region, removedLogType, logs.GenericRemoved{
				ID: (*element).GetID(),
			})
		},
	)
}
