package crawler

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"mythic-plus-crawler/internal/blizzapi"
	"mythic-plus-crawler/internal/database"
	"mythic-plus-crawler/internal/logs"
)

func CrawlMythicPlusDungeons(api *blizzapi.BlizzApi, db *database.Database) error {
	log.Info().Msgf("crawling mythic plus dungeons for region %v", api.Region.Slug)

	dungeonsIndexes, err := api.GetMPlusDungeonIndex()

	if err != nil {
		return fmt.Errorf("error while getting dungeon index: %w", err)
	}

	dungeons := make(map[uint]*database.MythicPlusDungeon)

	for _, entry := range dungeonsIndexes.Dungeons {
		blizzID := uint(entry.ID)
		dungeonInfo, err := api.GetMPlusDungenInfo(entry.ID)

		if err != nil {
			return fmt.Errorf("error while getting information for dungeon %v: %w", blizzID, err)
		}

		upgradeTimings := make(map[uint8]int)
		upgradeTimings[1] = -1
		upgradeTimings[2] = -1
		upgradeTimings[3] = -1
		for _, timing := range dungeonInfo.KeystoneUpgrades {
			upgradeTimings[uint8(timing.UpgradeLevel)] = timing.QualifyingDuration
		}

		dungeons[blizzID] = &database.MythicPlusDungeon{
			HasBlizzID:          database.HasBlizzID{BlizzID: blizzID},
			HasRegion:           database.HasRegion{Region: &api.Region},
			Name:                dungeonInfo.Name,
			MapID:               uint(dungeonInfo.Map.ID),
			MapName:             dungeonInfo.Map.Name,
			ZoneSlug:            dungeonInfo.Zone.Slug,
			DungeonID:           uint(dungeonInfo.Dungeon.ID),
			DungeonName:         dungeonInfo.Dungeon.Name,
			QualifierTime:       upgradeTimings[1],
			QualifierDoubleTime: upgradeTimings[2],
			QualifierTripleTime: upgradeTimings[3],
		}
	}

	err = HandleDBUpdatesGenericLogs(
		db, &api.Region, &database.MythicPlusDungeon{}, []string{}, &dungeons,
		logs.LogType.MythicPlusDungeonAdded,
		logs.LogType.MythicPlusDungeonUpdated,
		logs.LogType.MythicPlusDungeonRemoved,
	)

	if err != nil {
		return fmt.Errorf("error while handling dungeon updates: %w", err)
	}

	return nil
}
