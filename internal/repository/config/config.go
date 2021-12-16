package config

import "time"

func PublicBookingStartTime() time.Time {
	t, _ := time.Parse(StartTimeFormat, configuration().GoLive.Public.StartIsoDatetime)
	return t
}

func ServerAddr() string {
	return ":" + configuration().Server.Port
}

func PublicBookingCode() string {
	return configuration().GoLive.Public.BookingCode
}

func StaffBookingCode() string {
	return configuration().GoLive.Staff.BookingCode
}

func StaffBookingStartTime() time.Time {
	start, _ := time.Parse(StartTimeFormat, configuration().GoLive.Staff.StartIsoDatetime)
	return start
}

func StaffRole() string {
	return configuration().GoLive.Staff.StaffRole
}

func IsCorsDisabled() bool {
	return configuration().Security.DisableCors
}

func JWTPublicKey() string {
	return configuration().Security.TokenPublicKeyPEM
}

func JWTCookieName() string {
	return configuration().Security.TokenCookieName
}
