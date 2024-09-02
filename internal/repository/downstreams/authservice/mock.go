package authservice

import (
	"context"
	"fmt"

	"github.com/eurofurence/reg-room-service/internal/application/common"
)

type Mock interface {
	AuthService

	Reset()
	Enable()
	Recording() []string
	SimulateGetError(err error)
	SetupResponse(idToken string, acToken string, response UserInfoResponse)
}

type MockImpl struct {
	responses        map[string]UserInfoResponse
	recording        []string
	simulateGetError error
	simulatorEnabled bool
}

var (
	_ AuthService = (*MockImpl)(nil)
	_ Mock        = (*MockImpl)(nil)
)

func newMock() Mock {
	return &MockImpl{
		responses: make(map[string]UserInfoResponse),
		recording: make([]string, 0),
	}
}

func (m *MockImpl) UserInfo(ctx context.Context) (UserInfoResponse, error) {
	// idToken missing is allowed if request came in via Authorization header with access token
	idToken, _ := ctx.Value(common.CtxKeyIDToken{}).(string)

	accessToken, ok := ctx.Value(common.CtxKeyAccessToken{}).(string)
	if m.simulatorEnabled && ok {
		key := fmt.Sprintf("userinfo %s %s", idToken, accessToken)
		m.recording = append(m.recording, key)

		if m.simulateGetError != nil {
			return UserInfoResponse{}, m.simulateGetError
		}
		response, ok := m.responses[key]
		if !ok {
			return UserInfoResponse{}, UnauthorizedError
		}

		return response, nil
	} else {
		return UserInfoResponse{}, DownstreamError
	}
}

func (m *MockImpl) IsEnabled() bool {
	return m.simulatorEnabled
}

// only used in tests

func (m *MockImpl) Reset() {
	m.recording = make([]string, 0)
	m.simulateGetError = nil
	m.simulatorEnabled = false
}

func (m *MockImpl) Enable() {
	m.simulatorEnabled = true
}

func (m *MockImpl) Recording() []string {
	return m.recording
}

func (m *MockImpl) SimulateGetError(err error) {
	m.simulateGetError = err
}

func (m *MockImpl) SetupResponse(idToken string, acToken string, response UserInfoResponse) {
	key := fmt.Sprintf("userinfo %s %s", idToken, acToken)
	m.responses[key] = response
}
