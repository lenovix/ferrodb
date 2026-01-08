package adminapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ferrodb/internal/engine"
)

var eng *engine.Engine

func SetEngine(e *engine.Engine) {
	eng = e
}

// helper
func jsonOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// --- health ---

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	jsonOK(w)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// --- dbs ---

func listDBHandler(w http.ResponseWriter, _ *http.Request) {
	jsonOK(w)
	json.NewEncoder(w).Encode(map[string]any{
		"db_count": eng.DBCount(),
	})
}

// --- router ---
// /api/db/{id}/...

func dbRouter(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) < 3 {
		writeJSONError(w, "invalid path", http.StatusNotFound)
		return
	}

	db, err := strconv.Atoi(parts[2])
	if err != nil {
		writeJSONError(w, "invalid db", http.StatusBadRequest)
		return
	}

	// /api/db/{id}/keys
	if len(parts) == 4 && parts[3] == "keys" {
		listKeys(w, db)
		return
	}

	// /api/db/{id}/key/{name}
	if len(parts) >= 5 && parts[3] == "key" {
		handleKey(w, r, db, parts[4:])
		return
	}

	writeJSONError(w, "not found", http.StatusNotFound)
}

// --- keys ---

func listKeys(w http.ResponseWriter, db int) {
	jsonOK(w)

	res := eng.Execute(db, "KEYS *")
	keys := []string{}

	if res != "(nil)" && res != "" {
		for _, k := range strings.Split(res, "\n") {
			if k != "" {
				keys = append(keys, k)
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]any{
		"db":   db,
		"keys": keys,
	})
}

// --- key ---

func handleKey(w http.ResponseWriter, r *http.Request, db int, rest []string) {
	if len(rest) < 1 {
		writeJSONError(w, "missing key", http.StatusBadRequest)
		return
	}

	key := rest[0]

	switch r.Method {
	case http.MethodGet:
		getKey(w, db, key)
	case http.MethodPost:
		setKey(w, r, db, key)
	case http.MethodDelete:
		delKey(w, db, key)
	default:
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getKey(w http.ResponseWriter, db int, key string) {
	jsonOK(w)

	val := eng.Execute(db, "GET "+key)
	if val == "(nil)" {
		writeJSONError(w, "key not found", http.StatusNotFound)
		return
	}

	ttl := eng.Execute(db, "TTL "+key)

	json.NewEncoder(w).Encode(map[string]any{
		"key":   key,
		"value": val,
		"ttl":   ttl,
	})
}

func setKey(w http.ResponseWriter, r *http.Request, db int, key string) {
	jsonOK(w)

	var body struct {
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, "invalid json body", http.StatusBadRequest)
		return
	}

	// allow empty string value
	cmd := "SET " + key + " " + body.Value
	eng.Execute(db, cmd)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"key":    key,
	})
}

func delKey(w http.ResponseWriter, db int, key string) {
	jsonOK(w)

	res := eng.Execute(db, "DEL "+key)
	deleted := 0
	if res == "OK" {
		deleted = 1
	}

	json.NewEncoder(w).Encode(map[string]int{
		"deleted": deleted,
	})
}

func writeJSONError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}
