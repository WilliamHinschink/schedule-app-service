package config

import (
	"github.com/gocql/gocql"
	"log"
	"os"
	"schedule-app-service/utils"
)

type CassandraConfig struct {
	log *log.Logger
}

const (
	KEYSPACE_NAME      = "relationalkeyspace"
	CASSANDRA_URL      = "CASSANDRA_URL"
	CASSANDRA_USERNAME = "CASSANDRA_USERNAME"
	CASSANDRA_PASSWORD = "CASSANDRA_PASSWORD"
	PAGE_SIZE          = 10
)

func (config *CassandraConfig) InitCluster() (*gocql.Session, error) {
	cassandra := utils.GetEnvDefault(CASSANDRA_URL, "localhost:9042")
	cluster := gocql.NewCluster(cassandra)
	cluster.Keyspace = KEYSPACE_NAME
	cluster.Consistency = gocql.One

	session, err := cluster.CreateSession()

	if err != nil {
		config.log.Println(err)
		return session, err
	}
	return session, nil
}

func (config *CassandraConfig) envVar(key string) string {
	cassandra, exists := os.LookupEnv(key)
	if !exists {
		config.log.Fatal("Valor de " + key + " no encontrado")
	} else {
		config.log.Println("Recuperado valor de " + key)
	}
	return cassandra
}

func NewCassandraDatabase(log *log.Logger) *CassandraConfig {
	return &CassandraConfig{log: log}
}
