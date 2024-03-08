package general

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/alvadorncorp/bunny-go/internal/bunny"
)

const pullZonePath = "/pullzone"

type ListPullZoneParams struct {
	bunny.PageParams
	Search bunny.OptionalValue[string]
}

type PullZone struct {
	ID                              string `json:"Id"`
	Name                            string
	OriginUrl                       string
	Enabled                         bool
	Hostnames                       []PullZoneHostname
	StorageZoneID                   int
	EdgeScriptID                    int
	AllowedReferrer                 []string
	BlockedReferrer                 []string
	BlockedIP                       []string
	EnabledGeoZoneUS                bool
	EnabledGeoZoneEU                bool
	EnabledGeoZoneAsia              bool
	EnabledGeoZoneSA                bool
	ZoneSecurityEnabled             bool
	ZoneSecurityKey                 string
	ZoneSecurityIncludeHashRemoteIP bool
	IgnoreQueryStrings              bool
	MonthlyBandwidthLimit           int
	MonthlyBandwidthUsed            int
}

type PullZoneHostname struct{}

func (c *bunnyClient) ListPullZones(ctx context.Context, params ListPullZoneParams) (bunny.Page[PullZone], error) {
	apiURL := c.baseAPIUrl + pullZonePath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)

	if err != nil {
		return bunny.Page[PullZone]{}, err
	}

	qp := url.Values{
		"page":    []string{fmt.Sprint(params.Page.ValueOrDefault(bunny.DefaultPage))},
		"perPage": []string{fmt.Sprint(params.PerPage.ValueOrDefault(bunny.DefaultPerPage))},
		"search":  []string{params.Search.ValueOrDefault("")},
	}
	req.URL.RawQuery = qp.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return bunny.Page[PullZone]{}, err
	}

	pullZonePage := bunny.Page[PullZone]{}
	if err := json.NewDecoder(res.Body).Decode(&pullZonePage); err != nil {
		return bunny.Page[PullZone]{}, err
	}

	return pullZonePage, nil
}

func (c *bunnyClient) CreatePullZone(ctx context.Context, zone PullZone) (*PullZone, error) {
	apiURL := c.baseAPIUrl + "/pullzone"

	if zone.ID == "" {
		apiURL += "/" + zone.ID
		zone.ID = ""
	}

	buf := make([]byte, 0, defaultBufferSize)
	if err := json.NewEncoder(bytes.NewBuffer(buf)).Encode(zone); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(buf))

	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	pullZone := PullZone{}
	if err := json.NewDecoder(res.Body).Decode(&pullZone); err != nil {
		return nil, err
	}

	return &pullZone, nil
}

func (c *bunnyClient) GetPullZone(ctx context.Context, id string) (*PullZone, error) {
	apiURL := c.baseAPIUrl + pullZonePath + id
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)

	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	pullZone := PullZone{}
	if err := json.NewDecoder(res.Body).Decode(&pullZone); err != nil {
		return nil, err
	}

	return &pullZone, nil
}

func (c *bunnyClient) DeletePullZone(ctx context.Context, id string) (*PullZone, error) {
	apiURL := c.baseAPIUrl + pullZonePath + id
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, apiURL, nil)

	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	pullZone := PullZone{}
	if err := json.NewDecoder(res.Body).Decode(&pullZone); err != nil {
		return nil, err
	}

	return &pullZone, nil
}
