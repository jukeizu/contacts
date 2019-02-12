package contacts

import (
	"context"

	"github.com/jukeizu/contacts/api/protobuf-spec/contactspb"
)

type service struct {
	Repository Repository
}

func NewService(contactStorage Repository) contactspb.ContactsServer {
	return &service{contactStorage}
}

func (s service) SetAddress(ctx context.Context, req *contactspb.SetAddressRequest) (*contactspb.SetAddressReply, error) {
	contact, err := s.Repository.SetAddress(req)

	return &contactspb.SetAddressReply{Contact: contact}, err
}

func (s service) SetPhone(ctx context.Context, req *contactspb.SetPhoneRequest) (*contactspb.SetPhoneReply, error) {
	contact, err := s.Repository.SetPhone(req)

	return &contactspb.SetPhoneReply{Contact: contact}, err
}

func (s service) Query(ctx context.Context, req *contactspb.QueryRequest) (*contactspb.QueryReply, error) {
	contacts, err := s.Repository.Query(req)

	return &contactspb.QueryReply{Contacts: contacts}, err
}

func (s service) RemoveContact(ctx context.Context, req *contactspb.RemoveContactRequest) (*contactspb.RemoveContactReply, error) {
	removed, err := s.Repository.RemoveContact(req)

	return &contactspb.RemoveContactReply{Removed: removed}, err
}
