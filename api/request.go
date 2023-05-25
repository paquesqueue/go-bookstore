package api

type RequestBook struct {
	Title      string   `json:"title"`
	Authors    []string `json:"authors"`
	Publisher  string   `json:"publisher"`
	Isbn       string   `json:"isbn"`
	Price      int64    `json:"price"`
	Quantity   int64    `json:"quantity"`
	Created_by string   `json:"created_by"`
}

type RequestUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
}

type RequestGetAll struct {
	PageId   int64 `json:"page_id"`
	PageSize int64 `json:"page_size"`
}

type GetAllParams struct {
	Limit  int64
	Offset int64
}
