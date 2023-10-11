package cell

import "errors"

var (
	ErrNotFound = errors.New("entity not found")
)

type Repository interface {
	GetOne(sheetID, cellID string) (Cell, error)
	GetManyBySheetID(sheetID string) ([]Cell, error)
	Insert(cell Cell) error
	Update(cell Cell) error
}
