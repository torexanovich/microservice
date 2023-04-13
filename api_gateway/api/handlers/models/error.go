package models

type Error struct {
	Code        int
	Error       error
	Description string
}

type SuccessInfo struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}