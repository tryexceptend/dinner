package models

type FootCategory int

const (
	Soup FootCategory = iota + 1
	Salad
	Meat
	SideDish
)

type Food struct {
	Name     string
	Category FootCategory
}
