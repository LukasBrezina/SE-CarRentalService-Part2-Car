package types

type Car struct {
	ID        int     `json:"id"`
	Model     string  `json:"model"`
	Brand     string  `json:"brand"`
	CollectAt string  `json:"collectAt"`
	AccountId int     `json:"accountId"`
	Year      int     `json:"year"`
	Price     float32 `json:"price"`
	PS        int     `json:"ps"`
}
