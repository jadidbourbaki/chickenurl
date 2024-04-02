package main

import (
	"time"

	"github.com/andrew-d/csmrand"
)

// these should be in a config file
const keyLength = 6
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// uses a very standard url shortening algorithm
func urlToShort() string {
	csmrand.Seed(time.Now().UnixNano()) // cryptographically secure randomness
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[csmrand.Intn(len(charset))]
	}

	return string(shortKey)
}
