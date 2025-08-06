package raindrop

type pleaseParse struct{}

type collection struct {
	Id int `json:"$id"`
}
type postPayload struct {
	PleaseParse pleaseParse `json:"pleaseParse"`
	Collection  collection  `json:"collection"`
	Link        string      `json:"link"`
}
