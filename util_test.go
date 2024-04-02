package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortUrl(t *testing.T) {
	shortenedUrl := urlToShort()
	isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(shortenedUrl)

	assert.True(t, isAlphanumeric)
	assert.Equal(t, len(shortenedUrl), 6)
}
