package config

type conf struct {
	GoLive goLiveConfig `yaml:"go_live"`
}

type goLiveConfig struct {
	StartIsoDatetime string `yaml:"start_iso_datetime"`
	BookingCode      string `yaml:"booking_code"`
}

const (
	StartTimeFormat = "2006-01-02T15:04:05-07:00"
)