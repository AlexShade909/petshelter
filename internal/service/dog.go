package service

import (
	"petshelter/internal/models"
	"sort"
)

func Nicknames(dogs map[string]models.Dog) []string {
	nicknames := make([]string, 0, len(dogs))
	for nickname := range dogs {
		nicknames = append(nicknames, nickname)
	}
	sort.Strings(nicknames)
	return nicknames
}
