package crawler

import (
	"fmt"
	"strings"

	"github.com/r3labs/diff/v3"
	"github.com/rs/zerolog/log"
	"mythic-plus-crawler/internal/blizzapi"
	"mythic-plus-crawler/internal/database"
	"mythic-plus-crawler/internal/logs"
)

type realmCrawler struct {
	api *blizzapi.BlizzApi
	db  *database.Database
}

func CrawlRealms(api *blizzapi.BlizzApi, db *database.Database) error {
	log.Info().Msgf("crawling realms for region %v", api.Region.Slug)

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
				HasRegion:  database.HasRegion{Region: &api.Region},
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
					HasRegion:      database.HasRegion{Region: &api.Region},
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

	err = HandleDBUpdates(db, &api.Region, &database.ConnectedRealm{}, []string{}, &connectedRealms,
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

	err = HandleDBUpdates(db, &api.Region, &database.Realm{}, []string{}, &realms,
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
