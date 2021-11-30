package countdownctl

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/eurofurence/reg-room-service/api/v1/countdown"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/logging"
	"github.com/eurofurence/reg-room-service/web/util/media"
	"github.com/form3tech-oss/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-http-utils/headers"
)

const isoDateTimeFormat = "2006-01-02T15:04:05-07:00"
const demoPublicSecret = "[demo-secret]"
const demoStaffSecret = "[demo-staff-secret]"

func Create(server chi.Router) {
	server.Get("/api/rest/v1/countdown", getCountdown)
}

func getCurrentTime(ctx context.Context, mockTime string) time.Time {
	if mockTime != "" {
		logging.Ctx(ctx).Info("mock time specified")
		current, err := time.Parse(config.StartTimeFormat, mockTime)
		if err == nil {
			return current
		}
	}
	return time.Now()
}

func getPublicSecret(mockTime string) string {
	if mockTime != "" {
		return demoPublicSecret
	}
	return config.PublicBookingCode()
}

func getStaffSecret(mockTime string) string {
	if mockTime != "" {
		return demoStaffSecret
	}
	return config.StaffBookingCode()
}

func hasStaffClaim(r *http.Request) bool {
	if config.StaffClaimKey() == "" || config.StaffClaimValue() == "" {
		return false
	}

	if user, ok := r.Context().Value("user").(*jwt.Token); ok {
		claims := user.Claims.(jwt.MapClaims)

		if value, ok := claims[config.StaffClaimKey()]; ok {
			return fmt.Sprintf("%v", value) == config.StaffClaimValue()
		}
	}

	return false
}

func getCountdown(w http.ResponseWriter, r *http.Request) {
	logging.Ctx(r.Context()).Info("countdown")

	mockTime := r.URL.Query().Get("currentTimeIso")
	currentTime := getCurrentTime(r.Context(), mockTime)

	publicTargetTime := config.PublicBookingStartTime()
	staffTargetTime := config.StaffBookingStartTime()

	publicCountdownSeconds := int64(math.Round(publicTargetTime.Sub(currentTime).Seconds()))
	staffCountdownSeconds := int64(math.Round(staffTargetTime.Sub(currentTime).Seconds()))

	dto := countdown.CountdownResultDto{}
	dto.CurrentTimeIsoDateTime = currentTime.Format(isoDateTimeFormat)
	dto.CountdownSeconds = publicCountdownSeconds
	dto.TargetTimeIsoDateTime = publicTargetTime.Format(isoDateTimeFormat)
	if publicCountdownSeconds <= 0 {
		dto.Secret = getPublicSecret(mockTime)
		dto.CountdownSeconds = 0
	} else if hasStaffClaim(r) {
		dto.CountdownSeconds = staffCountdownSeconds
		dto.TargetTimeIsoDateTime = staffTargetTime.Format(isoDateTimeFormat)
		if staffCountdownSeconds <= 0 {
			dto.Secret = getStaffSecret(mockTime)
			dto.CountdownSeconds = 0
		}
	}

	w.Header().Add(headers.ContentType, media.ContentTypeApplicationJson)
	w.WriteHeader(http.StatusOK)
	writeJson(r.Context(), w, dto)
}

func writeJson(ctx context.Context, w http.ResponseWriter, v interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		logging.Ctx(ctx).Warn(fmt.Sprintf("error while encoding json response: %v", err))
	}
}
