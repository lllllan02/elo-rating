package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"time"

	"github.com/duke-git/lancet/v2/mathutil"
)

func main() {
	contest := 1946
	now := time.Now()
	defer func() {
		fmt.Printf("time.Since(now): %v\n", time.Since(now))
	}()

	// read rating json
	data, err := os.ReadFile(fmt.Sprintf("%d.json", contest))
	if err != nil {
		panic(err)
	}

	var contestants []*Contestant
	if err = json.Unmarshal(data, &contestants); err != nil {
		panic(err)
	}

	// calculate seed
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

	// 赛前 rating 从大到小
	sort.Slice(contestants, func(i, j int) bool { return contestants[i].Rating > contestants[j].Rating })

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
	sum = 0
	size := float64(len(contestants))
	zeroSumCount := mathutil.Min(int(4*math.Round(math.Sqrt(size))), int(size))
	for i := 0; i < zeroSumCount; i++ {
		sum += contestants[i].Delta
	}
	inc = mathutil.Min(mathutil.Max(-sum/zeroSumCount, -10), 0)
	for _, c := range contestants {
		c.Delta += inc
	}

	validateDeltas(contestants)

	// write result
	file, err := os.Create(fmt.Sprintf("%d.res", contest))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for _, c := range contestants {
		bytes, _ := json.Marshal(c)
		file.WriteString(string(bytes))
		file.WriteString("\n")
	}
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

func validateDeltas(contestants []*Contestant) {
	sort.Slice(contestants, func(i, j int) bool { return contestants[i].Rank < contestants[j].Rank })

	for i := 0; i < len(contestants); i++ {
		for j := i + 1; j < len(contestants); j++ {
			if contestants[i].Rating > contestants[j].Rating {
				if contestants[i].Rating+contestants[i].Delta < contestants[j].Rating+contestants[j].Delta {
					fmt.Printf("First rating invariant fialed: (%v->%v) vs. (%v->%v)\n",
						contestants[i].Rating,
						contestants[i].Rating+contestants[i].Delta,
						contestants[j].Rating,
						contestants[j].Rating+contestants[j].Delta,
					)
				}
			}

			if contestants[i].Rating < contestants[j].Rating {
				if contestants[i].Delta < contestants[j].Delta {
					fmt.Printf("Second rating invariant fialed: %v vs. %v\n",
						contestants[i].Delta,
						contestants[j].Delta,
					)
				}
			}
		}
	}
}
