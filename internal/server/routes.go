package server

import (
	"encoding/json"
	"esp32/cmd/web"
	"esp32/internal/database"
	"io"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/rs/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", templ.Handler(web.Chart()))
	mux.HandleFunc("/api/data", s.ChartDataHandler)
	mux.HandleFunc("/api/add", s.AddSensorDataHandler)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./cmd/web/favicon.ico")
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return c.Handler(mux)
}

func (s *Server) AddSensorDataHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data struct {
		SensorName string              `json:"SensorName"`
		SensorData database.SensorData `json:"SensorData"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Received data: SensorName: %s, SensorData: %+v", data.SensorName, data.SensorData)

	err = s.db.StoreSensorData(data.SensorName, data.SensorData)
	if err != nil {
		log.Printf("Failed to store sensor data. Err: %v", err)
		http.Error(w, "Failed to store sensor data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sensor data stored successfully"))
}

func (s *Server) ChartDataHandler(w http.ResponseWriter, r *http.Request) {
	SensorName := r.URL.Query().Get("sensor")
	if SensorName == "" {
		http.Error(w, "Missing SensorName query parameter", http.StatusBadRequest)
		return
	}

	resp, err := s.db.GetSensorData(SensorName)
	if err != nil {
		log.Printf("Failed to get sensor data. Err: %v", err)
		http.Error(w, "Failed to get sensor data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
