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
	Start() error
}

type setPhone struct {
	command *command
	logger  log.Logger
	config  handler.HandlerConfig
}

func (c *command) SetPhone(config handler.HandlerConfig) SetPhone {
	logger := log.With(c.logger, "command", "setphone")

	return &setPhone{c, logger, config}
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

func (c *setPhone) Start() error {
	h, err := handler.NewCommandHandler(c.logger, c.config)
	if err != nil {
		return err
	}

	h.Start(c)

	return nil
}
