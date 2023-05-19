package entity

// СТРУКТУРА ПАРСИНГА ИЗ GOOGLE SHEETS
type Stock struct {
	StockName string `json:"Склад"`
	Number    string `json:"Артикул"`
	Units     string `json:"Единица"`
	Quantity  string `json:"Конечный остаток"`
}
