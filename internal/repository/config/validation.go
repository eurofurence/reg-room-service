package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

func addError(errs validationErrors, key string, value interface{}, message string) {
	errs[key] = append(errs[key], fmt.Sprintf("value '%v' %s", value, message))
}

func validateServerConfiguration(errs validationErrors, sc serverConfig) {
	if sc.Port == "" {
		addError(errs, "server.port", sc.Port, "cannot be empty")
	} else {
		port, err := strconv.ParseUint(sc.Port, 10, 16)
		if err != nil {
			addError(errs, "server.port", sc.Port, "is not a valid port number")
		} else if port <= 1024 {
			addError(errs, "server.port", sc.Port, "must be a nonprivileged port")
		}
	}
}

func validateGoLiveConfiguration(errs validationErrors, gc goLiveConfig) {
	var publicStartTime time.Time
	// public section
	if gc.Public.BookingCode == "" {
		addError(errs, "go_live.public.booking_code", gc.Public.BookingCode, "cannot be empty")
	}
	if gc.Public.StartIsoDatetime == "" {
		addError(errs, "go_live.public.start_iso_datetime", gc.Public.StartIsoDatetime, "cannot be empty")
	} else {
		var err error
		publicStartTime, err = time.Parse(StartTimeFormat, gc.Public.StartIsoDatetime)
		if err != nil {
			addError(errs, "go_live.public.start_iso_datetime", gc.Public.StartIsoDatetime, "is not a valid go live time, format is "+StartTimeFormat)
		}
	}

	// staff section
	// XXX TODO: test coverage for "staff" config section
	if gc.Staff.StartIsoDatetime == "" && gc.Staff.BookingCode == "" && gc.Staff.StaffRole == "" {
		// section is optional
		return
	}
	if gc.Staff.BookingCode == "" {
		addError(errs, "go_live.staff.booking_code", gc.Staff.BookingCode, "cannot be empty")
	}
	if gc.Staff.StartIsoDatetime == "" {
		addError(errs, "go_live.staff.start_iso_datetime", gc.Staff.StartIsoDatetime, "cannot be empty")
	} else {
		staffStartTime, err := time.Parse(StartTimeFormat, gc.Staff.StartIsoDatetime)
		if err != nil {
			addError(errs, "go_live.staff.start_iso_datetime", gc.Staff.StartIsoDatetime, "is not a valid go live time, format is "+StartTimeFormat)
		}
		if publicStartTime.Before(staffStartTime) {
			addError(errs, "go_live.staff.start_iso_datetime", gc.Staff.StartIsoDatetime, "must be earlier than go_live.public.start_iso_datetime")
		}
	}
	if gc.Staff.StaffRole == "" {
		addError(errs, "go_live.staff.staff_role", gc.Staff.StaffRole, "cannot be empty")
	}
}

func validateSecurityConfiguration(errs validationErrors, sc securityConfig) {
	if sc.TokenPublicKeyPEM == "" {
		addError(errs, "security.token_public_key_PEM", sc.TokenPublicKeyPEM, "cannot be empty")
	} else {
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(sc.TokenPublicKeyPEM))
		if err != nil {
			addError(errs, "security.token_public_key_PEM", "(omitted)", "is not a valid RSA256 PEM key")
		} else {
			if key.Size() != 256 && key.Size() != 512 {
				addError(errs, "security.token_public_key_PEM", "(omitted)", "has wrong key size, must be 256 (2048bit) or 512 (4096bit)")
			}
		}
	}
}
