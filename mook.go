package main

import (
	"math/rand"
	"slices"

	"gonum.org/v1/gonum/stat/distuv"
)

type Contest struct {
	Id          int         `json:"id,omitempty"`
	Size        int         `json:"size,omitempty"`
	Contestants Contestants `json:"contestants,omitempty"`
}

type UserContest struct {
	Id           int `json:"id,omitempty"`
	Points       int `json:"points,omitempty"`
	Rank         int `json:"rank,omitempty"`
	BeforeRating int `json:"before_rating,omitempty"`
	AfterRating  int `json:"after_rating,omitempty"`
	Delta        int `json:"delta,omitempty"`
}

type User struct {
	Id         int            `json:"id,omitempty"`
	Capability int            `json:"capability,omitempty"`
	Rating     int            `json:"rating,omitempty"`
	Contests   []*UserContest `json:"contests,omitempty"`
}

func (user *User) mook() *Contestant {
	return &Contestant{
		Rating: user.Rating,
		Points: randomPoints(user.Capability),
	}
}

type Users []*User

func (users Users) sortByRatingDesc() {
	slices.SortFunc(users, func(a, b *User) int { return b.Rating - a.Rating })
}

func (users Users) sortByCapabilityDesc() {
	slices.SortFunc(users, func(a, b *User) int { return b.Capability - a.Capability })
}

func (users Users) Shuffle() {
	rand.Shuffle(len(users), func(i, j int) {
		users[i], users[j] = users[j], users[i]
	})
}

func InitialUsers(size int) Users {
	normal := distuv.Normal{Mu: 70, Sigma: 10}

	users := make(Users, 0, size)
	for i := 0; i < size; i++ {
		users = append(users, &User{
			Id:         i + 1,
			Capability: min(100, int(normal.Rand())),
			Rating:     1000,
		})
	}

	return users
}

func randomContestantsSzie(mins, maxs int) int {
	return rand.Intn(maxs-mins+1) + mins
}

func randomPoints(capability int) int {
	return rand.Intn((101-capability)*10) + capability*10
}
