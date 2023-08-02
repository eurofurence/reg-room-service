package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/eurofurence/reg-room-service/internal/apierrors"
	"github.com/eurofurence/reg-room-service/internal/logging"
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
		expectedError           error
		expectedRequestCounter  int
		expectedResponseCounter int
		expectedStatus          int
	}{
		{
			name:       "Should do nothing when no request handler was provided",
			reqHandler: nil,
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return tRes, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				return nil
			},
			expectedRequestCounter:  0,
			expectedResponseCounter: 0,
			expectedStatus:          http.StatusInternalServerError,
		},
		{
			name: "Should do nothing when no response handler was provided",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return tRes, nil
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				return tReq, nil
			},
			respHandler:             nil,
			expectedRequestCounter:  0,
			expectedResponseCounter: 0,
			expectedStatus:          http.StatusInternalServerError,
		},
		{
			name: "Should increase counter when all values are set",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return tRes, nil
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				tReq.Counter++
				return tReq, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				return nil
			},
			expectedRequestCounter:  1,
			expectedResponseCounter: 1,
			expectedStatus:          http.StatusOK,
		},
		{
			name: "Should return bad request when request validation failed",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return tRes, nil
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				tReq.Counter++
				return nil, errors.New("error error error")
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				return nil
			},
			expectedRequestCounter:  1,
			expectedResponseCounter: 0,
			expectedStatus:          http.StatusBadRequest,
		},
		{
			name: "Should return internal server error when endpoint returns an error",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return nil, errors.New("Endpoint failed")
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				tReq.Counter++
				return tReq, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				return nil
			},
			expectedRequestCounter:  1,
			expectedResponseCounter: 0,
			expectedStatus:          http.StatusInternalServerError,
		},
		{
			name: "Should return internal server error when Response Handler returns an error",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return tRes, nil
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				tReq.Counter++
				return tReq, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				return errors.New("Error sending response")
			},
			expectedRequestCounter:  1,
			expectedResponseCounter: 1,
			expectedStatus:          http.StatusInternalServerError,
		},
		{
			name: "Should successfully return result when nothing failed",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return tRes, nil
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				tReq.Counter++
				return tReq, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				require.NoError(t, json.NewEncoder(w).Encode(res))
				return errors.New("Error sending response")
			},
			expectedRequestCounter:  1,
			expectedResponseCounter: 1,
			expectedStatus:          http.StatusOK,
		},
		{
			name: "Should return specific error if business logic returns StatusError",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return nil, apierrors.NewBadRequest("request was bad :(")
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				tReq.Counter++
				return tReq, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				require.NoError(t, json.NewEncoder(w).Encode(res))
				return errors.New("Error sending response")
			},
			expectedRequestCounter:  1,
			expectedResponseCounter: 0,
			expectedError:           apierrors.NewBadRequest("request was bad :("),
			expectedStatus:          http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tReq.Counter = 0
			tRes.Counter = 0
			router := chi.NewRouter()
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

			require.Equal(t, tc.expectedStatus, resp.StatusCode)

		})
	}

}

func TestStatusErrors(t *testing.T) {
	tRes := &testResponse{
		Counter: 0,
	}

	tReq := &testRequest{
		Counter: 0,
	}

	tests := []struct {
		name                   string
		endpoint               Endpoint[testRequest, testResponse]
		reqHandler             RequestHandler[testRequest]
		respHandler            ResponseHandler[testResponse]
		expectedError          error
		expectedRequestCounter int
		expectedStatus         int
	}{
		{
			name: "Should return bad request if business logic returns StatusError",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return nil, apierrors.NewBadRequest("request was bad :(")
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				tReq.Counter++
				return tReq, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				require.NoError(t, json.NewEncoder(w).Encode(res))
				return errors.New("Error sending response")
			},
			expectedRequestCounter: 1,
			expectedError:          apierrors.NewBadRequest("request was bad :("),
			expectedStatus:         http.StatusBadRequest,
		},
		{
			name: "Should return unauthorized if business logic returns StatusError",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return nil, apierrors.NewUnauthorized("unauthorized token")
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				tReq.Counter++
				return tReq, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				require.NoError(t, json.NewEncoder(w).Encode(res))
				return errors.New("Error sending response")
			},
			expectedRequestCounter: 1,
			expectedError:          apierrors.NewUnauthorized("unauthorized token"),
			expectedStatus:         http.StatusUnauthorized,
		},
		{
			name: "Should return forbidden if business logic returns StatusError",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return nil, apierrors.NewForbidden("forbidden")
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				tReq.Counter++
				return tReq, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				require.NoError(t, json.NewEncoder(w).Encode(res))
				return errors.New("Error sending response")
			},
			expectedRequestCounter: 1,
			expectedError:          apierrors.NewForbidden("forbidden"),
			expectedStatus:         http.StatusForbidden,
		},
		{
			name: "Should return not found if business logic returns StatusError",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return nil, apierrors.NewNotFound("not found")
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				tReq.Counter++
				return tReq, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				require.NoError(t, json.NewEncoder(w).Encode(res))
				return errors.New("Error sending response")
			},
			expectedRequestCounter: 1,
			expectedError:          apierrors.NewNotFound("not found"),
			expectedStatus:         http.StatusNotFound,
		},
		{
			name: "Should return conflict if business logic returns StatusError",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return nil, apierrors.NewConflict("conflict")
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				tReq.Counter++
				return tReq, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				require.NoError(t, json.NewEncoder(w).Encode(res))
				return errors.New("Error sending response")
			},
			expectedRequestCounter: 1,
			expectedError:          apierrors.NewConflict("conflict"),
			expectedStatus:         http.StatusConflict,
		},
		{
			name: "Should return internal server error if business logic returns StatusError",
			endpoint: func(ctx context.Context, request *testRequest, logger logging.Logger) (*testResponse, error) {
				return nil, apierrors.NewInternalServerError("internal server error")
			},
			reqHandler: func(r *http.Request) (*testRequest, error) {
				tReq.Counter++
				return tReq, nil
			},
			respHandler: func(ctx context.Context, res *testResponse, w http.ResponseWriter) error {
				res.Counter++
				require.NoError(t, json.NewEncoder(w).Encode(res))
				return errors.New("Error sending response")
			},
			expectedRequestCounter: 1,
			expectedError:          apierrors.NewInternalServerError("internal server error"),
			expectedStatus:         http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tReq.Counter = 0
			tRes.Counter = 0
			router := chi.NewRouter()
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
			require.Equal(t, 0, tRes.Counter)

			require.Equal(t, tc.expectedStatus, resp.StatusCode)

		})
	}

}
