package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type Server struct {
	http     *http.Server
	listener net.Listener
	grpc     *grpc.Server
}

type Configuration func(s *Server) error

func New(configs ...Configuration) (r *Server, err error) {
	r = &Server{}
	for _, cfg := range configs {
		if err = cfg(r); err != nil {
			return
		}
	}
	return
}

func (s *Server) Run() (err error) {
	if s.http != nil {
		go func() {
			if err = s.http.ListenAndServe(); err != nil {
				fmt.Printf("ERR_SERVE_HTTP: %v\n", err)
				return
			}
		}()
	}

	if s.grpc != nil {
		go func() {
			if err = s.grpc.Serve(s.listener); err != nil {
				fmt.Printf("ERR_SERVE_GRPC: %v\n", err)
				return
			}
		}()
	}
	return
}
func (s *Server) Stop(ctx context.Context) (err error) {
	if s.http != nil {
		if err = s.http.Shutdown(ctx); err != nil {
			return
		}
	}

	if s.grpc != nil {
		s.grpc.GracefulStop()
	}

	return
}
func WithHTTPServer(handler http.Handler, port string) Configuration {
	return func(s *Server) error {
		s.http = &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		}
		return nil
	}
}

func WithGRPCServer(port string) Configuration {
	return func(s *Server) (err error) {
		s.listener, err = net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
		if err != nil {
			return
		}
		s.grpc = &grpc.Server{}

		return
	}
}
