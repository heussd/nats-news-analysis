package likemodel

type ClassificationResult struct {
	Label       string  `json:"label"`
	Probability float64 `json:"probability"`
	Certain     bool    `json:"certain"`
}
