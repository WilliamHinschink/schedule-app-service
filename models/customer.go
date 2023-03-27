package models

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"io"
)

type Customer struct {
	ID        gocql.UUID `json:"id"`
	Firstname string     `json:"firstname"`
}

type Customers []*Customer

func GetCustomers() Customers {
	return customerList
}

func (c *Customers) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(c)
}

var customerList = []*Customer{
	&Customer{
		ID:        gocql.TimeUUID(),
		Firstname: "William",
	},
	&Customer{
		ID:        gocql.TimeUUID(),
		Firstname: "Noah",
	},
}
