package contacts

import (
	"context"

	"github.com/jukeizu/contacts/api/protobuf-spec/contactspb"
)

type service struct {
	ContactStorage ContactStorage
}

func NewService(contactStorage ContactStorage) contactspb.ContactsServer {
	return &service{contactStorage}
}

func (s service) SetAddress(ctx context.Context, req *contactspb.SetAddressRequest) (*contactspb.SetAddressReply, error) {
	contact, err := s.ContactStorage.SetAddress(req)

	return &contactspb.SetAddressReply{Contact: contact}, err
}

func (s service) SetPhone(ctx context.Context, req *contactspb.SetPhoneRequest) (*contactspb.SetPhoneReply, error) {
	contact, err := s.ContactStorage.SetPhone(req)

	return &contactspb.SetPhoneReply{Contact: contact}, err
}

func (s service) Query(ctx context.Context, req *contactspb.QueryRequest) (*contactspb.QueryReply, error) {
	contacts, err := s.ContactStorage.Query(req)

	return &contactspb.QueryReply{Contacts: contacts}, err
}

func (s service) RemoveContact(ctx context.Context, req *contactspb.RemoveContactRequest) (*contactspb.RemoveContactReply, error) {
	err := s.ContactStorage.RemoveContact(req)
	if err != nil {
		return &contactspb.RemoveContactReply{Removed: false}, err
	}

	return &contactspb.RemoveContactReply{Removed: true}, nil
}
