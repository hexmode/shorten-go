package main

import (
	"math/rand"
	"time"
)

func generateKey(size int) string {
	chars := "abcdefghijkmnopqrstuvwxyz23456789ABCDEFGHJKMNPQRSTUVWXYZ"

	rand.Seed(time.Now().UTC().UnixNano())

	out := ""
	for i := 0; len(out) < size; i++ {
		out = out + string(chars[rand.Intn(len(chars))])
	}

	return out
}
