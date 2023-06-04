package response

type Response struct {
	Data  any   `json:"data"`
	Error Error `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
