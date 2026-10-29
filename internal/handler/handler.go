package handler

import (
	"encoding/json"
	"fmt"
	"httpServer/internal/database"
	"httpServer/internal/telemetry"
	"io"
	"net/http"
	"strings"

	"github.com/VictoriaMetrics/metrics"
)

var repository = database.NewRepository()

func Telemetry(w http.ResponseWriter, r *http.Request) {
	metrics.WritePrometheus(w, true)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := repository.GetUsers()
	writeJson(w, http.StatusOK, users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {

	}
	defer r.Body.Close()

	bodyText := strings.Split(strings.TrimSpace(string(body)), "\n")
	if len(bodyText) != 2 {
		writeError(w, 500, "ERROR")
		return
	}
	name := bodyText[0]
	id := bodyText[1]
	repository.CreateUser(name, id)
	writeJson(w, http.StatusCreated, name)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	telemetry.RequestCounter.Inc()
	fmt.Fprintf(w, "hello")
}

func writeJson(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJson(w, status, map[string]string{"error": message})
}
