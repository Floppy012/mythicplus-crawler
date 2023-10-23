package blizzapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"mythic-plus-crawler/internal/config"
	"mythic-plus-crawler/internal/database"
	"mythic-plus-crawler/internal/utils"

	"github.com/FuzzyStatic/blizzard/v3"
	"github.com/FuzzyStatic/blizzard/v3/wowgd"
	"github.com/FuzzyStatic/blizzard/v3/wowsearch"
)

type CountingTransport struct {
	roundTripper http.RoundTripper
	count        uint32
	mu           sync.Mutex
}

func (t *CountingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.count++
	return t.roundTripper.RoundTrip(req)
}

func (t *CountingTransport) GetCount() uint32 {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.count
}

type BlizzApi struct {
	client    *blizzard.Client
	config    *config.Config
	transport *CountingTransport
	Region    database.Region
}

func Create(config *config.Config, region database.Region) (*BlizzApi, error) {
	blizzRegion, err := utils.BlizzardRegionFromString(region.Slug)
	if err != nil {
		return nil, err
	}
	var locale blizzard.Locale
	switch blizzRegion {
	case blizzard.US:
		locale = blizzard.EnUS
	case blizzard.EU:
		locale = blizzard.EnGB
	case blizzard.KR:
		locale = blizzard.KoKR
	case blizzard.TW:
		locale = blizzard.ZhTW
	case blizzard.CN:
		locale = blizzard.ZhCN
	default:
		return nil, errors.New("unknown region")
	}

	transport := &CountingTransport{
		roundTripper: http.DefaultTransport,
	}

	client, err := blizzard.NewClient(blizzard.Config{
		ClientID:     config.BlizzardAPI.ClientID,
		ClientSecret: config.BlizzardAPI.ClientSecret,
		HTTPClient: &http.Client{
			Transport: transport,
		},
		Region: blizzRegion,
		Locale: locale,
		Retries: blizzard.RetriesConfig{
			Enabled:  true,
			Attempts: 5,
			Delay:    500 * time.Millisecond,
		},
		RateLimit: blizzard.BlizzardRateLimit,
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
		transport,
		region,
	}, nil
}

func (api *BlizzApi) makeContext() (context.Context, context.CancelFunc) {
	return context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Duration(api.config.BlizzardAPI.RequestTimeout)*time.Millisecond),
	)
}

func (api *BlizzApi) GetRequestCount() uint32 {
	return api.transport.GetCount()
}

func (api *BlizzApi) GetRealms() (*[]*wowgd.ConnectedRealmsSearch, error) {

	var out []*wowgd.ConnectedRealmsSearch

	for page := 1; ; page++ {
		ctx, cancel := api.makeContext()
		result, _, err := api.client.WoWConnectedRealmSearch(ctx, wowsearch.Page(page))
		cancel()

		if err != nil {
			return nil, fmt.Errorf("error while fetching affixes: %w", err)
		}

		out = append(out, result)

		if page == result.PageCount {
			break
		}
	}

	return &out, nil
}

func (api *BlizzApi) GetMPlusAffixesIndex() (*wowgd.MythicKeystoneAffixIndex, error) {
	ctx, cancel := api.makeContext()
	result, _, err := api.client.WoWMythicKeystoneAffixIndex(ctx)
	cancel()

	return result, err
}

func (api *BlizzApi) GetMPlusAffixInfo(affixId int) (*wowgd.MythicKeystoneAffix, error) {
	ctx, cancel := api.makeContext()
	result, _, err := api.client.WoWMythicKeystoneAffix(ctx, affixId)
	cancel()

	return result, err
}

func (api *BlizzApi) GetMPlusDungeonIndex() (*wowgd.MythicKeystoneDungeonIndex, error) {
	ctx, cancel := api.makeContext()
	result, _, err := api.client.WoWMythicKeystoneDungeonIndex(ctx)
	cancel()

	return result, err
}

func (api *BlizzApi) GetMPlusDungenInfo(dungeonId int) (*wowgd.MythicKeystoneDungeon, error) {
	ctx, cancel := api.makeContext()
	result, _, err := api.client.WoWMythicKeystoneDungeon(ctx, dungeonId)
	cancel()

	return result, err
}
