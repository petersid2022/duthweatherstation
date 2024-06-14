package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type SensorData struct {
	Timestamp   time.Time `json:"timestamp"`
	Temperature float64   `json:"temperature,omitempty"`
	Humidity    float64   `json:"humidity,omitempty"`
	GasLevel    float64   `json:"gas_level,omitempty"`
	Pressure    float64   `json:"pressure,omitempty"`
}

type Service interface {
	Close() error
	StoreSensorData(sensor string, data SensorData) error
	GetSensorData(sensor string) ([]SensorData, error)
}

type service struct {
	db *sql.DB
}

var (
	dbname     = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dbname)
	return s.db.Close()
}

func (s *service) StoreSensorData(sensor string, data SensorData) error {
	var insertQuery string
	var err error

	log.Printf("Storing data for sensor: %s, Data: %+v", sensor, data)

	switch sensor {
	case "dht11":
		insertQuery = `
			INSERT INTO dht11_data (record_datetime, humidity) 
			VALUES (NOW(), ?)
		`
		_, err = s.db.ExecContext(context.Background(), insertQuery, data.Humidity)
	case "mq135":
		insertQuery = `
			INSERT INTO mq135_data (record_datetime, gas_level) 
			VALUES (NOW(), ?)
		`
		_, err = s.db.ExecContext(context.Background(), insertQuery, data.GasLevel)
	case "bmp180":
		insertQuery = `
			INSERT INTO bmp180_data (record_datetime, temperature, pressure) 
			VALUES (NOW(), ?, ?)
		`
		_, err = s.db.ExecContext(context.Background(), insertQuery, data.Temperature, data.Pressure)
	default:
		return fmt.Errorf("unknown sensor type: %s", sensor)
	}

	if err != nil {
		log.Printf("Failed to insert %s sensor data: %v", sensor, err)
		return fmt.Errorf("failed to insert %s sensor data: %v", sensor, err)
	}

	log.Printf("Successfully stored data for sensor: %s", sensor)
	return nil
}

func (s *service) GetSensorData(sensor string) ([]SensorData, error) {
	var sensorDataList []SensorData

	queries := map[string]string{
		"dht11": `
			SELECT record_datetime, humidity
			FROM dht11_data
			ORDER BY record_datetime DESC
		`,
		"mq135": `
			SELECT record_datetime, gas_level
			FROM mq135_data
			ORDER BY record_datetime DESC
		`,
		"bmp180": `
			SELECT record_datetime, temperature, pressure
			FROM bmp180_data
			ORDER BY record_datetime DESC
		`,
	}

	query, ok := queries[sensor]
	if !ok {
		return nil, fmt.Errorf("unknown sensor type: %s", sensor)
	}

	rows, err := s.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s sensor data: %v", sensor, err)
	}
	defer rows.Close()

	for rows.Next() {
		var sensorData SensorData
		var recordDatetime []byte

		switch sensor {
		case "dht11":
			err = rows.Scan(&recordDatetime, &sensorData.Humidity)
			if err != nil {
				return nil, fmt.Errorf("failed to scan row for %s: %v", sensor, err)
			}
		case "mq135":
			err = rows.Scan(&recordDatetime, &sensorData.GasLevel)
			if err != nil {
				return nil, fmt.Errorf("failed to scan row for %s: %v", sensor, err)
			}
		case "bmp180":
			err = rows.Scan(&recordDatetime, &sensorData.Temperature, &sensorData.Pressure)
			if err != nil {
				return nil, fmt.Errorf("failed to scan row for %s: %v", sensor, err)
			}
		}

		// Convert recordDatetime to time.Time
		sensorData.Timestamp, err = time.Parse("2006-01-02 15:04:05", string(recordDatetime))
		if err != nil {
			return nil, fmt.Errorf("failed to parse datetime for %s: %v", sensor, err)
		}

		sensorDataList = append(sensorDataList, sensorData)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration for %s: %v", sensor, err)
	}

	return sensorDataList, nil
}
