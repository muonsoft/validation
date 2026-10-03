package scenario

type Input struct {
	Title   string
	Entries []Entry
}
type Entry struct {
	Code     string
	Quantity int
}
