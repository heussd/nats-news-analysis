package likemodel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/heussd/nats-news-analysis/pkg/utils"
)

var (
	LikeModelServer      = utils.GetEnv("LIKE_MODEL_SERVER", "http://like-model-server:8000")
	ProbabilityThreshold = utils.GetEnvFloat("PROBABILITY_THRESHOLD", "0.95")
)

func init() {
	fmt.Printf("Using %s...\n", LikeModelServer)
	for status := 0; status != 200; {
		resp, _ := http.Get(LikeModelServer)
		if resp == nil {
			fmt.Printf("Waiting for %s to come up...\n", LikeModelServer)
			time.Sleep(2 * time.Second)
		} else {
			status = resp.StatusCode
		}
	}
}

func Classify(text string) (bool, error) {
	client := &http.Client{}

	jsonData := map[string]interface{}{
		"text": text,
	}
	jsonBytes, _ := json.Marshal(jsonData)

	req, err := http.NewRequest("POST", LikeModelServer+"/classify", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}

	var response *http.Response
	response, err = client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	var result ClassificationResult
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return false, err
	}

	if !result.Certain {
		return false, nil
	}

	if result.Probability < ProbabilityThreshold {
		return false, nil
	}

	return true, nil
}
