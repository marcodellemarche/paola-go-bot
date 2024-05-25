package chatgpt

import (
	"os"
	"strconv"
	"time"
)

const defaultCapacity = 100
const defaultRatePerHour = 100

type ChatGPT struct {
	openaiToken string
	assistantID string
	model       string
	RateLimiter *RateLimiter[int64]
}

func New(openaiToken string, assistantID string) *ChatGPT {
	capacity := intOrDefault(os.Getenv("OPENAI_LIMITER_CAPACITY"), defaultCapacity)
	ratePerHour := intOrDefault(os.Getenv("OPENAI_LIMITER_RATE_PER_HOUR"), defaultRatePerHour)

	return &ChatGPT{
		openaiToken: openaiToken,
		assistantID: assistantID,
		model:       "gpt-3.5-turbo",
		RateLimiter: NewRateLimiter[int64](capacity, ratePerHour, time.Hour),
	}
}

func intOrDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return result
}
