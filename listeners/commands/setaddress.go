package commands

import (
	"context"
	"regexp"
	"strings"

	pb "github.com/jukeizu/contacts/api/contacts"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/handler"
)

type SetAddress interface {
	handler.Command
}

type setAddress struct {
	command *command
	logger  log.Logger
}

func (c *command) SetAddress() SetAddress {
	logger := log.With(c.logger, "component", "command.contacts.setaddress")

	return &setAddress{c, logger}
}

func (c *setAddress) IsCommand(request handler.Request) (bool, error) {
	return regexp.MatchString("!setaddress '(.*?)' ([^\\s]*)", request.Content)
}

func (c *setAddress) Handle(request handler.Request) (handler.Results, error) {
	input := strings.SplitAfterN(request.Content, "!setaddress ", 2)[1]

	split := strings.SplitN(input, "'", 3)

	name, address := split[1], strings.TrimSpace(split[2])

	setAddressRequest := pb.SetAddressRequest{
		ServerId: request.ServerId,
		Name:     name,
		Address:  address,
	}

	reply, err := c.command.Client.SetAddress(context.Background(), &setAddressRequest)
	if err != nil {
		return handler.Results{}, err
	}

	result := handler.Result{
		Content: formatContact(reply.Contact),
	}

	return handler.Results{result}, nil
}
