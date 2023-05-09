package models

type Policy struct {
	User   string `json:"user"`
	Domain string `json:"domain"`
	Action string `json:"action"`
}

type UpdatePol struct {
	NewPolicy Policy `json:"new"`
	OldPolicy Policy `json:"old"`
}
