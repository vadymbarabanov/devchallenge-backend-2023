package cell

import (
	"dev-challenge/internal/evaluator"
	"dev-challenge/internal/parser"
	"errors"
	"strconv"
)

type Service struct {
	cellRepo Repository
}

func NewService(cellRepo Repository) *Service {
	return &Service{
		cellRepo: cellRepo,
	}
}

func (s *Service) GetCell(sheetID, cellID string) (Cell, error) {
	cell, err := s.cellRepo.GetOne(sheetID, cellID)
	if err != nil {
		return Cell{}, err
	}

	return cell, nil
}

func (s *Service) GetCellsBySheetID(sheetID string) ([]Cell, error) {
	cells, err := s.cellRepo.GetManyBySheetID(sheetID)
	if err != nil {
		return nil, err
	}

	return cells, nil
}

func (s *Service) UpsertCell(c Cell) (Cell, error) {
	formulaTree, err := parser.Parse(c.Value)
	if err != nil {
		return Cell{}, err
	}

	result, err := evaluator.Evaluate(formulaTree, func(cellID string) (string, error) {
		cell, err := s.cellRepo.GetOne(c.SheetID, cellID)
		if err != nil {
			return "", err
		}
		return cell.Value, nil
	})
	if err != nil {
		return Cell{}, err
	}

	c.Result = strconv.FormatFloat(result, 'f', -1, 32)

	_, err = s.cellRepo.GetOne(c.SheetID, c.CellID)
	if errors.Is(err, ErrNotFound) {
		if err := s.cellRepo.Insert(c); err != nil {
			return Cell{}, err
		}
	} else {
		if err := s.cellRepo.Update(c); err != nil {
			return Cell{}, err
		}
	}

	return c, nil
}
