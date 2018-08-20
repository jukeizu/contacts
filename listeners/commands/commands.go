package commands

import (
	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/contacts/api/contacts"
	"github.com/jukeizu/handler"
)

type Config struct {
	Endpoint string
	Handler  handler.HandlerConfig
}

type Command interface {
	Query() Query
	SetAddress() SetAddress
	SetPhone() SetPhone
	RemoveContact() RemoveContact
}

type command struct {
	logger log.Logger
	Client pb.ContactsClient
}

func NewCommand(logger log.Logger, client pb.ContactsClient) Command {
	return &command{logger, client}
}
