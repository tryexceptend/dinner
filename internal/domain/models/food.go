package models

type FootCategory int

// Типы еды
const (
	Soup FootCategory = iota + 1
	Salad
	Meat
	SideDish
)

// Описание одного блюда
type Food struct {
	Name     string
	Category FootCategory
}
