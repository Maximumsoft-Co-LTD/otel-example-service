package dto

type Response[T any] struct {
	Code   string `json:"code"`
	Msg    string `json:"message"`
	Result T      `json:"result,omitempty"`
}

type Pagination[T any] struct {
	Limit      int   `json:"limit"`
	Page       int   `json:"page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	Items      []T   `json:"items"`
}

type LoginResp struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire_at"`
}
