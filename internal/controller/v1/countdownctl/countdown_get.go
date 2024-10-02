package countdownctl

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/application/web"
	"github.com/eurofurence/reg-room-service/internal/repository/config"
	"math"
	"net/http"
	"time"

	modelsv1 "github.com/eurofurence/reg-room-service/internal/api/v1"
)

const (
	mockTimeFormat    = "2006-01-02T15:04:05-07:00"
	isoDateTimeFormat = "2006-01-02T15:04:05-07:00"
	demoPublicSecret  = "[demo-secret]"
	demoStaffSecret   = "[demo-staff-secret]"
)

// GetCountdownRequest is the request for the GetCountdown operation.
type GetCountdownRequest struct {
	mockTime *time.Time
}

// GetCountdown returns the countdown information.
// If the countdown has reached 0, also reveals the configured secret.
func (*Handler) GetCountdown(ctx context.Context, req *GetCountdownRequest, w http.ResponseWriter) (*modelsv1.Countdown, error) {
	conf, err := config.GetApplicationConfig()
	if conf == nil || err != nil {
		aulogging.Warn(ctx, "application config not found - failing request")
		return nil, common.NewInternalServerError(ctx, common.InternalErrorMessage, common.Details("configuration missing"))
	}

	if req.mockTime != nil {
		if hasStaffClaim(ctx, conf) && conf.GoLive.Staff.StartISODatetime != "" {
			target := parseTime(ctx, conf.GoLive.Staff.StartISODatetime)
			return countdown(*req.mockTime, target, demoStaffSecret), nil
		} else {
			target := parseTime(ctx, conf.GoLive.Public.StartISODatetime)
			return countdown(*req.mockTime, target, demoPublicSecret), nil
		}
	} else {
		current := time.Now()
		if hasStaffClaim(ctx, conf) && conf.GoLive.Staff.StartISODatetime != "" {
			target := parseTime(ctx, conf.GoLive.Staff.StartISODatetime)
			return countdown(current, target, conf.GoLive.Staff.BookingCode), nil
		} else {
			target := parseTime(ctx, conf.GoLive.Public.StartISODatetime)
			return countdown(current, target, conf.GoLive.Public.BookingCode), nil
		}
	}
}

func (*Handler) GetCountdownRequest(r *http.Request, w http.ResponseWriter) (*GetCountdownRequest, error) {
	currentTimeIsoParam := r.URL.Query().Get("currentTimeIso")
	if currentTimeIsoParam != "" {
		ctx := r.Context()
		aulogging.Warn(ctx, "mock time specified")
		mockTime, err := time.Parse(mockTimeFormat, currentTimeIsoParam)
		if err != nil {
			return nil, common.NewBadRequest(ctx, common.RequestParseFailed, common.Details("mock time specified but failed to parse"))
		}
		return &GetCountdownRequest{mockTime: &mockTime}, nil
	}

	return &GetCountdownRequest{}, nil
}

func (*Handler) GetCountdownResponse(ctx context.Context, res *modelsv1.Countdown, w http.ResponseWriter) error {
	return web.EncodeWithStatus(http.StatusOK, res, w)
}

func countdown(current time.Time, target time.Time, secret string) *modelsv1.Countdown {
	seconds := int64(math.Round(target.Sub(current).Seconds()))

	result := &modelsv1.Countdown{
		CurrentTimeIsoDateTime: current.Format(isoDateTimeFormat),
		TargetTimeIsoDateTime:  target.Format(isoDateTimeFormat),
		CountdownSeconds:       seconds,
	}

	if seconds <= 0 {
		result.CountdownSeconds = 0
		result.Secret = secret
	}

	return result
}

func hasStaffClaim(ctx context.Context, conf *config.Config) bool {
	group := conf.GoLive.Staff.Group
	if group != "" {
		if common.HasGroup(ctx, group) {
			user := common.GetSubject(ctx)
			if user == "" {
				aulogging.Warn(ctx, "staff claim found but user name not found - not allowing early access")
				return false
			}

			aulogging.Infof(ctx, "staff claim found for user '%s' - allowing early access", user)
			return true
		}
	}

	return false
}

func parseTime(ctx context.Context, targetStr string) time.Time {
	t, err := time.Parse(isoDateTimeFormat, targetStr)
	if err != nil {
		aulogging.Warn(ctx, "target time configuration invalid - returning a time in the far future")
		return time.Unix(1<<63-62135596801, 999999999) // maximally in the future
	}
	return t
}
