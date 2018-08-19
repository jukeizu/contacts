package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/jukeizu/contacts/api/contacts"
	"github.com/jukeizu/contacts/services/contacts"
	mdb "github.com/shawntoffel/GoMongoDb"
	"github.com/shawntoffel/services-core/command"
	"github.com/shawntoffel/services-core/config"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

var serviceArgs command.CommandArgs

func init() {
	serviceArgs = command.ParseArgs()
}

type Config struct {
	Port           int
	ContactStorage mdb.DbConfig
}

func main() {
	logger := logging.GetLogger("services.contacts", os.Stdout)

	c := Config{}
	err := config.ReadConfig(serviceArgs.ConfigFile, &c)
	if err != nil {
		panic(err)
	}

	storage, err := contacts.NewContactStorage(c.ContactStorage)
	if err != nil {
		panic(err)
	}

	service := contacts.NewService(storage)
	service = contacts.NewLoggingService(logger, service)

	errChannel := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errChannel <- fmt.Errorf("%s", <-c)

	}()

	go func() {
		port := fmt.Sprintf(":%d", c.Port)

		listener, err := net.Listen("tcp", port)
		if err != nil {
			logger.Log("error", err.Error())
		}

		s := grpc.NewServer()
		pb.RegisterContactsServer(s, service)

		logger.Log("transport", "grpc", "address", port, "msg", "listening")

		errChannel <- s.Serve(listener)
	}()

	logger.Log("stopped", <-errChannel)
}
