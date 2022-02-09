package report

type Report struct {
	Abundance []int  `json:"abundance"`
	Gender    Gender `json:"gender"`
	Case      []Case `json:"case"`
}
type Case struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}
type Gender struct {
	Male   int `json:"male"`
	Female int `json:"female"`
}
