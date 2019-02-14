package treediagram

import (
	"time"

	"github.com/jukeizu/contract"
	"github.com/rs/zerolog"
)

type handlerLogger struct {
	handler Handler
	logger  zerolog.Logger
}

func NewHandlerLogger(handler Handler, logger zerolog.Logger) Handler {
	return &handlerLogger{handler, logger}
}

func (l *handlerLogger) SetAddress(request contract.Request) (response *contract.Response, err error) {
	defer func(begin time.Time) {
		l := l.logger.With().
			Str("intent", "setaddress").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("called")
	}(time.Now())

	response, err = l.handler.SetAddress(request)

	return
}

func (l *handlerLogger) SetPhone(request contract.Request) (response *contract.Response, err error) {
	defer func(begin time.Time) {
		l := l.logger.With().
			Str("intent", "setphone").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("called")
	}(time.Now())

	response, err = l.handler.SetPhone(request)

	return
}

func (l *handlerLogger) Query(request contract.Request) (response *contract.Response, err error) {
	defer func(begin time.Time) {
		l := l.logger.With().
			Str("intent", "query").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("called")
	}(time.Now())

	response, err = l.handler.Query(request)

	return
}

func (l *handlerLogger) RemoveContact(request contract.Request) (response *contract.Response, err error) {
	defer func(begin time.Time) {
		l := l.logger.With().
			Str("intent", "removecontact").
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			l.Error().Err(err).Msg("")
			return
		}

		l.Info().Msg("called")
	}(time.Now())

	response, err = l.handler.RemoveContact(request)

	return
}

func (l *handlerLogger) Start() error {
	l.logger.Info().Msg("starting handler")
	return l.handler.Start()
}

func (l *handlerLogger) Stop() error {
	l.logger.Info().Msg("stopping handler")
	return l.handler.Stop()
}
