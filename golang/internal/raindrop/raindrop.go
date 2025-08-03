package raindrop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/heussd/nats-news-analysis/pkg/utils"
)

var (
	RaindropCollection      = utils.GetEnv("RAINDROP_COLLECTION", "")
	RaindropCollectionId, _ = strconv.Atoi(RaindropCollection)
	RaindropAccessToken     = utils.GetEnv("RAINDROP_ACCESS_TOKEN", "")
)

func init() {
	utils.RequireNotEmpty(RaindropCollection, "RAINDROP_COLLECTION")
	utils.RequireNotEmpty(RaindropAccessToken, "RAINDROP_ACCESS_TOKEN")
}

func Add(url string) (err error) {

	payload := postPayload{
		PleaseParse: pleaseParse{},
		Collection: collection{
			Id: RaindropCollectionId,
		},
		Link: url,
	}

	var jsonBytes []byte
	if jsonBytes, err = json.Marshal(payload); err != nil {
		return fmt.Errorf("cannot unmarshall: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.raindrop.io/rest/v1/raindrop?=",
		bytes.NewBuffer(jsonBytes))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+RaindropAccessToken)

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return fmt.Errorf("error doing request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from raindrop: %s", res.Status)
	}

	return nil
}
