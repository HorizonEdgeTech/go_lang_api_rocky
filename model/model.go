package model

import "time"

type Rider struct {
	Surname           string    `json:"Surname"`
	OtherName         string    `json:"OtherName"`
	DOB               time.Time `json:Dob`
	IDNumber          int       `json:"IDNumber"`
	MobileNumber      string    `json:"MobileNumber"`
	AlternativeNumber string    `json:"AlternativeNumber"`
}
