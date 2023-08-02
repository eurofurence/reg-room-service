package countdownctl

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	v1 "github.com/eurofurence/reg-room-service/internal/api/v1"
	"github.com/eurofurence/reg-room-service/internal/web/filter/jwt"
	"github.com/eurofurence/reg-room-service/internal/web/util/media"

	"github.com/go-chi/chi/v5"
	"github.com/go-http-utils/headers"

	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/logging"
)

const (
	isoDateTimeFormat = "2006-01-02T15:04:05-07:00"
	demoPublicSecret  = "[demo-secret]"
	demoStaffSecret   = "[demo-staff-secret]"
)

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
	if config.StaffRole() == "" {
		return false
	}

	ctx := r.Context()

	staff, _ := jwt.HasRole(ctx, config.StaffRole())
	if staff {
		user, err := jwt.GetName(ctx)
		if err != nil {
			logging.Ctx(ctx).Warn("staff claim found but user name not found - not allowing early access")
			return false
		}

		logging.Ctx(ctx).Info(fmt.Sprintf("staff claim found for user '%s' - allowing early access", user))
		return true
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

	dto := v1.CountdownResultDto{}
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
