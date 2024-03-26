package main

type Contestant struct {
	Rank       int `json:"rank,omitempty"`
	Rating     int `json:"rating,omitempty"`
	NeedRating int `json:"need_rating,omitempty"`
}
