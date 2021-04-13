package config

import "time"

func BookingStartTime() time.Time {
	t, _ := time.Parse(StartTimeFormat, configuration().GoLive.Public.StartIsoDatetime)
	return t
}

func ServerAddr() string {
	return ":" + configuration().Server.Port
}

func BookingCode() string {
	return configuration().GoLive.Public.BookingCode
}

func StaffBookingCode() string {
	return configuration().GoLive.Staff.BookingCode
}

func StaffBookingStartTime() time.Time {
	start, _ := time.Parse(StartTimeFormat, configuration().GoLive.Staff.StartIsoDatetime)
	return start
}

func StaffClaimKey() string {
	return configuration().GoLive.Staff.ClaimKey
}

func StaffClaimValue() string {
	return configuration().GoLive.Staff.ClaimValue
}

func IsCorsDisabled() bool {
	return configuration().Security.DisableCors
}

func JWTPublicKey() string {
	return configuration().Security.TokenPublicKeyPEM
}
