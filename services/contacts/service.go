package contacts

import (
	"context"

	pb "github.com/jukeizu/contacts/api/contacts"
)

type service struct {
	ContactStorage ContactStorage
}

func NewService(contactStorage ContactStorage) pb.ContactsServer {
	return &service{contactStorage}
}

func (s service) SetAddress(ctx context.Context, req *pb.SetAddressRequest) (*pb.SetAddressReply, error) {
	contact, err := s.ContactStorage.SetAddress(req)

	return &pb.SetAddressReply{Contact: contact}, err
}

func (s service) SetPhone(ctx context.Context, req *pb.SetPhoneRequest) (*pb.SetPhoneReply, error) {
	contact, err := s.ContactStorage.SetPhone(req)

	return &pb.SetPhoneReply{Contact: contact}, err
}

func (s service) Query(ctx context.Context, req *pb.QueryRequest) (*pb.QueryReply, error) {
	contacts, err := s.ContactStorage.Query(req)

	return &pb.QueryReply{Contacts: contacts}, err
}

func (s service) RemoveContact(ctx context.Context, req *pb.RemoveContactRequest) (*pb.RemoveContactReply, error) {
	err := s.ContactStorage.RemoveContact(req)
	if err != nil {
		return &pb.RemoveContactReply{Removed: false}, err
	}

	return &pb.RemoveContactReply{Removed: true}, err
}
