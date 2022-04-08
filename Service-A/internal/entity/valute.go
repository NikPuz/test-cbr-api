package entity

type Valute struct {
	NumCode  int     `xml:"NumCode", json:"num_code"`
	CharCode string  `xml:"CharCode", json:"char_code"`
	Nominal  int     `xml:"Nominal", json:"nominal"`
	Name     string  `xml:"Name", json:"name"`
	Value    float64 `xml:"Value", json:"value"`
}

type ValCurs struct {
	Date   string   `xml:"Date,attr", json:"date"`
	Valute []Valute `xml:"Valute", json:"valute"`
}