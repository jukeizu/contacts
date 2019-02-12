package contacts

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/contacts/api/protobuf-spec/contactspb"
	migration "github.com/jukeizu/contacts/contacts/migrations"
	_ "github.com/lib/pq"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "contact"
)

type Repository interface {
	SetAddress(*contactspb.SetAddressRequest) (*contactspb.Contact, error)
	SetPhone(*contactspb.SetPhoneRequest) (*contactspb.Contact, error)
	Query(*contactspb.QueryRequest) ([]*contactspb.Contact, error)
	RemoveContact(*contactspb.RemoveContactRequest) error
	Migrate() error
}

type repository struct {
	Db *sql.DB
}

func NewRepository(url string) (Repository, error) {
	conn := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", url, DatabaseName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err

	}

	r := repository{
		Db: db,
	}

	return &r, nil
}

func (r *repository) Migrate() error {
	_, err := r.Db.Exec(`CREATE DATABASE IF NOT EXISTS ` + DatabaseName)
	if err != nil {
		return err
	}

	g, err := gossage.New(r.Db)
	if err != nil {
		return err
	}

	err = g.RegisterMigrations(migration.CreateTableContact20190212052552{})
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) SetAddress(req *contactspb.SetAddressRequest) (*contactspb.Contact, error) {
	contact := &contactspb.Contact{}

	return contact, nil
}

func (r *repository) SetPhone(req *contactspb.SetPhoneRequest) (*contactspb.Contact, error) {
	contact := &contactspb.Contact{}

	return contact, nil
}

func (r *repository) Query(query *contactspb.QueryRequest) ([]*contactspb.Contact, error) {
	contacts := []*contactspb.Contact{}

	return contacts, nil
}

func (r *repository) RemoveContact(req *contactspb.RemoveContactRequest) error {
	return nil
}
