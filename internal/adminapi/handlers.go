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
		http.NotFound(w, r)
		return
	}

	db, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "invalid db", http.StatusBadRequest)
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

	http.NotFound(w, r)
}

// --- keys ---

func listKeys(w http.ResponseWriter, db int) {
	jsonOK(w)

	res := eng.Execute(db, "KEYS *")
	keys := []string{}

	if res != "(nil)" && res != "" {
		keys = strings.Split(res, "\n")
	}

	json.NewEncoder(w).Encode(map[string]any{
		"db":   db,
		"keys": keys,
	})
}

// --- key ---

func handleKey(w http.ResponseWriter, r *http.Request, db int, rest []string) {
	if len(rest) < 1 {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	key := rest[0]

	switch r.Method {
	case "GET":
		getKey(w, db, key)
	case "POST":
		setKey(w, r, db)
	case "DELETE":
		delKey(w, db, key)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getKey(w http.ResponseWriter, db int, key string) {
	jsonOK(w)

	val := eng.Execute(db, "GET "+key)
	if val == "(nil)" {
		http.NotFound(w, nil)
		return
	}

	ttl := eng.Execute(db, "TTL "+key)

	json.NewEncoder(w).Encode(map[string]any{
		"key":   key,
		"value": val,
		"ttl":   ttl,
	})
}

func setKey(w http.ResponseWriter, r *http.Request, db int) {
	jsonOK(w)

	var body struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	eng.Execute(db, "SET "+body.Key+" "+body.Value)
	w.WriteHeader(http.StatusCreated)
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
