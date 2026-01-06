package engine

import (
	"strconv"
	"strings"

	"ferrodb/internal/parser"
	"ferrodb/internal/persistence"
	"ferrodb/internal/storage"
)

type Engine struct {
	store *storage.MemoryStore
	aof   *persistence.AOF
}

func New() *Engine {
	store := storage.NewMemoryStore()

	aof, err := persistence.OpenAOF("data/ferrodb.aof")
	if err != nil {
		panic(err)
	}

	engine := &Engine{
		store: store,
		aof:   aof,
	}

	// 🔁 Replay data from disk
	aof.Replay(func(line string) {
		engine.executeInternal(line, false)
	})

	return engine
}

func (e *Engine) Execute(input string) string {
	return e.executeInternal(input, true)
}

func (e *Engine) executeInternal(input string, persist bool) string {
	cmd := parser.Parse(input)

	switch cmd.Name {
	case "SET":
		if len(cmd.Args) < 2 {
			return "ERR SET requires key and value"
		}
		e.store.Set(cmd.Args[0], cmd.Args[1])
		if persist {
			e.aof.Write(input)
		}
		return "OK"

	case "GET":
		if len(cmd.Args) < 1 {
			return "ERR GET requires key"
		}
		val, ok := e.store.Get(cmd.Args[0])
		if !ok {
			return "(nil)"
		}
		return val

	case "DEL":
		if len(cmd.Args) < 1 {
			return "ERR DEL requires key"
		}
		e.store.Del(cmd.Args[0])
		if persist {
			e.aof.Write(input)
		}
		return "OK"

	case "EXPIRE":
		if len(cmd.Args) < 2 {
			return "ERR EXPIRE requires key and seconds"
		}

		seconds, err := strconv.ParseInt(cmd.Args[1], 10, 64)
		if err != nil || seconds <= 0 {
			return "ERR invalid TTL"
		}

		ok := e.store.Expire(cmd.Args[0], seconds)
		if !ok {
			return "(nil)"
		}

		return "OK"

	case "HELP":
		return strings.Join([]string{
			"SET key value",
			"GET key",
			"DEL key",
			"EXPIRE key seconds",
			"EXIT",
		}, "\n")

	default:
		return "ERR unknown command"
	}
}
