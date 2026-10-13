package scenario

type Request struct {
	Title   string  `json:"title"`
	Contact Contact `json:"contact"`
	Visits  []Visit `json:"visits"`
}
type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type Visit struct {
	Address string `json:"address"`
	Tasks   []Task `json:"tasks"`
}
type Task struct {
	ServiceCode string `json:"serviceCode"`
	Variant     string `json:"variant"`
	Duration    int    `json:"duration"`
}
