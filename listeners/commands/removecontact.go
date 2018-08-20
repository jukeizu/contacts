package commands

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/contacts/api/contacts"
	"github.com/jukeizu/handler"
)

type RemoveContact interface {
	handler.Command
	Start() error
}

type removeContact struct {
	command *command
	logger  log.Logger
	config  handler.HandlerConfig
}

func (c *command) RemoveContact(config handler.HandlerConfig) RemoveContact {
	logger := log.With(c.logger, "component", "command.contacts.removecontact")

	return &removeContact{c, logger, config}
}

func (c *removeContact) IsCommand(request handler.Request) (bool, error) {
	return regexp.MatchString("!setphone '(.*?)' ([^\\s]*)", request.Content)
}

func (c *removeContact) Handle(request handler.Request) (handler.Results, error) {
	input := strings.SplitAfterN(request.Content, "!removecontact ", 2)[1]

	name := strings.Split(input, "'")[1]

	removeContactRequest := pb.RemoveContactRequest{
		ServerId: request.ServerId,
		Name:     name,
	}

	reply, _ := c.command.Client.RemoveContact(context.Background(), &removeContactRequest)

	removeResponse := ""

	if reply.Removed {
		removeResponse = "removed"
	} else {
		removeResponse = "could not remove"
	}

	result := handler.Result{
		Content: fmt.Sprintf("%s '%s'", removeResponse, name),
	}

	return handler.Results{result}, nil
}

func (c *removeContact) Start() error {
	h, err := handler.NewCommandHandler(c.logger, c.config)
	if err != nil {
		return err
	}

	h.Start(c)

	return nil
}
