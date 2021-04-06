package config

import (
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"strconv"
	"time"
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
	if gc.BookingCode == "" {
		addError(errs, "go_live.booking_code", gc.BookingCode, "cannot be empty")
	}
	if gc.StartIsoDatetime == "" {
		addError(errs, "go_live.start_iso_datetime", gc.StartIsoDatetime, "cannot be empty")
	} else {
		_, err := time.Parse(StartTimeFormat, gc.StartIsoDatetime)
		if err != nil {
			addError(errs, "go_live.start_iso_datetime", gc.StartIsoDatetime, "is not a valid go live time, format is "+StartTimeFormat)
		}
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
			if key.Size() != 256 {
				addError(errs, "security.token_public_key_PEM", "(omitted)", "has wrong key size, must be RSA256 (2048bit)")
			}
		}
	}
}