package sheet

import (
	"dev-challenge/internal/cell"
)

type Service struct {
	cellService *cell.Service
}

func NewService(cellService *cell.Service) *Service {
	return &Service{
		cellService: cellService,
	}
}

func (s *Service) GetSheet(sheetID string) (Sheet, error) {
	cells, err := s.cellService.GetCellsBySheetID(sheetID)
	if err != nil {
		return Sheet{}, err
	}

	return Sheet(cells), nil
}
