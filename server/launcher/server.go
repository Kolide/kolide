// Package launcher provides a gRPC server to handle launcher requests.
package launcher

import (
	"net/http"
	"strings"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	launcher "github.com/kolide/launcher/service"

	grpc "google.golang.org/grpc"

	"github.com/kolide/fleet/server/health"
	"github.com/kolide/fleet/server/kolide"
)

// Handler extends the grpc.Server, providing Handler that allows us to serve
// both gRPC and twirp over http traffic.
type Handler struct {
	grpcServer       *grpc.Server
	twirpHTTPHandler http.Handler
}

// New creates a gRPC and Twirp server to handle remote requests from launcher.
func New(
	tls kolide.OsqueryService,
	logger log.Logger,
	grpcServer *grpc.Server,
	healthCheckers map[string]health.Checker,
) *Handler {
	var svc launcher.KolideService
	{
		svc = &launcherWrapper{
			tls:            tls,
			logger:         logger,
			healthCheckers: healthCheckers,
		}
		svc = launcher.LoggingMiddleware(logger)(svc)
	}
	endpoints := launcher.MakeServerEndpoints(svc)

	// Build the go-kit GRPC server for Launcher and then connect it back
	// to our primary grpcServer.
	launcherGRPCserver := launcher.NewGRPCServer(endpoints, logger)
	launcher.RegisterGRPCServer(grpcServer, launcherGRPCserver)

	// Build a Twirp HTTP handler with our endpoints.
	twirpHTTPHandler := launcher.NewTwirpHTTPHandler(endpoints, logger)
	return &Handler{grpcServer, twirpHTTPHandler}
}

// Handler will route gRPC traffic to the gRPC server, Twirp traffic to
// the Twirp server (via the mux) and other http traffic will be routed to
// normal http handler functions.
func (h *Handler) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			ctx := r.Context()
			ctx = kithttp.PopulateRequestContext(ctx, r)
			h.grpcServer.ServeHTTP(w, r.WithContext(ctx))
		} else if strings.Contains(r.URL.Path, launcher.TwirpHTTPApiPathPrefix) {
			ctx := r.Context()
			ctx = kithttp.PopulateRequestContext(ctx, r)
			h.twirpHTTPHandler.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (h *Handler) GracefulStop() {
	h.grpcServer.GracefulStop()
	return
}
