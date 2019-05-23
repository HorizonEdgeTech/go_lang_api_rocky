package model

type Rider struct {
	Surname           string    `json:"Surname"`
	OtherName         string    `json:"OtherName"`
	IDNumber          int       `json:"IDNumber"`
	MobileNumber      string    `json:"MobileNumber"`
	AlternativeNumber string    `json:"AlternativeNumber"`
	Nationality string `json:"Nationality"`
	NoOfDependents int `json:"NoOfDependents"`
}
