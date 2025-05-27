package http

import (
	"context"
	"errors"
	"fmt"
	corsx "github.com/rs/cors"
	"net"
	"net/http"
)

// newServer returns a http.Server configured with the provided handler and a
// base context which will make the handler request have it.
func newServer(ctx context.Context, handler http.Handler) *http.Server {
	server := &http.Server{
		Handler: handler,
	}

	// serve requests with our own context
	server.BaseContext = func(ln net.Listener) context.Context {
		return ctx
	}

	return server
}

// listenAndServe creates the server and listens and then serves it.
// Once is listens, closes the readyCh so the clients can start requesting data.
func listenAndServe(ctx context.Context, address string, server *http.Server, readyCh chan struct{}) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	close(readyCh)

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(ctx)
	}()

	err = server.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// CORSOpt is a set of CORS options
type CORSOpt struct {
	AllowedOrigins     []string
	AllowedMethods     []string
	AllowedHeaders     []string
	OptionsPassthrough bool
}

// RunServer execute http server
func RunServer(ctx context.Context, serverURL string, httpHandler http.Handler, corsOpt *CORSOpt) error {
	fmt.Println("http system start")
	defer func() {
		fmt.Println("http system stop")
	}()

	readyCh := make(chan struct{})
	go func() {
		<-readyCh
		fmt.Println("http system ready to serve")
	}()

	corsHTTPHandler := buildCORS(corsOpt).Handler(httpHandler)

	server := newServer(ctx, corsHTTPHandler)

	return listenAndServe(ctx, serverURL, server, readyCh)
}

func buildCORS(corsOpt *CORSOpt) *corsx.Cors {
	var cors *corsx.Cors
	if corsOpt == nil {
		cors = corsx.AllowAll()
	} else {
		cors = corsx.New(corsx.Options{
			AllowedOrigins:     corsOpt.AllowedOrigins,
			AllowedMethods:     corsOpt.AllowedMethods,
			AllowedHeaders:     corsOpt.AllowedHeaders,
			OptionsPassthrough: corsOpt.OptionsPassthrough,
		})
	}
	return cors
}
