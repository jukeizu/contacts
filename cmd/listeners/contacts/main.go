package main

import (
	"os"

	"golang.org/x/sync/errgroup"

	pb "github.com/jukeizu/contacts/api/contacts"
	"github.com/jukeizu/contacts/listeners/commands"
	"github.com/shawntoffel/services-core/command"
	"github.com/shawntoffel/services-core/config"
	"github.com/shawntoffel/services-core/logging"
	"google.golang.org/grpc"
)

var commandArgs command.CommandArgs

func init() {
	commandArgs = command.ParseArgs()
}

func main() {
	logger := logging.GetLogger("command.contacts", os.Stdout)

	commandConfig := commands.Config{}

	err := config.ReadConfig(commandArgs.ConfigFile, &commandConfig)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial(commandConfig.Endpoint, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := pb.NewContactsClient(conn)
	command := commands.NewCommand(logger, client)

	group := errgroup.Group{}

	group.Go(command.Query(commandConfig.QueryHandler).Start)
	group.Go(command.SetAddress(commandConfig.SetAddressHandler).Start)
	group.Go(command.SetPhone(commandConfig.SetPhoneHandler).Start)
	group.Go(command.RemoveContact(commandConfig.RemoveContactHandler).Start)

	err = group.Wait()

	logger.Log("stopped", err)
}
