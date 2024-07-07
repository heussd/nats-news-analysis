package raindrop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"nats-news-analysis/internal/config"
)

type pleaseParse struct{}

type collection struct {
	Id int `json:"$id"`
}
type postPayload struct {
	PleaseParse pleaseParse `json:"pleaseParse"`
	Collection  collection  `json:"collection"`
	Link        string      `json:"link"`
}

var collectionId, _ = strconv.Atoi(config.RaindropCollection)

func Add(url string) (ok bool, err error) {

	payload := postPayload{
		PleaseParse: pleaseParse{},
		Collection: collection{
			Id: collectionId,
		},
		Link: url,
	}

	var jsonBytes []byte
	if jsonBytes, err = json.Marshal(payload); err != nil {
		return false, fmt.Errorf("cannot unmarshall: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.raindrop.io/rest/v1/raindrop?=",
		bytes.NewBuffer(jsonBytes))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.RaindropToken)

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return false, fmt.Errorf("error doing request: %w", err)
	}

	defer res.Body.Close()

	return res.StatusCode == 200, nil
}
