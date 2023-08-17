package alert

// Event is a model of Alert Event
type Event struct {
	Level     string `json:"level"`
	Content   string `json:"content"`
	AlertTime string `json:"alertTime"`
}
