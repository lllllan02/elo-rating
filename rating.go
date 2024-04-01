package main

import (
	"log"
	"math"
	"slices"
)

type Contestant struct {
	Rank   int `json:"rank,omitempty"`   // 真实排名
	Rating int `json:"rating,omitempty"` // 赛前分数
	Points int `json:"points,omitempty"` // 比赛得分

	Seed        float64 `json:"seed,omitempty"`         // 理论排名
	Delta       int     `json:"delta,omitempty"`        // 分数变化
	NeedRating  int     `json:"need_rating,omitempty"`  // 理论分数
	FinalRating int     `json:"final_rating,omitempty"` // 赛后分数
}

type Contestants []*Contestant

func (contestants Contestants) Process() {
	if len(contestants) == 0 {
		return
	}

	contestants.reassignRanks()

	for _, a := range contestants {
		a.Seed = 1
		for _, b := range contestants {
			if a != b {
				a.Seed += b.getEloWinProbability(a)
			}
		}
	}

	for _, contestant := range contestants {
		midRank := math.Sqrt(float64(contestant.Rank) * contestant.Seed)
		contestant.NeedRating = contestants.getRatingToRank(midRank)
		contestant.Delta = (contestant.NeedRating - contestant.Rating) / 2
	}

	contestants.sortByRatingDesc()

	{
		sum := 0
		for _, contestant := range contestants {
			sum += contestant.Delta
		}
		inc := -sum/len(contestants) - 1
		for _, contestant := range contestants {
			contestant.Delta += inc
		}
	}

	{
		sum, zeroSumCount := 0, min(int(4*math.Round(math.Sqrt(float64(len(contestants))))), len(contestants))
		for i := 0; i < zeroSumCount; i++ {
			sum += contestants[i].Delta
		}
		inc := min(max(-sum/zeroSumCount, -10), 0)
		for _, contestant := range contestants {
			contestant.Delta += inc
		}
	}

	for _, contestant := range contestants {
		contestant.FinalRating = contestant.Rating + contestant.Delta
	}

	contestants.validateDeltas()
}

func (contestants Contestants) reassignRanks() {
	contestants.sortByPointsDesc()

	for _, contestant := range contestants {
		contestant.Rank, contestant.Delta = 0, 0
	}

	first, points := 0, contestants[0].Points
	for i := 1; i < len(contestants); i++ {
		if contestants[i].Points < points {
			for j := first; j < i; j++ {
				contestants[j].Rank = i
			}
			first = i
			points = contestants[i].Points
		}
	}

	{
		rank := len(contestants)
		for j := first; j < len(contestants); j++ {
			contestants[j].Rank = rank
		}
	}
}

func (contestants Contestants) sortByRatingDesc() {
	slices.SortFunc(contestants, func(a, b *Contestant) int { return b.Rating - a.Rating })
}

func (contestants Contestants) sortByPointsDesc() {
	slices.SortFunc(contestants, func(a, b *Contestant) int { return b.Points - a.Points })
}

func (contestants Contestants) getRatingToRank(rank float64) int {
	left, right := -8000, 8000

	for right-left > 1 {
		mid := (left + right) / 2

		if contestants.getSeed(mid) < rank {
			right = mid
		} else {
			left = mid
		}
	}

	return left
}

func (contestants Contestants) getSeed(rating int) float64 {
	extraContestant := &Contestant{Rating: rating}

	result := 1.
	for _, contestant := range contestants {
		result += contestant.getEloWinProbability(extraContestant)
	}

	return result
}

func (contestants Contestants) validateDeltas() {
	contestants.sortByPointsDesc()

	for i := 0; i < len(contestants); i++ {
		for j := i + 1; j < len(contestants); j++ {
			if contestants[i].Rating > contestants[j].Rating {
				if contestants[i].FinalRating < contestants[j].FinalRating {
					log.Panicf("First rating invariant failed: %d vs. %d.", i, j)
				}
			}
			if contestants[i].Rating < contestants[j].Rating {
				if contestants[i].Delta < contestants[j].Delta {
					log.Panicf("Second rating invariant failed: %d vs. %d.", i, j)
				}
			}
		}
	}
}

func (a *Contestant) getEloWinProbability(b *Contestant) float64 {
	return getEloWinProbability(float64(a.Rating), float64(b.Rating))
}

func getEloWinProbability(ra, rb float64) float64 {
	return 1. / (1 + math.Pow(10, (rb-ra)/400.))
}
