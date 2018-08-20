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
}

type removeContact struct {
	command *command
	logger  log.Logger
}

func (c *command) RemoveContact() RemoveContact {
	logger := log.With(c.logger, "component", "command.contacts.removecontact")

	return &removeContact{c, logger}
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
