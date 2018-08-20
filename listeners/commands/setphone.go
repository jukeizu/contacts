package commands

import (
	"context"
	"regexp"
	"strings"

	pb "github.com/jukeizu/contacts/api/contacts"

	"github.com/go-kit/kit/log"
	"github.com/jukeizu/handler"
)

type SetPhone interface {
	handler.Command
}

type setPhone struct {
	command *command
	logger  log.Logger
}

func (c *command) SetPhone() SetPhone {
	logger := log.With(c.logger, "component", "command.contacts.setphone")

	return &setPhone{c, logger}
}

func (c *setPhone) IsCommand(request handler.Request) (bool, error) {
	return regexp.MatchString("!setphone '(.*?)' ([^\\s]*)", request.Content)
}

func (c *setPhone) Handle(request handler.Request) (handler.Results, error) {
	input := strings.SplitAfterN(request.Content, "!setphone", 2)[1]

	split := strings.SplitN(input, "'", 3)

	name, phone := split[1], strings.TrimSpace(split[2])

	setPhoneRequest := pb.SetPhoneRequest{
		ServerId: request.ServerId,
		Name:     name,
		Phone:    phone,
	}

	reply, err := c.command.Client.SetPhone(context.Background(), &setPhoneRequest)
	if err != nil {
		return handler.Results{}, err
	}

	result := handler.Result{
		Content: formatContact(reply.Contact),
	}

	return handler.Results{result}, nil
}
