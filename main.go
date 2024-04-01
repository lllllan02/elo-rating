package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	users := InitialUsers(8000)

	users.sortByCapabilityDesc()
	print(users, "contest/cap.json")

	for i := 0; i < 50; i++ {
		fmt.Printf("mook contest: %d\n", i+1)
		users.Shuffle()
		size := randomContestantsSzie(500, 2500)

		m := make(map[int]*Contestant, size)
		contestants := make(Contestants, 0, size)
		for j := 0; j < size; j++ {
			m[j] = users[j].mook()
			contestants = append(contestants, m[j])
		}
		contestants.Process()

		for j := 0; j < size; j++ {
			contestant := m[j]
			users[j].Rating = m[j].FinalRating
			users[j].Contests = append(users[j].Contests, &UserContest{
				Id:           i + 1,
				Points:       contestant.Points,
				Rank:         contestant.Rank,
				BeforeRating: contestant.Rating,
				AfterRating:  contestant.FinalRating,
				Delta:        contestant.Delta,
			})
		}

		contestants.sortByPointsDesc()
		print(Contest{Id: i + 1, Size: len(contestants), Contestants: contestants}, fmt.Sprintf("contest/%d/%d.json", i/100+1, i+1))
	}

	users.sortByRatingDesc()
	print(users, "contest/users.json")
}

func print(data any, filename string) {
	if dir, _ := filepath.Split(filename); dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		panic(err)
	}
}
