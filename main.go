package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"mythic-plus-crawler/internal/blizzapi"
	"mythic-plus-crawler/internal/config"
	"mythic-plus-crawler/internal/crawler"
	"mythic-plus-crawler/internal/database"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	db, err := database.Connect(cfg)

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	regionApiClients := make(map[string]*blizzapi.BlizzApi)

	var activeRegions []database.Region
	db.Gorm.Where(&database.Region{Active: true}).Find(&activeRegions)

	for _, region := range activeRegions {
		log.Info().Msgf("creating api for region %v", region.Slug)
		api, err := blizzapi.Create(cfg, region)
		if err != nil {
			log.Fatal().
				Stack().
				Err(fmt.Errorf("error while creating blizzard api client for region %v: %w", region.Slug, err)).
				Send()
		}

		regionApiClients[region.Slug] = api
	}

	//log.Info().Msg("crawling realms ...")
	//for regionSlug, api := range regionApiClients {
	//	err = crawler.CrawlRealms(api, db)
	//
	//	if err != nil {
	//		log.Fatal().
	//			Stack().
	//			Err(fmt.Errorf("error while crawling realms fro region %v: %w", regionSlug, err)).
	//			Send()
	//	}
	//
	//	log.Info().Msgf("made %v requests to %v api", api.GetRequestCount(), regionSlug)
	//}

	//log.Info().Msg("crawling mythic plus affixes ...")
	//for regionSlug, api := range regionApiClients {
	//	err = crawler.CrawlMythicPlusAffixes(api, db)
	//
	//	if err != nil {
	//		log.Fatal().
	//			Stack().
	//			Err(fmt.Errorf("error while crawling mythic plus affixes from region %v: %w", regionSlug, err)).
	//			Send()
	//	}
	//
	//	log.Info().Msgf("made %v requests to %v api", api.GetRequestCount(), regionSlug)
	//}

	log.Info().Msg("crawling mythic plus dungeons ...")
	for regionSlug, api := range regionApiClients {
		err = crawler.CrawlMythicPlusDungeons(api, db)

		if err != nil {
			log.Fatal().
				Stack().
				Err(fmt.Errorf("error while crawling mythic plus dungeons from region %v: %w", regionSlug, err)).
				Send()
		}

		log.Info().Msgf("made %v requests to %v api", api.GetRequestCount(), regionSlug)
	}
}
