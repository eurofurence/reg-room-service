package v1

type CountdownResultDto struct {
	CurrentTimeIsoDateTime string `json:"currentTime"`
	TargetTimeIsoDateTime  string `json:"targetTime"`
	CountdownSeconds       int64  `json:"countdown"`
	Secret                 string `json:"secret,omitempty"`
}
