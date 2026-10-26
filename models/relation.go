package models

type Relation struct {
	Id          uint64
	Name        string
	Type        string
	ReverseType string
	Lag         int
	HasLag      bool
	Description string
	FromId      uint64
	FromTitle   string
	ToId        uint64
	ToTitle     string
}
