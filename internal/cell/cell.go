package cell

type Cell struct {
	CellID  string `json:"-"`
	SheetID string `json:"-"`
	Value   string `json:"value"`
	Result  string `json:"result"`
}
