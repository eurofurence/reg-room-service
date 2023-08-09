package common

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type testRequest struct {
	Counter int
}

type testResponse struct {
	Counter int
}

func setupHandler(ep Endpoint[testRequest, testResponse], rh RequestHandler[testRequest], resph ResponseHandler[testResponse]) http.Handler {
	return CreateHandler(ep, rh, resph)
}

func TestCreateHandler(t *testing.T) {
	tReq := &testRequest{
		Counter: 0,
	}
	tRes := &testResponse{
		Counter: 0,
	}

	tests := []struct {
		name                    string
		endpoint                Endpoint[testRequest, testResponse]
		reqHandler              RequestHandler[testRequest]
		respHandler             ResponseHandler[testResponse]
		shouldPanic             bool
		expectedError           error
		expectedRequestCounter  int
		expectedResponseCounter int
	}{
		{
			name:       "Should panic when no request handler was provided",
			reqHandler: nil,
			endpoint:   nil,
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				return nil
			},
			shouldPanic: true,
		},
		{
			name: "Should panic when no response handler was provided",
			endpoint: func(ctx context.Context, request *testRequest, w http.ResponseWriter) (*testResponse, error) {
				return tRes, nil
			},
			reqHandler: func(r *http.Request, w http.ResponseWriter) (*testRequest, error) {
				return tReq, nil
			},
			shouldPanic: true,
			respHandler: nil,
		},
		{
			name: "Should panic when no endpoint was provided",
			reqHandler: func(r *http.Request, w http.ResponseWriter) (*testRequest, error) {
				return tReq, nil
			},
			shouldPanic: true,
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				return nil
			},
		},
		{
			name: "Should increase counter when all values are set",
			endpoint: func(ctx context.Context, request *testRequest, w http.ResponseWriter) (*testResponse, error) {
				return tRes, nil
			},
			reqHandler: func(r *http.Request, w http.ResponseWriter) (*testRequest, error) {
				tReq.Counter++
				return tReq, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				return nil
			},
			expectedRequestCounter:  1,
			expectedResponseCounter: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tReq.Counter = 0
			tRes.Counter = 0
			router := chi.NewRouter()

			if tc.shouldPanic {
				require.Panics(t, func() {
					setupHandler(tc.endpoint, tc.reqHandler, tc.respHandler)
				})

				// stop execution of test logic
				return
			}

			router.Method(http.MethodGet, "/", setupHandler(tc.endpoint, tc.reqHandler, tc.respHandler))

			srv := httptest.NewServer(router)
			defer srv.Close()

			req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, fmt.Sprintf("%s/", srv.URL), nil)
			require.NoError(t, err)

			cl := &http.Client{
				Timeout: time.Second * 10,
			}

			resp, err := cl.Do(req)
			require.NoError(t, err)

			require.NotNil(t, resp)

			b, err := io.ReadAll(resp.Body)
			require.NoError(t, resp.Body.Close())
			require.NoError(t, err)

			fmt.Println(string(b))

			require.Equal(t, tc.expectedRequestCounter, tReq.Counter)
			require.Equal(t, tc.expectedResponseCounter, tRes.Counter)
		})
	}
}
