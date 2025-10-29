package http

type APIResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
	Errors     any    `json:"errors,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more,omitempty"`
}

type InfiniteScrollRequest struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit"`
}
