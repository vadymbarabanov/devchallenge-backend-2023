package database

import (
	"database/sql"
	"dev-challenge/internal/cell"
	"errors"
	"log"
)

type CellRepo struct {
	db *sql.DB
}

func NewCellRepository(db *sql.DB) *CellRepo {
	return &CellRepo{
		db: db,
	}
}

func (cr *CellRepo) CreateTableIfNotExists() {
	_, err := cr.db.Exec("create table if not exists sheetcell (sheet_id text not null, cell_id text not null, value text not null, result text)")
	if err != nil {
		log.Println(err)
	}
}

func (cr *CellRepo) GetOne(sheetID, cellID string) (cell.Cell, error) {
	c := cell.Cell{
		CellID:  cellID,
		SheetID: sheetID,
	}

	query := "select value, result from sheetcell where sheet_id = $1 and cell_id = $2"
	if err := cr.db.QueryRow(query, sheetID, cellID).Scan(&c.Value, &c.Result); err != nil {
		return cell.Cell{}, cell.ErrNotFound
	}

	return c, nil
}

func (cr *CellRepo) GetManyBySheetID(sheetID string) ([]cell.Cell, error) {
	query := "select cell_id, value, result from sheetcell where sheet_id = $1"
	rows, err := cr.db.Query(query, sheetID)
	if err != nil {
		return nil, err
	}

	cells := make([]cell.Cell, 0)
	for rows.Next() {
		c := cell.Cell{
			SheetID: sheetID,
		}

		if err := rows.Scan(&c.CellID, &c.Value, &c.Result); err != nil {
			log.Println(err)
		}

		cells = append(cells, c)
	}

	return cells, nil
}

func (cr *CellRepo) Insert(c cell.Cell) error {
	if c.CellID == "" || c.SheetID == "" || c.Value == "" {
		return errors.New("insertion error: invalid cell")
	}

	query := "insert into sheetcell (sheet_id, cell_id, value, result) values ($1, $2, $3, $4)"
	_, err := cr.db.Exec(query, c.SheetID, c.CellID, c.Value, c.Result)
	return err
}

func (cr *CellRepo) Update(c cell.Cell) error {
	if c.CellID == "" || c.SheetID == "" || c.Value == "" {
		return errors.New("insertion error: invalid cell")
	}

	query := "update sheetcell set value = $1, result = $2 where sheet_id = $3 and cell_id = $4"
	_, err := cr.db.Exec(query, c.Value, c.Result, c.SheetID, c.CellID)
	return err
}
