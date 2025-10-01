package models

type Profile struct {
	FullName    *string `json:"fullname"`
	PhoneNumber *string `json:"phone"`
	Img         *string `json:"img"`
}

type Follow struct {
	ID       int    `json:"id,omitempty"`
	Fullname string `json:"fullname"`
	Img      string `json:"img"`
}
