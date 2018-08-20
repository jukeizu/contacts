package commands

import (
	"github.com/go-kit/kit/log"
	pb "github.com/jukeizu/contacts/api/contacts"
	"github.com/jukeizu/handler"
)

type Config struct {
	Endpoint             string
	QueryHandler         handler.HandlerConfig
	SetAddressHandler    handler.HandlerConfig
	SetPhoneHandler      handler.HandlerConfig
	RemoveContactHandler handler.HandlerConfig
}

type Command interface {
	Query(handler.HandlerConfig) Query
	SetAddress(handler.HandlerConfig) SetAddress
	SetPhone(handler.HandlerConfig) SetPhone
	RemoveContact(handler.HandlerConfig) RemoveContact
}

type command struct {
	logger log.Logger
	Client pb.ContactsClient
}

func NewCommand(logger log.Logger, client pb.ContactsClient) Command {
	return &command{logger, client}
}
