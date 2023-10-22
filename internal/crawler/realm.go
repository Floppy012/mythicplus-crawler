package crawler

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"mythic-plus-crawler/internal/blizzapi"
	"mythic-plus-crawler/internal/database"
	"mythic-plus-crawler/internal/logs"
	"mythic-plus-crawler/internal/utils"

	"github.com/r3labs/diff/v3"
	"gorm.io/gorm/clause"
)

type idAndBlizzId interface {
	database.IGormModel
	database.IHasBlizzID
}

type realmCrawler struct {
	api *blizzapi.BlizzApi
	db  *database.Database
}

func CrawlRealms(api *blizzapi.BlizzApi, db *database.Database) error {
	log.Info().Msgf("crawling realms for region %v", api.Region.Slug)
	crawler := &realmCrawler{api, db}

	connectedRealms := make(map[uint]*database.ConnectedRealm)
	realms := make(map[uint]*database.Realm)

	pages, err := api.GetRealms()

	if err != nil {
		return fmt.Errorf("error while getting realms from api: %w", err)
	}

	for _, page := range *pages {
		for _, connectedRealm := range page.Results {
			dbConnectedRealm := &database.ConnectedRealm{
				HasBlizzID: database.HasBlizzID{
					BlizzID: uint(connectedRealm.Data.ID),
				},
				Region:     &api.Region,
				Queue:      connectedRealm.Data.HasQueue,
				Online:     strings.ToUpper(connectedRealm.Data.Status.Type) == "UP",
				Population: connectedRealm.Data.Population.Type,
			}

			connectedRealms[uint(connectedRealm.Data.ID)] = dbConnectedRealm

			for _, realm := range connectedRealm.Data.Realms {
				realms[uint(realm.ID)] = &database.Realm{
					HasBlizzID: database.HasBlizzID{
						BlizzID: uint(realm.ID),
					},
					Region:         &api.Region,
					Timezone:       realm.Timezone,
					ConnectedRealm: dbConnectedRealm,
					Name:           realm.Name.EnUS,
					Slug:           realm.Slug,
					Tournament:     realm.IsTournament,
					Locale:         realm.Locale,
					Type:           realm.Type.Type,
				}
			}
		}
	}

	err = handleUpdates(crawler, &database.ConnectedRealm{}, []string{}, &connectedRealms,
		func(element *database.ConnectedRealm) database.Log[logs.ConnectedRealmAddedPayload] {
			return logs.Create(&api.Region, logs.LogType.ConnectedRealmAdded, logs.ConnectedRealmAddedPayload{
				ID: element.GetID(),
			})
		},
		func(element *database.ConnectedRealm, changelog *diff.Changelog) database.Log[logs.ConnectedRealmUpdatedPayload] {
			return logs.Create(&api.Region, logs.LogType.ConnectedRealmUpdated, logs.ConnectedRealmUpdatedPayload{
				ID:        element.GetID(),
				Changelog: *changelog,
			})
		},
		func(element *database.ConnectedRealm) database.Log[logs.ConnectedRealmRemovedPayload] {
			return logs.Create(&api.Region, logs.LogType.ConnectedRealmRemoved, logs.ConnectedRealmRemovedPayload{
				ID: element.GetID(),
			})
		},
	)

	if err != nil {
		return fmt.Errorf("error while handling connected realms: %w", err)
	}

	err = handleUpdates(crawler, &database.Realm{}, []string{}, &realms,
		func(element *database.Realm) database.Log[logs.RealmAddedPayload] {
			return logs.Create(&api.Region, logs.LogType.RealmAdded, logs.RealmAddedPayload{
				ID: element.GetID(),
			})
		},
		func(element *database.Realm, changelog *diff.Changelog) database.Log[logs.RealmUpdatedPayload] {
			return logs.Create(&api.Region, logs.LogType.RealmUpdated, logs.RealmUpdatedPayload{
				ID:        element.GetID(),
				Changelog: *changelog,
			})
		},
		func(element *database.Realm) database.Log[logs.RealmRemovedPayload] {
			return logs.Create(&api.Region, logs.LogType.RealmRemoved, logs.RealmRemovedPayload{
				ID: element.GetID(),
			})
		},
	)

	if err != nil {
		return fmt.Errorf("error while handling realms: %w", err)
	}

	return nil
}

func handleUpdates[T idAndBlizzId, TAddedLog any, TUpdatedLog any, TRemovedLog any](
	c *realmCrawler,
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
	c.db.Gorm.Model(model).
		Where("blizz_id IN ?", blizzIds).
		Find(&existingModels)

	var updateLogs []database.Log[TUpdatedLog]
	for _, existing := range existingModels {
		update := *(*elements)[existing.GetBlizzID()]
		delete(*elements, existing.GetBlizzID())

		c.db.Gorm.Model(model).Clauses(clause.Returning{}).Where("id", existing.GetID()).Updates(update)
		// TODO why does the "clause.Returning" above not update the missing values?
		c.db.Gorm.Model(model).Where("id", existing.GetID()).Find(&update)

		changelog, err := utils.DiffDatabaseModel(existing, update)

		if err != nil {
			return fmt.Errorf("error while creating database model diff: %w", err)
		}

		if len(changelog) > 0 {
			updateLogs = append(updateLogs, updatedLogProvider(&existing, &changelog))
		}
	}

	if len(updateLogs) > 0 {
		c.db.Gorm.Create(updateLogs)
	}

	// Creates
	var creates []*T
	for _, newElement := range *elements {
		creates = append(creates, newElement)
	}

	if len(creates) > 0 {
		c.db.Gorm.Create(creates)
	}

	var addLogs []database.Log[TAddedLog]
	for _, created := range creates {
		addLogs = append(addLogs, addedLogProvider(created))
	}

	if len(addLogs) > 0 {
		c.db.Gorm.Create(addLogs)
	}

	// Removes
	var removes []T
	query := c.db.Gorm.Model(model).
		Select(selects).
		Where("region_id = ?", c.api.Region.ID)

	if len(blizzIds) > 0 {
		query = query.Where("blizz_id NOT IN ?", blizzIds)
	}

	query.Find(&removes)

	var removeLogs []database.Log[TRemovedLog]
	for _, missingElement := range removes {
		removeLogs = append(removeLogs, removedLogProvider(&missingElement))
	}

	if len(removeLogs) > 0 {
		c.db.Gorm.Create(removeLogs)
	}

	if len(removes) > 0 {
		c.db.Gorm.Delete(removes)
	}

	return nil
}
