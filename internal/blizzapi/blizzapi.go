package blizzapi

import (
	"context"
	"fmt"
	"mythic-plus-crawler/internal/config"
	"net/http"
	"time"

	"github.com/FuzzyStatic/blizzard/v3"
	"github.com/FuzzyStatic/blizzard/v3/wowgd"
)

type BlizzApi struct {
	client *blizzard.Client
	config *config.Config
}

func Create(config *config.Config) (*BlizzApi, error) {
	client, err := blizzard.NewClient(blizzard.Config{
		ClientID:     config.BlizzardAPI.ClientID,
		ClientSecret: config.BlizzardAPI.ClientSecret,
		HTTPClient:   http.DefaultClient,
		Region:       blizzard.EU,
		Locale:       blizzard.EnGB,
	})

	if err != nil {
		return nil, fmt.Errorf("error while creating new blizzard api client: %w", err)
	}

	deadline := time.Now().Add(time.Duration(config.BlizzardAPI.RequestTimeout) * time.Millisecond)
	ctx, cancelCtx := context.WithDeadline(context.Background(), deadline)
	err = client.AccessTokenRequest(ctx)
	cancelCtx()
	if err != nil {
		return nil, fmt.Errorf("error while requesting access token from blizzard: %w", err)
	}

	return &BlizzApi{
		client,
		config,
	}, nil
}

func (api *BlizzApi) makeContext() (context.Context, context.CancelFunc) {
	return context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Duration(api.config.BlizzardAPI.RequestTimeout)*time.Millisecond),
	)
}

func (api *BlizzApi) GetRealms() (*wowgd.ConnectedRealmsSearch, error) {
	ctx, cancel := api.makeContext()

	realms, _, err := api.client.WoWConnectedRealmSearch(ctx)
	cancel()
	if err != nil {
		return nil, fmt.Errorf("error while fetching affixes: %w", err)
	}

	return realms, nil
}
