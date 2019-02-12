package contacts

import (
	"context"
	"time"

	"github.com/jukeizu/contacts/api/protobuf-spec/contactspb"
	"github.com/rs/zerolog"
)

type loggingService struct {
	logger  zerolog.Logger
	Service contactspb.ContactsServer
}

func NewLoggingService(logger zerolog.Logger, s contactspb.ContactsServer) contactspb.ContactsServer {
	return &loggingService{logger, s}
}

func (l loggingService) SetAddress(ctx context.Context, req *contactspb.SetAddressRequest) (reply *contactspb.SetAddressReply, err error) {
	defer func(begin time.Time) {
		logger := l.logger.With().
			Str("method", "SetAddress").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			logger.Error().Err(err).Msg("")
			return
		}

		logger.Info().Msg("called")
	}(time.Now())

	reply, err = l.Service.SetAddress(ctx, req)

	return
}

func (l loggingService) SetPhone(ctx context.Context, req *contactspb.SetPhoneRequest) (reply *contactspb.SetPhoneReply, err error) {
	defer func(begin time.Time) {
		logger := l.logger.With().
			Str("method", "SetPhone").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			logger.Error().Err(err).Msg("")
			return
		}

		logger.Info().Msg("called")
	}(time.Now())

	reply, err = l.Service.SetPhone(ctx, req)

	return
}

func (l loggingService) Query(ctx context.Context, req *contactspb.QueryRequest) (reply *contactspb.QueryReply, err error) {
	defer func(begin time.Time) {
		logger := l.logger.With().
			Str("method", "Query").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			logger.Error().Err(err).Msg("")
			return
		}

		logger.Info().Msg("called")
	}(time.Now())

	reply, err = l.Service.Query(ctx, req)

	return
}

func (l loggingService) RemoveContact(ctx context.Context, req *contactspb.RemoveContactRequest) (reply *contactspb.RemoveContactReply, err error) {
	defer func(begin time.Time) {
		logger := l.logger.With().
			Str("method", "RemoveContact").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			logger.Error().Err(err).Msg("")
			return
		}

		logger.Info().Msg("called")
	}(time.Now())

	reply, err = l.Service.RemoveContact(ctx, req)

	return
}
