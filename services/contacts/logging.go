package contacts

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	pb "github.com/jukeizu/contacts/api/contacts"
)

type loggingService struct {
	logger  log.Logger
	Service pb.ContactsServer
}

func NewLoggingService(logger log.Logger, s pb.ContactsServer) pb.ContactsServer {
	return &loggingService{logger, s}
}

func (s loggingService) SetAddress(ctx context.Context, req *pb.SetAddressRequest) (reply *pb.SetAddressReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SetAddress",
			"request", *req,
			"reply", *reply,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	reply, err = s.Service.SetAddress(ctx, req)

	return
}

func (s loggingService) SetPhone(ctx context.Context, req *pb.SetPhoneRequest) (reply *pb.SetPhoneReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SetPhone",
			"request", *req,
			"reply", *reply,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	reply, err = s.Service.SetPhone(ctx, req)

	return
}

func (s loggingService) Query(ctx context.Context, req *pb.QueryRequest) (reply *pb.QueryReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Query",
			"request", *req,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	reply, err = s.Service.Query(ctx, req)

	return
}

func (s loggingService) RemoveContact(ctx context.Context, req *pb.RemoveContactRequest) (reply *pb.RemoveContactReply, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "RemoveContact",
			"request", *req,
			"reply", *reply,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	reply, err = s.Service.RemoveContact(ctx, req)

	return
}
