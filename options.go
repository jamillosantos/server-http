package srvhttp

import (
	"context"
	"net"
	"time"
)

type serverOpts struct {
	listener          net.Listener
	bindAddress       string
	appName           string
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	maxHeaderBytes    int
	baseContext       func(net.Listener) context.Context
	connContext       func(ctx context.Context, c net.Conn) context.Context
}

func defaultOpts() serverOpts {
	return serverOpts{
		bindAddress:       ":8080",
		appName:           "http-server",
		readTimeout:       15 * time.Second,
		readHeaderTimeout: 10 * time.Second,
		writeTimeout:      30 * time.Second,
		idleTimeout:       0,
		maxHeaderBytes:    1024 * 8, // 8k
	}
}

type Option func(o *serverOpts)

func WithAppName(appName string) Option {
	return func(o *serverOpts) {
		o.appName = appName
	}
}

func WithBindAddress(bindAddress string) Option {
	return func(o *serverOpts) {
		o.bindAddress = bindAddress
	}
}

func Withlistener(l net.Listener) Option {
	return func(o *serverOpts) {
		o.listener = l
	}
}

func WithReadTimeout(readTimeout time.Duration) Option {
	return func(o *serverOpts) {
		o.readTimeout = readTimeout
	}
}

func WithReadHeaderTimeout(readHeaderTimeout time.Duration) Option {
	return func(o *serverOpts) {
		o.readHeaderTimeout = readHeaderTimeout
	}
}

func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(o *serverOpts) {
		o.writeTimeout = writeTimeout
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(o *serverOpts) {
		o.idleTimeout = idleTimeout
	}
}

func WithMaxHeaderBytes(maxHeaderBytes int) Option {
	return func(o *serverOpts) {
		o.maxHeaderBytes = maxHeaderBytes
	}
}

func WithBaseContext(baseContext func(net.Listener) context.Context) Option {
	return func(o *serverOpts) {
		o.baseContext = baseContext
	}
}

func WithConnContext(connContext func(ctx context.Context, c net.Conn) context.Context) Option {
	return func(o *serverOpts) {
		o.connContext = connContext
	}
}
