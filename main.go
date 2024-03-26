package main

import (
	"encoding/json"
	"math"
	"os"
)

func main() {
	// read rating json
	data, err := os.ReadFile("1943.json")
	if err != nil {
		panic(err)
	}

	var contestants []*Contestant
	if err = json.Unmarshal(data, &contestants); err != nil {
		panic(err)
	}

	// calculate win seed
	for _, ci := range contestants {
		ci.Seed = 1
		for _, cj := range contestants {
			if ci != cj {
				ci.Seed += getEloWinProbability(float64(cj.Rating), float64(ci.Rating))
			}
		}
	}

	// calculate need rating
	for _, c := range contestants {
		midRank := math.Sqrt(float64(c.Rank) * c.Seed)
		c.NeedRating = getRatingToRank(midRank, contestants)
		c.Delta = (c.NeedRating - c.Rating) / 2
	}

	// total sum should not be more than zero.
	sum := 0
	for _, c := range contestants {
		sum += c.Delta
	}
	inc := -sum/len(contestants) - 1
	for _, c := range contestants {
		c.Delta += inc
	}

	// sum of top-4*sqrt should be adjusted to zero.
	// sum = 0
	// size := float64(len(contestants))
	// zeroSumCount := int(math.Min(4*math.Round(math.Sqrt(size)), size))

	// for _, c := range contestants {
	// 	fmt.Printf("%v->%v(%v)\n", c.Rating, c.AfterRating, c.Rating+c.Delta)
	// }
}

func getEloWinProbability(a, b float64) float64 {
	return 1.0 / (1.0 + math.Pow(10, (b-a)/400.0))
}

func getRatingToRank(rank float64, contestants []*Contestant) int {
	left, right := 1, 8000

	for right-left > 1 {
		mid := (left + right) / 2

		if getSeed(mid, contestants) < rank {
			right = mid
		} else {
			left = mid
		}
	}

	return left
}

func getSeed(rating int, contestants []*Contestant) float64 {
	result := float64(1)

	for _, c := range contestants {
		result += getEloWinProbability(float64(c.Rating), float64(rating))
	}

	return result
}
