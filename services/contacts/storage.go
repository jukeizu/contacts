package contacts

import (
	pb "github.com/jukeizu/contacts/api/contacts"
	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
)

type ContactStorage interface {
	mdb.Storage

	SetAddress(*pb.SetAddressRequest) (*pb.Contact, error)
	SetPhone(*pb.SetPhoneRequest) (*pb.Contact, error)
	Query(*pb.QueryRequest) ([]*pb.Contact, error)
	RemoveContact(*pb.RemoveContactRequest) error
}

type storage struct {
	mdb.Store
}

func NewContactStorage(dbConfig mdb.DbConfig) (ContactStorage, error) {
	store, err := mdb.NewStorage(dbConfig)

	s := storage{}
	s.Session = store.Session
	s.Collection = store.Collection

	return &s, err
}

func (s *storage) SetAddress(req *pb.SetAddressRequest) (*pb.Contact, error) {
	contact := &pb.Contact{}

	_, err := s.Collection.Upsert(buildQuery(req.ServerId, req.Name), bson.M{"$set": bson.M{"address": req.Address}})
	if err != nil {
		return contact, err
	}

	return s.findContact(req.ServerId, req.Name)
}

func (s *storage) SetPhone(req *pb.SetPhoneRequest) (*pb.Contact, error) {
	contact := &pb.Contact{}

	_, err := s.Collection.Upsert(buildQuery(req.ServerId, req.Name), bson.M{"$set": bson.M{"phone": req.Phone}})
	if err != nil {
		return contact, err
	}

	return s.findContact(req.ServerId, req.Name)
}

func (s *storage) Query(query *pb.QueryRequest) ([]*pb.Contact, error) {
	contacts := []*pb.Contact{}

	err := s.Collection.Find(bson.M{"serverid": query.ServerId}).Sort("name").All(&contacts)

	return contacts, err
}

func (s *storage) RemoveContact(req *pb.RemoveContactRequest) error {
	return s.Collection.Remove(buildQuery(req.ServerId, req.Name))
}

func (s *storage) findContact(serverId string, name string) (*pb.Contact, error) {
	contact := &pb.Contact{}

	err := s.Collection.Find(buildQuery(serverId, name)).One(&contact)

	return contact, err
}

func buildQuery(serverId string, name string) bson.M {
	bsonQuery := []bson.M{
		bson.M{"serverid": serverId},
		bson.M{"name": name},
	}

	return bson.M{"$and": bsonQuery}
}
