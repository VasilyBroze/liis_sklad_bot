package entity

type BotSettings struct {
	Google_sheet_stock_url  string `json:"google_stock_url"`
	Google_sheet_stock_list string `json:"google_stock_list"`
	Bot_token               string `json:"bot_token"`
}
