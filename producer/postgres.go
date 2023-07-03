package producer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	_ "github.com/lib/pq"
)

type PostgresClient struct {
	db *sql.DB
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

var (
	once        sync.Once
	dbClient    *PostgresClient
	dbClientErr error
)

func NewPostgresClient(config PostgresConfig) (*PostgresClient, error) {
	once.Do(func() {
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.User, config.Password, config.DBName)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			dbClientErr = err
			return
		}

		err = db.Ping()
		if err != nil {
			dbClientErr = err
			return
		}

		dbClient = &PostgresClient{
			db: db,
		}
	})

	if dbClientErr != nil {
		return nil, dbClientErr
	}

	return dbClient, nil
}

func (c *PostgresClient) Close() {
	c.db.Close()
}

func (c *PostgresClient) InsertOrder(order Order) error {
	query, values := generateInsertQuery("orders", order)

	_, err := c.db.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func generateInsertQuery(tableName string, entity interface{}) (string, []interface{}) {
	var columns []string
	var placeholders []string
	var values []interface{}

	// Get the type and value of the entity
	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)

	// Iterate over the fields of the struct
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		fieldValue := entityValue.Field(i).Interface()

		// Skip fields that are not exported or have the "db" tag set to "-"
		if field.PkgPath != "" || field.Tag.Get("db") == "-" {
			continue
		}

		var insertValue interface{}
		if field.Tag.Get("marshal") != "" {
			insertValue, _ = json.Marshal(fieldValue)
		} else {
			insertValue = fieldValue
		}

		//log.Println(field, fieldValue)
		columns = append(columns, field.Tag.Get("db"))
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(columns)))
		values = append(values, insertValue)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	return query, values
}

func (c *PostgresClient) GetOrderFromPostgres(id string) (*Order, error) {
	// Execute query to retrieve order by ID
	query := "SELECT * FROM orders WHERE order_uid = $1"
	row := c.db.QueryRow(query, id)

	order := &Order{}
	var delivery, payment, items []byte

	err := row.Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry,
		&delivery, &payment, &items,
		&order.Locale, &order.InternalSignature, &order.CustomerID,
		&order.DeliveryService, &order.ShardKey, &order.SMID,
		&order.DateCreated, &order.OOFShard,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(delivery, &order.Delivery)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(payment, &order.Payment)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(items, &order.Items)
	if err != nil {
		return nil, err
	}

	return order, nil
}
