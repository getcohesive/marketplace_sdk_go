package request

type BaseParams struct {
	IdempotencyKey string `json:"idempotency_key"`
}

