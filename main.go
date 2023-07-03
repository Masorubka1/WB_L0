package main

import (
	"log"
	"wildberries/l0/http_server"
	"wildberries/l0/producer"
)

func CreatePostgressConfig() producer.PostgresConfig {
	return producer.PostgresConfig{
		Host:     "postgresql",
		Port:     "5432",
		User:     "root",
		Password: "root",
		DBName:   "ordersDB",
	}
}

func main() {
	tmp := "Name"
	log.Println("hello", tmp)
	client, err := producer.NewPostgresClient(CreatePostgressConfig())
	if err != nil {
		log.Println("Failed connect to postgress instance")
	}
	cache := producer.NewMemcacheClient("host.docker.internal:11211")

	sc := producer.NewNatsStream()
	sc.RunNatsSteaming(client, cache)

	http_server.StartHTTPServer(client, cache)
}
