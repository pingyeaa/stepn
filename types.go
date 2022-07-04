package main

type OrderList struct {
	Code int     `json:"code"`
	Data []*Shoe `json:"data"`
}

type Shoe struct {
	ID        int    `json:"id"`
	TypeID    int    `json:"type_id"`
	TypeName  string `json:"type_name"`
	Color     string `json:"color"`
	Otd       int    `json:"otd"`
	Time      int    `json:"time"`
	PropID    int    `json:"propID"`
	Img       string `json:"img"`
	DataID    int    `json:"dataID"`
	SellPrice int    `json:"sellPrice"`
	Hp        int    `json:"hp"`
	Level     int    `json:"level"`
	Quantity  int    `json:"quality"`
	Mint      int    `json:"mint"`
	AddRatio  int    `json:"addRatio"`
	V1        int    `json:"v1"`
	V2        int    `json:"v2"`
}
