package crawler

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"mythic-plus-crawler/internal/blizzapi"
	"mythic-plus-crawler/internal/database"
	"mythic-plus-crawler/internal/logs"
)

func CrawlMythicPlusAffixes(api *blizzapi.BlizzApi, db *database.Database) error {
	log.Info().Msgf("crawling mythic plus affixes for region %v", api.Region.Slug)

	affixesIndex, err := api.GetMPlusAffixesIndex()

	model := &database.MythicPlusAffix{}

	if err != nil {
		return fmt.Errorf("error while getting affixes index: %w", err)
	}

	affixes := make(map[uint]*database.MythicPlusAffix)

	for _, entry := range affixesIndex.Affixes {
		blizzID := uint(entry.ID)
		affixInfo, err := api.GetMPlusAffixInfo(entry.ID)

		if err != nil {
			return fmt.Errorf("error while getting information for affix %v: %w", blizzID, err)
		}

		affixes[blizzID] = &database.MythicPlusAffix{
			HasBlizzID: database.HasBlizzID{
				BlizzID: blizzID,
			},
			HasRegion:   database.HasRegion{Region: &api.Region},
			Name:        entry.Name,
			Description: affixInfo.Description,
		}
	}

	err = HandleDBUpdatesGenericLogs(
		db, &api.Region, model, []string{}, &affixes,
		logs.LogType.MythicPlusAffixAdded,
		logs.LogType.MythicPlusAffixUpdated,
		logs.LogType.MythicPlusAffixRemoved,
	)

	if err != nil {
		return fmt.Errorf("error while handling affix updates: %w", err)
	}

	return nil
}
