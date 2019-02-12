package contacts

import (
	"github.com/jukeizu/contacts/api/protobuf-spec/contactspb"
)

type ContactStorage interface {
	SetAddress(*contactspb.SetAddressRequest) (*contactspb.Contact, error)
	SetPhone(*contactspb.SetPhoneRequest) (*contactspb.Contact, error)
	Query(*contactspb.QueryRequest) ([]*contactspb.Contact, error)
	RemoveContact(*contactspb.RemoveContactRequest) error
}

type storage struct {
}

func NewContactStorage() (ContactStorage, error) {
	s := storage{}
	return &s, nil
}

func (s *storage) SetAddress(req *contactspb.SetAddressRequest) (*contactspb.Contact, error) {
	contact := &contactspb.Contact{}

	return contact, nil
}

func (s *storage) SetPhone(req *contactspb.SetPhoneRequest) (*contactspb.Contact, error) {
	contact := &contactspb.Contact{}

	return contact, nil
}

func (s *storage) Query(query *contactspb.QueryRequest) ([]*contactspb.Contact, error) {
	contacts := []*contactspb.Contact{}

	return contacts, nil
}

func (s *storage) RemoveContact(req *contactspb.RemoveContactRequest) error {
	return nil
}
