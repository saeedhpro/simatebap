package report


type Report struct {
	Abundance []int `json:"abundance"`
	Gender Gender `json:"gender"`
}

type Gender struct {
	Male int `json:"male"`
	Female int `json:"female"`
}