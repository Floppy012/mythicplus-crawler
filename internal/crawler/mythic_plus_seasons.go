package crawler

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/FuzzyStatic/blizzard/v3/wowgd"
	"github.com/rs/zerolog/log"
	"mythic-plus-crawler/internal/blizzapi"
	"mythic-plus-crawler/internal/database"
	"mythic-plus-crawler/internal/logs"
)

func CrawlMythicPlusSeasons(api *blizzapi.BlizzApi, db *database.Database) error {
	log.Info().Msgf("crawling mythic plus seasons for region %v", api.Region.Slug)

	seasonsIndex, err := api.GetMPlusSeasonIndex()

	if err != nil {
		return fmt.Errorf("error while getting seasons index: %w", err)
	}

	seasonInfos := make(map[uint]*wowgd.MythicKeystoneSeason)
	seasons := make(map[uint]*database.MythicPlusSeason)

	for _, entry := range seasonsIndex.Seasons {
		blizzID := uint(entry.ID)
		seasonInfo, err := api.GetMPlusSeasonInfo(entry.ID)

		if err != nil {
			return fmt.Errorf("error while getgting information for season %v: %w", blizzID, err)
		}

		seasonInfos[blizzID] = seasonInfo

		var endTime sql.NullTime
		if endTimestamp := seasonInfo.EndTimestamp; endTimestamp != nil {
			endTime = sql.NullTime{
				Time:  time.UnixMilli(*endTimestamp),
				Valid: true,
			}
		} else {
			endTime = sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			}
		}

		seasons[blizzID] = &database.MythicPlusSeason{
			HasBlizzID: database.HasBlizzID{BlizzID: blizzID},
			HasRegion:  database.HasRegion{Region: &api.Region},
			StartTime:  time.UnixMilli(seasonInfo.StartTimestamp),
			EndTime:    endTime,
		}
	}

	err = HandleDBUpdatesGenericLogs(
		db, &api.Region, &database.MythicPlusSeason{}, []string{}, &seasons,
		logs.LogType.MythicPlusSeasonAdded,
		logs.LogType.MythicPlusSeasonUpdated,
		logs.LogType.MythicPlusSeasonRemoved,
	)

	if err != nil {
		return fmt.Errorf("error while updating seasons for region %v: %w", api.Region.Slug, err)
	}

	for blizzID, season := range seasons {
		seasonInfo := seasonInfos[blizzID]

		var existingPeriods []database.MythicPlusPeriod
		db.Gorm.Model(&database.MythicPlusPeriod{}).
			Where(&database.MythicPlusPeriod{MythicPlusSeasonID: season.ID}).
			Find(&existingPeriods)

		existingPeriodIds := make(map[uint]struct{})

		for _, existingPeriod := range existingPeriods {
			existingPeriodIds[existingPeriod.HasBlizzID.BlizzID] = struct{}{}
		}

		var newPeriods []*database.MythicPlusPeriod
		for _, period := range seasonInfo.Periods {
			periodBlizzID := uint(period.ID)
			if _, ok := existingPeriodIds[periodBlizzID]; ok {
				continue
			}

			periodInfo, err := api.GetMPlusPeriodInfo(period.ID)

			if err != nil {
				log.Warn().
					Err(err).
					Msgf(
						"error while fetching period information of period %v season %v region %v",
						periodBlizzID, blizzID, api.Region.Slug,
					)

				continue
			}

			newPeriods = append(newPeriods, &database.MythicPlusPeriod{
				HasBlizzID:       database.HasBlizzID{BlizzID: periodBlizzID},
				MythicPlusSeason: season,
				StartTime:        time.UnixMilli(periodInfo.StartTimestamp),
				EndTime:          time.UnixMilli(periodInfo.EndTimestamp),
			})
		}

		if len(newPeriods) > 0 {
			db.Gorm.Create(newPeriods)
		}

		var addedLogs []database.Log[logs.GenericAdded]
		for _, newPeriod := range newPeriods {
			addedLogs = append(addedLogs, logs.Create(&api.Region, logs.LogType.MythicPlusPeriodAdded, logs.GenericAdded{
				ID: newPeriod.GetID(),
			}))
		}

		if len(addedLogs) > 0 {
			db.Gorm.Create(addedLogs)
		}
	}

	currentSeasonBlizzID := uint(seasonsIndex.CurrentSeason.ID)
	activeSeason := api.Region.ActiveMythicPlusSeason

	if activeSeason == nil || activeSeason.GetBlizzID() != currentSeasonBlizzID {
		var oldId *uint
		if activeSeason != nil {
			oldId = &activeSeason.ID
		}
		api.Region.ActiveMythicPlusSeasonID = &seasons[currentSeasonBlizzID].ID
		db.Gorm.Save(api.Region)

		logEntry := logs.Create(&api.Region, logs.LogType.MythicPlusActiveSeasonChanged, logs.MythicPlusActiveSeasonChanged{
			OldSeasonID: oldId,
			NewSeasonID: &seasons[currentSeasonBlizzID].ID,
		})

		db.Gorm.Create(&logEntry)
	}

	return nil
}
