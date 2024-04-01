package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	contest := 592

	data, err := os.ReadFile(fmt.Sprintf("%d.json", contest))
	if err != nil {
		log.Panic(err)
	}

	var contestants Contestants
	if err := json.Unmarshal(data, &contestants); err != nil {
		panic(err)
	}

	contestants.Process()

	fmt.Printf("%12s, %12s, %12s, %12s, %12s, %12s, %12s, %12s\n", "i", "rank", "rating", "n_rating", "seed", "delta", "f_rating", "d_rating")
	for i, contestant := range contestants {
		fmt.Printf("%12d, %12d, %12d, %12d, %12f, %12d, %12d, %12d\n", i, contestant.Rank, contestant.Rating, contestant.NeedRating, contestant.Seed, contestant.Delta, contestant.FinalRating, contestant.Rating+contestant.Delta-contestant.FinalRating)
	}
}
