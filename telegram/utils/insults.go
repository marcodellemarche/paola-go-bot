package utils

import (
	"math/rand"
)

var insults []string = []string{
	"MA CHE CAZZO VOI A LESBICAAA",
	"Sei scemo o mangi i sassi?",
	"Sì il cazzo quello ti piace",
	"Ma che oooooh!",
	"E io che cazzo ne so scusi",
	"Senti che puzzo Bruno",
}

func RandomInsult() string {
	randomIndex := rand.Intn(len(insults))

	return insults[randomIndex]
}
