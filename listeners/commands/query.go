package commands

import (
	"bytes"
	"context"

	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/contacts/api/contacts"
	"github.com/jukeizu/handler"
)

type Query interface {
	handler.Command
}

type query struct {
	command *command
	logger  log.Logger
}

func (c *command) Query() Query {
	logger := log.With(c.logger, "component", "command.contacts.query")

	return &query{c, logger}
}

func (c *query) IsCommand(request handler.Request) (bool, error) {
	return request.Content == "!listcontacts", nil
}

func (c *query) Handle(request handler.Request) (handler.Results, error) {
	queryReply, err := c.command.Client.Query(context.Background(), &pb.QueryRequest{ServerId: request.ServerId})
	if err != nil {
		return handler.Results{}, err
	}

	contacts := queryReply.Contacts

	if len(contacts) == 0 {
		result := handler.Result{
			Content: "no results :cry:",
		}

		return handler.Results{result}, nil
	}

	buffer := bytes.Buffer{}
	for _, contact := range contacts {
		buffer.WriteString(formatContact(contact))
	}

	result := handler.Result{
		Content: buffer.String(),
	}

	return handler.Results{result}, nil
}
