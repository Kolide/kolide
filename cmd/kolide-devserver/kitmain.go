package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"

	kitlog "github.com/go-kit/kit/log"

	"github.com/kolide/kolide-ose/datastore"
	"github.com/kolide/kolide-ose/kitserver"
)

// this main is temporary. testing the new MakeHandler from kitserver
func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
		ctx      = context.Background()
		logger   kitlog.Logger
	)
	flag.Parse()
	logger = kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.NewContext(logger).With("ts", kitlog.DefaultTimestampUTC)

	ds, _ := datastore.New("mock", "")
	svc, _ := kitserver.NewService(ds)

	httpLogger := kitlog.NewContext(logger).With("component", "http")

	apiHandler := kitserver.MakeHandler(ctx, svc, httpLogger)
	http.Handle("/", accessControl(apiHandler))

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)

}

// cors headers
func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
