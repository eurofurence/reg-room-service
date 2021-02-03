package config

import "time"

func BookingStartTime() time.Time {
	t, _ := time.Parse(StartTimeFormat, configuration().GoLive.StartIsoDatetime)
	return t
}

func BookingCode() string {
	return configuration().GoLive.BookingCode
}