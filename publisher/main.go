package main

import (
	"encoding/json"
	"log"
	"os"
	"wildberries/l0/producer"

	"github.com/nats-io/stan.go"
)

var err error

func main() {
	vv, err := os.ReadFile("testdata/test1.json")
	if err != nil {
		log.Fatal(err)
	}
	var orders []producer.Order
	err = json.Unmarshal([]byte(vv), &orders)
	if err != nil {
		log.Fatal(err)
	}

	sc, err := stan.Connect(producer.ClusterID, producer.ClientID+"1", stan.NatsURL(producer.NATSStreamingURL))
	if err != nil {
		log.Fatal(err)
	}
	for idx, order := range orders {
		o, err := json.Marshal(order)
		if err != nil {
			if err != nil {
				log.Panic(err)
			}
		}
		err = sc.Publish(producer.Channel, o)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("message [%d] send succesfull,  uuid:[%s]", idx, order.OrderUID)
	}
}
