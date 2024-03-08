package general

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type EdgeRule struct {
	ID                  string
	ActionType          string
	ActionParameter1    string
	ActionParameter2    string
	Trigger             []EdgeRuleTrigger
	TriggerMatchingType int
	Description         string
	Enabled             bool
}

type EdgeRuleTrigger struct {
	Type                int
	PatternMatches      []string
	PatternMatchingType int
	Parameter1          string
}

func (c *bunnyClient) CreateEdgeRule(ctx context.Context, zoneID string, rule EdgeRule) (*EdgeRule, error) {
	apiURL := c.baseAPIUrl + pullZonePath + "/" + zoneID + "/edgerules/addOrUpdate"

	if rule.ID == "" {
		apiURL += "/" + rule.ID
		rule.ID = ""
	}

	buf := make([]byte, 0, defaultBufferSize)
	if err := json.NewEncoder(bytes.NewBuffer(buf)).Encode(rule); err != nil {
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

	edgeRule := EdgeRule{}
	if err := json.NewDecoder(res.Body).Decode(&edgeRule); err != nil {
		return nil, err
	}

	return &edgeRule, nil
}
