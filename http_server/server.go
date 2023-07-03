package http_server

import (
	"encoding/json"
	"log"
	"time"

	"wildberries/l0/producer"

	"github.com/gin-gonic/gin"
)

var cache *producer.MemcacheClient
var client *producer.PostgresClient

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/order/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Try to get order from cache
		item, err := cache.Get(id)
		if err == nil && len(item) > 0 {
			var order producer.Order
			//var delivery, payment, items []byte
			err = json.Unmarshal(item, &order)
			if err == nil {
				c.JSON(200, order)
				return
			}
			log.Println("Error unmarshaling order from cache:", err)
		} else {
			log.Println("Error getting order from cache:", err)
		}

		// If order is not found in cache, retrieve it from Postgres
		order, err := client.GetOrderFromPostgres(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "Order not found"})
			return
		}

		// Save order to cache
		orderBytes, err := json.Marshal(order)
		if err == nil {
			err = cache.Set(id, orderBytes, time.Second)
			if err != nil {
				log.Println("Error setting order in cache:", err)
			}
		} else {
			log.Println("Error marshaling order:", err)
		}

		c.JSON(200, order)
	})

	return router
}

func StartHTTPServer(local_client *producer.PostgresClient, local_cache *producer.MemcacheClient) {
	cache = local_cache
	client = local_client
	router := setupRouter()

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
