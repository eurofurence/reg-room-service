package web

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"net/http"
)

type (
	RequestHandler[Req any]  func(r *http.Request, w http.ResponseWriter) (*Req, error)
	ResponseHandler[Res any] func(ctx context.Context, res *Res, w http.ResponseWriter) error
	Endpoint[Req, Res any]   func(ctx context.Context, request *Req, w http.ResponseWriter) (*Res, error)
)

func CreateHandler[Req, Res any](
	endpoint Endpoint[Req, Res],
	requestHandler RequestHandler[Req],
	responseHandler ResponseHandler[Res],
) http.Handler {
	if endpoint == nil {
		panic("unable to set up service: no endpoint provided")
	}

	if requestHandler == nil {
		panic("unable to set up service: request handler must not be nil")
	}

	if responseHandler == nil {
		panic("unable to set up service: response handler must not be nil")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, common.CtxKeyRequestURL{}, r.URL)

		defer func() {
			err := r.Body.Close()
			if err != nil {
				aulogging.WarnErrf(ctx, err, "error while closing the request body: %v", err)
			}
		}()

		request, err := requestHandler(r, w)
		if err != nil {
			SendErrorResponse(ctx, w, err)
			return
		}

		response, err := endpoint(ctx, request, w)
		if err != nil {
			SendErrorResponse(ctx, w, err)
			return
		}

		if err := responseHandler(ctx, response, w); err != nil {
			// cannot SendErrorResponse(ctx, w, err) - likely already have started sending in responseHandler
			aulogging.ErrorErrf(ctx, err, "An error occurred during the handling of the response - response may have been incomplete. [error]: %v", err)
		}
	})
}
