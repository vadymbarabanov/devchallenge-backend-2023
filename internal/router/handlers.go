package router

import (
	"dev-challenge/internal/cell"
	"encoding/json"
	"net/http"
	"strings"
)

func (rt *Router) establishRoutes() {
	// /api/v1/:sheet_id
	rt.Get(`^\/api\/v1\/(?P<sheet_id>[\w-]+)$`, rt.handleGetSheet)

	// /api/v1/:sheet_id/:cell_id
	rt.Get(`^\/api\/v1\/(?P<sheet_id>[\w-]+)\/(?P<cell_id>[\w-]+)$`, rt.handleGetCell)

	// /api/v1/:sheet_id/:cell_id
	rt.Post(`^\/api\/v1\/(?P<sheet_id>[\w-]+)\/(?P<cell_id>[\w-]+)$`, rt.handlePostCell)
}

func (rt *Router) handleGetCell(ctx *Ctx) {
	sheetID, okSheetID := ctx.Params["sheet_id"]
	cellID, okCellID := ctx.Params["cell_id"]

	if !okSheetID || !okCellID {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	sheetID = strings.ToLower(sheetID)
	cellID = strings.ToLower(cellID)

	cell, err := rt.cellService.GetCell(sheetID, cellID)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		ctx.Response.Write([]byte("Cell " + http.StatusText(http.StatusNotFound)))
		return
	}

	respondJSON(ctx.Response, &cell)
}

func (rt *Router) handleGetSheet(ctx *Ctx) {
	sheetID, okSheetID := ctx.Params["sheet_id"]

	if !okSheetID {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	sheetID = strings.ToLower(sheetID)

	sheet, err := rt.sheetService.GetSheet(sheetID)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		ctx.Response.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	if len(sheet) == 0 {
		ctx.Response.WriteHeader(http.StatusNotFound)
		ctx.Response.Write([]byte("Sheet " + http.StatusText(http.StatusNotFound)))
		return
	}

	respondJSON(ctx.Response, &sheet)
}

func (rt *Router) handlePostCell(ctx *Ctx) {
	sheetID, okSheetID := ctx.Params["sheet_id"]
	cellID, okCellID := ctx.Params["cell_id"]

	if !okSheetID || !okCellID {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	c := cell.Cell{
		CellID:  cellID,
		SheetID: sheetID,
	}

	if err := json.NewDecoder(ctx.Request.Body).Decode(&c); err != nil {
		ctx.Response.WriteHeader(http.StatusUnprocessableEntity)
		ctx.Response.Write([]byte("cannot process request body"))
		return
	}

	if strings.Trim(c.Value, " ") == "" {
		ctx.Response.WriteHeader(http.StatusUnprocessableEntity)
		ctx.Response.Write([]byte("value is required"))
		return
	}

	result, err := rt.cellService.UpsertCell(c)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusUnprocessableEntity)
		respondJSON(ctx.Response, map[string]string{
			"message": err.Error(),
			"value":   c.Value,
			"result":  "ERROR",
		})
		return
	}

	ctx.Response.WriteHeader(http.StatusCreated)
	respondJSON(ctx.Response, &result)
}
