package countdownctl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eurofurence/reg-room-service/api/v1/countdown"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"github.com/eurofurence/reg-room-service/internal/repository/logging"
	"github.com/eurofurence/reg-room-service/web/util/media"
	"github.com/go-chi/chi"
	"github.com/go-http-utils/headers"
	"math"
	"net/http"
	"time"
)

const isoDateTimeFormat = "2006-01-02T15:04:05-07:00"


func Create(server chi.Router) {
	server.Get("/api/rest/v1/countdown", getCountdown)
}

func getCountdown(w http.ResponseWriter, r *http.Request) {
	logging.Ctx(r.Context()).Info("countdown")
	targetTime := config.BookingStartTime()
	currentTime := time.Now()
	dto := countdown.CountdownResultDto{}
	dto.CountdownSeconds = int64(math.Round(targetTime.Sub(currentTime).Seconds()))
	dto.TargetTimeIsoDateTime = targetTime.Format(isoDateTimeFormat)
	dto.CurrentTimeIsoDateTime = currentTime.Format(isoDateTimeFormat)

	if dto.CountdownSeconds <= 0 {
		dto.Secret = config.BookingCode()
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
