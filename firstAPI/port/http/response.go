package http

// Response
// is the response data of get handler.
type Response struct {
	Id          uint64 `json:"id"`
	Description string `json:"description"`
	Email       string `json:"email"`
	State       string `json:"state"`
	Category    string `json:"category"`
	Image       string `json:"image"`
}
