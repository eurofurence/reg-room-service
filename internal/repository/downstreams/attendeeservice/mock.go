package attendeeservice

import (
	"context"
	"errors"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"github.com/eurofurence/reg-room-service/internal/repository/downstreams"
)

type Mock interface {
	AttendeeService

	Reset()
	Unavailable()
	SetupRegistered(subject string, badgeNo int64, status Status, nickname string, email string)
}

type MockImpl struct {
	IdsBySubject  map[string][]int64
	StatusById    map[int64]Status
	AttendeeById  map[int64]Attendee
	IsUnavailable bool
}

func NewMock() Mock {
	instance := &MockImpl{}
	instance.Reset()
	return instance
}

func (m *MockImpl) ListMyRegistrationIds(ctx context.Context) ([]int64, error) {
	if m.IsUnavailable {
		return make([]int64, 0), downstreams.ErrDownStreamUnavailable
	}

	claimsPtr := ctx.Value(common.CtxKeyClaims{})
	if claimsPtr == nil {
		// no auth -> no badge numbers, but also not an error
		return make([]int64, 0), nil
	}
	claims, ok := claimsPtr.(*common.AllClaims)
	if !ok {
		return make([]int64, 0), errors.New("internal error - found invalid data type in context - this indicates a bug")
	}

	subject := claims.Subject
	if subject == "" {
		return make([]int64, 0), errors.New("invalid authentication in context, lacks subject - probably indicates a bug")
	}

	ids, ok := m.IdsBySubject[subject]
	if !ok {
		// not known -> no badge numbers, but also not an error
		return make([]int64, 0), nil
	}

	return ids, nil
}

func (m *MockImpl) GetStatus(ctx context.Context, id int64) (Status, error) {
	if m.IsUnavailable {
		return StatusDeleted, downstreams.ErrDownStreamUnavailable
	}

	status, ok := m.StatusById[id]
	if !ok {
		return StatusDeleted, nil
	}

	return status, nil
}

func (m *MockImpl) GetAttendee(ctx context.Context, id int64) (Attendee, error) {
	if m.IsUnavailable {
		return Attendee{}, downstreams.ErrDownStreamUnavailable
	}

	attendee, ok := m.AttendeeById[id]
	if !ok {
		return Attendee{}, downstreams.ErrByStatus(nil, 404)
	}

	return attendee, nil
}

func (m *MockImpl) Reset() {
	m.IdsBySubject = make(map[string][]int64)
	m.StatusById = make(map[int64]Status)
	m.AttendeeById = make(map[int64]Attendee)
	m.IsUnavailable = false
}

func (m *MockImpl) Unavailable() {
	m.IsUnavailable = true
}

func (m *MockImpl) SetupRegistered(subject string, badgeNo int64, status Status, nickname string, email string) {
	m.IdsBySubject[subject] = []int64{badgeNo}
	m.StatusById[badgeNo] = status
	m.AttendeeById[badgeNo] = Attendee{
		ID:                   badgeNo,
		Nickname:             nickname,
		Email:                email,
		RegistrationLanguage: "en-US",
	}
}
