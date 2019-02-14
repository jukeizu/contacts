package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cheapRoc/grpc-zerolog"
	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/jukeizu/contacts/api/protobuf-spec/contactspb"
	"github.com/jukeizu/contacts/treediagram"
	"github.com/oklog/run"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/shawntoffel/gossage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

var Version = ""

var (
	flagMigrate = false
	flagVersion = false
	flagServer  = false
	flagHandler = false

	grpcPort       = "50052"
	httpPort       = "10002"
	dbUrl          = "root@localhost:26257"
	serviceAddress = "localhost:" + grpcPort
)

func parseConfig() {
	flag.StringVar(&grpcPort, "grpc.port", grpcPort, "grpc port for server")
	flag.StringVar(&httpPort, "http.port", httpPort, "http port for handler")
	flag.StringVar(&dbUrl, "db", dbUrl, "Database connection url")
	flag.StringVar(&serviceAddress, "service.addr", serviceAddress, "address of service if not local")
	flag.BoolVar(&flagServer, "server", false, "Run as server")
	flag.BoolVar(&flagHandler, "handler", false, "Run as handler")
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
		Str("component", "contacts").
		Str("version", Version).
		Logger()

	grpcLoggerV2 := grpczerolog.New(logger.With().Str("transport", "grpc").Logger())
	grpclog.SetLoggerV2(grpcLoggerV2)

	if !flagServer && !flagHandler {
		flagServer = true
		flagHandler = true
	}

	if flagMigrate {
		contactsRepository, err := NewRepository(dbUrl)
		if err != nil {
			logger.Error().Err(err).Caller().Msg("couldn't create contacts repository")
			os.Exit(1)
		}

		gossage.Logger = func(format string, a ...interface{}) {
			msg := fmt.Sprintf(format, a...)
			logger.Info().Str("component", "migrator").Msg(msg)
		}

		err = contactsRepository.Migrate()
		if err != nil {
			logger.Error().Err(err).Caller().Msg("couldn't migrate contacts repository")
			os.Exit(1)
		}
	}

	g := run.Group{}

	if flagServer {
		contactsRepository, err := NewRepository(dbUrl)
		if err != nil {
			logger.Error().Err(err).Caller().Msg("couldn't create contacts repository")
			os.Exit(1)
		}

		grpcServer := newGrpcServer(logger)
		contactsServer := NewServer(logger, grpcServer, contactsRepository)
		contactspb.RegisterContactsServer(grpcServer, contactsServer)

		g.Add(func() error {
			return contactsServer.Start(":" + grpcPort)
		}, func(error) {
			contactsServer.Stop()
		})
	}

	if flagHandler {
		clientConn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(),
			grpc.WithKeepaliveParams(
				keepalive.ClientParameters{
					Time:                30 * time.Second,
					Timeout:             10 * time.Second,
					PermitWithoutStream: true,
				},
			),
		)
		if err != nil {
			logger.Error().Err(err).Msg("couldn't dial service address")
			os.Exit(1)
		}

		client := contactspb.NewContactsClient(clientConn)
		handler := treediagram.NewHandler(logger, client, ":"+httpPort)

		g.Add(func() error {
			return handler.Start()
		}, func(error) {
			err := handler.Stop()
			if err != nil {
				logger.Error().Err(err).Caller().Msg("couldn't stop handler")
			}
		})
	}

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

func newGrpcServer(logger zerolog.Logger) *grpc.Server {
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
		),
		LoggingInterceptor(logger),
	)

	return grpcServer
}
