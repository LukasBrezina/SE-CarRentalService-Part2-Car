package types

type CreateCarRequest struct {
	Model     string  `json:"model" binding:"required"`
	Brand     string  `json:"brand" binding:"required"`
	CollectAt string  `json:"collectAt" binding:"required"`
	AccountId int     `json:"accountId"`
	Year      int     `json:"year" binding:"required"`
	Price     float32 `json:"price" binding:"required"`
	PS        int     `json:"ps" binding:"required"`
}
