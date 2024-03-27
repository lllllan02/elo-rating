package main

type Contestant struct {
	Rank        int `json:"rank,omitempty"`         // 真实排名
	Rating      int `json:"rating,omitempty"`       // 赛前分数
	AfterRating int `json:"after_rating,omitempty"` // 赛后分数

	Seed       float64 `json:"seed,omitempty"`        // 理论排名
	NeedRating int     `json:"need_rating,omitempty"` // 理论分数
	Delta      int     `json:"delta,omitempty"`       // 分数变化
}
