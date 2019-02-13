package main

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func LoggingInterceptor(logger zerolog.Logger) grpc.ServerOption {
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func(begin time.Time) {
			logger := logger.With().
				Str("method", info.FullMethod).
				Str("took", time.Since(begin).String()).
				Logger()

			if err != nil {
				logger.Error().Err(err).Msg("")
				return
			}

			logger.Info().Msg("called")
		}(time.Now())

		resp, err = handler(ctx, req)
		return
	}

	return grpc.UnaryInterceptor(interceptor)
}
