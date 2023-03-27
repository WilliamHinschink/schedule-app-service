package repository

import (
	b64 "encoding/base64"
	"github.com/gocql/gocql"
	"log"
	"schedule-app-service/models"
)

type Instance struct {
	log *log.Logger
}

const (
	LOG_ERROR          = "Error on save ::: "
	LOG_ERROR_FIND_ALL = "Error on find ::: "
)

const (
	TABLE        = "customer"
	FIELD_ID     = "id"
	FIELD_TEXT   = "firstname"
	SELECT       = "SELECT " + FIELD_ID + ", " + FIELD_TEXT + " FROM " + TABLE
	SELECT_BY_ID = "SELECT " + FIELD_ID + ", " + FIELD_TEXT + " FROM " + TABLE + " WHERE " + FIELD_ID + " = ?"
	INSERT       = "INSERT INTO " + TABLE + " (" + FIELD_ID + ", " + FIELD_TEXT + ") VALUES (?, ?)"
	DELETE       = "DELETE from " + TABLE + " WHERE " + FIELD_ID + " = ?"
	UPDATE       = "UPDATE " + TABLE + " SET " + FIELD_TEXT + " = ? WHERE " + FIELD_ID + " = ? IF EXISTS"
	PAGE_SIZE    = 10
)

func (i Instance) GetCustomers(session *gocql.Session, state string) ([]models.Customer, string, error) {
	return i.findAll(session, state)
}

func (i Instance) PostCustomer(session *gocql.Session, customer *models.Customer) (*models.Customer, error) {
	return i.save(session, customer)
}

func (i Instance) save(session *gocql.Session, customer *models.Customer) (*models.Customer, error) {
	customer.ID = gocql.TimeUUID()
	if err := session.Query(INSERT, customer.ID, customer.Firstname).Exec(); err != nil {
		i.log.Println(LOG_ERROR, err)
		return customer, err
	}
	return customer, nil
}

func (i Instance) GetCustomerById(session *gocql.Session, id gocql.UUID) (*models.Customer, error) {
	return i.findOne(session, id)
}

func (i Instance) findOne(session *gocql.Session, id gocql.UUID) (*models.Customer, error) {
	var t models.Customer
	if err := session.Query(SELECT_BY_ID, id).Scan(&t.ID, &t.Firstname); err != nil {
		i.log.Println(LOG_ERROR, err)
		return nil, err
	}
	return &t, nil
}

func (i Instance) findAll(session *gocql.Session, state string) ([]models.Customer, string, error) {
	var ts []models.Customer
	var t models.Customer
	session.SetPageSize(PAGE_SIZE)
	query := session.Query(SELECT)
	if state != "" {
		st, _ := b64.StdEncoding.DecodeString(state)
		query.PageState(st)
	}
	it := query.Iter()
	sw := it.WillSwitchPage()
	for !sw && it.Scan(&t.ID, &t.Firstname) {
		ts = append(ts, t)
		sw = it.WillSwitchPage()
	}
	if err := it.Close(); err != nil {
		i.log.Println(LOG_ERROR_FIND_ALL, err)
		return []models.Customer{}, b64.StdEncoding.EncodeToString(it.PageState()), nil
	} else if ts == nil {
		return []models.Customer{}, b64.StdEncoding.EncodeToString(it.PageState()), nil
	}
	return ts, b64.StdEncoding.EncodeToString(it.PageState()), nil
}

func (i Instance) DeleteCustomer(db *gocql.Session, id gocql.UUID) error {
	return i.deleteOne(db, id)
}

func (i Instance) deleteOne(session *gocql.Session, id gocql.UUID) error {
	if err := session.Query(DELETE, id).Exec(); err != nil {
		i.log.Println(LOG_ERROR_FIND_ALL, err)
		return err
	}
	return nil
}

func NewCustomerRepository(l *log.Logger) *Instance {
	return &Instance{l}
}
