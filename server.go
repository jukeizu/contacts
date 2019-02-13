package main

import (
	"context"
	"net"

	"github.com/jukeizu/contacts/api/protobuf-spec/contactspb"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Server struct {
	logger     zerolog.Logger
	grpcServer *grpc.Server
	repository Repository
}

func NewServer(logger zerolog.Logger, grpcServer *grpc.Server, contactStorage Repository) Server {
	logger = logger.With().Str("component", "server.contacts").Logger()

	return Server{logger, grpcServer, contactStorage}
}

func (s Server) SetAddress(ctx context.Context, req *contactspb.SetAddressRequest) (*contactspb.SetAddressReply, error) {
	contact, err := s.repository.SetAddress(req)

	return &contactspb.SetAddressReply{Contact: contact}, err
}

func (s Server) SetPhone(ctx context.Context, req *contactspb.SetPhoneRequest) (*contactspb.SetPhoneReply, error) {
	contact, err := s.repository.SetPhone(req)

	return &contactspb.SetPhoneReply{Contact: contact}, err
}

func (s Server) Query(ctx context.Context, req *contactspb.QueryRequest) (*contactspb.QueryReply, error) {
	contacts, err := s.repository.Query(req)

	return &contactspb.QueryReply{Contacts: contacts}, err
}

func (s Server) RemoveContact(ctx context.Context, req *contactspb.RemoveContactRequest) (*contactspb.RemoveContactReply, error) {
	removed, err := s.repository.RemoveContact(req)

	return &contactspb.RemoveContactReply{Removed: removed}, err
}

func (s Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.logger.Info().
		Str("transport", "grpc").
		Str("addr", addr).
		Msg("listening")

	return s.grpcServer.Serve(listener)
}

func (s Server) Stop() {
	if s.grpcServer == nil {
		return
	}

	s.logger.Info().
		Str("transport", "grpc").
		Msg("stopping")

	s.grpcServer.GracefulStop()
}
