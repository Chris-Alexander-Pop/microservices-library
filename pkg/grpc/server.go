package grpc

import (
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Port string `env:"GRPC_PORT" env-default:"9090"`
}

type Server struct {
	srv *grpc.Server
	cfg Config
	log *slog.Logger
}

func New(cfg Config, log *slog.Logger) *Server {
	srv := grpc.NewServer()
	reflection.Register(srv)
	return &Server{srv: srv, cfg: cfg, log: log}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":"+s.cfg.Port)
	if err != nil {
		return err
	}
	s.log.Info("starting grpc server", "port", s.cfg.Port)
	return s.srv.Serve(lis)
}

func (s *Server) GRPC() *grpc.Server {
	return s.srv
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
