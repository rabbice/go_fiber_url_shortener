package models

import (
	"time"
)

type Request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"custom_short"`
	Expiry      time.Duration `json:"expiry"`
}

type Response struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"custom_short"`
	Expiry      time.Duration `json:"expiry"`
	RateLimit   time.Duration `json:"rate_limit"`
	RateRemain  int           `json:"rate_remain"`
}
