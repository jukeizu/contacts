package main

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/contacts/api/protobuf-spec/contactspb"
	migration "github.com/jukeizu/contacts/migrations"
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
	RemoveContact(*contactspb.RemoveContactRequest) (bool, error)
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

	q := `INSERT INTO contact (serverid, name, address) 
		VALUES($1, $2, $3) 
		ON CONFLICT (serverid, name) DO UPDATE SET address = excluded.address, updated = NOW()
		RETURNING serverid, name, address, phone`

	err := r.Db.QueryRow(q,
		req.ServerId,
		req.Name,
		req.Address,
	).Scan(
		&contact.ServerId,
		&contact.Name,
		&contact.Address,
		&contact.Phone,
	)

	return contact, err
}

func (r *repository) SetPhone(req *contactspb.SetPhoneRequest) (*contactspb.Contact, error) {
	contact := &contactspb.Contact{}

	q := `INSERT INTO contact (serverid, name, phone) 
		VALUES($1, $2, $3) 
		ON CONFLICT (serverid, name) DO UPDATE SET phone = excluded.phone, updated = NOW()
		RETURNING serverid, name, address, phone`

	err := r.Db.QueryRow(q,
		req.ServerId,
		req.Name,
		req.Phone,
	).Scan(
		&contact.ServerId,
		&contact.Name,
		&contact.Address,
		&contact.Phone,
	)

	return contact, err
}

func (r *repository) Query(query *contactspb.QueryRequest) ([]*contactspb.Contact, error) {
	contacts := []*contactspb.Contact{}

	q := `SELECT serverid, name, address, phone FROM contact WHERE serverid = $1`

	rows, err := r.Db.Query(q, query.ServerId)
	if err != nil {
		return contacts, err
	}

	defer rows.Close()
	for rows.Next() {
		contact := contactspb.Contact{}
		err := rows.Scan(
			&contact.ServerId,
			&contact.Name,
			&contact.Address,
			&contact.Phone,
		)
		if err != nil {
			return contacts, err
		}

		contacts = append(contacts, &contact)
	}

	return contacts, nil
}

func (r *repository) RemoveContact(req *contactspb.RemoveContactRequest) (bool, error) {
	q := `DELETE FROM contact WHERE serverid = $1 AND name = $2`

	result, err := r.Db.Exec(q, req.ServerId, req.Name)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
