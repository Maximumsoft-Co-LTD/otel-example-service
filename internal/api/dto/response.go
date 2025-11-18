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

type DomainResp struct {
}

type MangaResp struct {
	ID          string   `json:"id"`
	Prefix      string   `json:"prefix"`
	OrderNo     string   `json:"order_no"`
	Title       string   `json:"title"`
	Pseudonym   string   `json:"pseudonym"`
	Description string   `json:"description"`
	Tags        []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Categories  []string `bson:"category,omitempty" json:"category,omitempty"`
	Action      bool     `bson:"action" json:"action"`
	Stats       struct {
		FavoritesCount int     `bson:"favorites_count" json:"favorites_count"`
		RatingAvg      float64 `bson:"rating_avg" json:"rating_avg"`
		RatingCount    int     `bson:"rating_count" json:"rating_count"`
	} `bson:"stats" json:"stats"`
}

type EpisodeResp struct {
}

type BucketResp struct {
}

type ImageResp struct {
}

type CategoryResp struct {
}

type TagResp struct {
}
