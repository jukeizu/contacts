package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cheapRoc/grpc-zerolog"
	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/jukeizu/contacts/api/protobuf-spec/contactspb"
	"github.com/oklog/run"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

var Version = ""

const (
	DefaultGrpcPort = "50051"
)

var (
	flagServer  = false
	flagMigrate = false
	flagVersion = false

	grpcPort = DefaultGrpcPort
	dbUrl    = "root@localhost:26257"
)

func parseConfig() {
	flag.StringVar(&grpcPort, "grpc.port", grpcPort, "grpc port")
	flag.StringVar(&dbUrl, "db", dbUrl, "Database connection url")
	flag.BoolVar(&flagServer, "server", false, "Start as server")
	flag.BoolVar(&flagMigrate, "migrate", false, "Run db migrations")
	flag.BoolVar(&flagVersion, "v", false, "version")

	flag.Parse()
}

func main() {
	parseConfig()

	if flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().
		Str("instance", xid.New().String()).
		Str("component", "service.contacts").
		Str("version", Version).
		Logger()

	grpcLoggerV2 := grpczerolog.New(logger.With().Str("transport", "grpc").Logger())
	grpclog.SetLoggerV2(grpcLoggerV2)

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    5 * time.Minute,
				Timeout: 10 * time.Second,
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             5 * time.Second,
				PermitWithoutStream: true,
			},
		))

	contactsRepository, err := NewRepository(dbUrl)
	if err != nil {
		logger.Error().Err(err).Caller().Msg("couldn't create contacts repository")
		os.Exit(1)
	}

	if flagMigrate {
		err := contactsRepository.Migrate()
		if err != nil {
			logger.Error().Err(err).Caller().Msg("couldn't migrate contacts repository")
			os.Exit(1)
		}
	}

	contactsService := NewService(contactsRepository)
	contactsService = NewLoggingService(logger, contactsService)
	contactspb.RegisterContactsServer(grpcServer, contactsService)

	port := ":" + grpcPort

	g := run.Group{}

	g.Add(func() error {
		logger.Info().
			Str("transport", "grpc").
			Str("addr", port).
			Msg("listening")
		listener, err := net.Listen("tcp", port)
		if err != nil {
			logger.Error().Err(err).Caller().Msg("couldn't listener")
			return err
		}

		return grpcServer.Serve(listener)
	}, func(error) {
		grpcServer.GracefulStop()
	})

	cancel := make(chan struct{})
	g.Add(func() error {
		return interrupt(cancel)
	}, func(error) {
		close(cancel)
	})

	logger.Info().Err(g.Run()).Msg("stopped")
}

func interrupt(cancel <-chan struct{}) error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)

	select {
	case <-cancel:
		return errors.New("stopping")
	case sig := <-c:
		return fmt.Errorf("%s", sig)
	}
}
