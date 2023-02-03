package srvhttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
)

var (
	ErrNotReady = errors.New("service is not ready")
)

type HTTPServer struct {
	cfg      serverOpts
	ready    atomic.Bool
	server   *http.Server
	serverWg sync.WaitGroup
	listener net.Listener
	init     func(mux *http.ServeMux)
}

func New(init func(mux *http.ServeMux), opts ...Option) *HTTPServer {
	cfg := defaultOpts()
	for _, opt := range opts {
		opt(&cfg)
	}

	return &HTTPServer{
		cfg:  cfg,
		init: init,
	}
}

func (f *HTTPServer) Name() string {
	return f.cfg.appName
}

func (f *HTTPServer) Listen(_ context.Context) error {
	opts := f.cfg

	mux := http.NewServeMux()

	f.server = &http.Server{
		Addr:              "",
		Handler:           mux,
		ReadTimeout:       opts.readTimeout,
		ReadHeaderTimeout: opts.readHeaderTimeout,
		WriteTimeout:      opts.writeTimeout,
		IdleTimeout:       opts.idleTimeout,
		MaxHeaderBytes:    opts.maxHeaderBytes,
		BaseContext:       opts.baseContext,
		ConnContext:       opts.connContext,
	}

	if f.init != nil {
		f.init(mux)
	}

	l := f.cfg.listener

	if l == nil {
		newL, err := net.Listen("tcp", f.cfg.bindAddress)
		if err != nil {
			return err
		}
		l = newL
	}

	f.serverWg.Add(1)
	go func() {
		f.ready.Store(true)
		defer func() {
			f.serverWg.Done()
			f.ready.Store(false)
		}()

		f.listener = l
		_ = f.server.Serve(l)
	}()

	return nil
}

func (f *HTTPServer) Close(ctx context.Context) error {
	err := f.server.Shutdown(ctx)
	f.serverWg.Wait()
	return err
}

// IsReady will return true if the service is ready to accept requests. This is compliant with the
// github.com/jamillosantos/application library.
func (f *HTTPServer) IsReady(_ context.Context) error {
	if v := f.ready.Load(); !v {
		return ErrNotReady
	}
	return nil
}
