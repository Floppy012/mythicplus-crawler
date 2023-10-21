package main

import (
	"fmt"
	"log"
	"mythic-plus-crawler/internal/blizzapi"
	"mythic-plus-crawler/internal/config"
	"mythic-plus-crawler/internal/database"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg)

	if err != nil {
		log.Fatal(err)
	}

	api, err := blizzapi.Create(cfg)

	if err != nil {
		log.Fatal(err)
	}

	realms, err := api.GetRealms()

	if err != nil {
		log.Fatal(err)
	}

	for _, connected := range realms.Results {
		for _, realm := range connected.Data.Realms {
			fmt.Printf(realm.Region.Name.EnUS)
			db.UpsertRegion(realm.Region.ID, realm.Region.Name.EnUS)
		}
	}
}
