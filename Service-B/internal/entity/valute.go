package entity

type Valute struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Nominal  int     `json:"nominal"`
	Name     string  `json:"name"`
	Value    float64 `json:"value"`
}

type ValCurs struct {
	Date   string   `json:"date"`
	Valute []Valute `json:"valute"`
}
