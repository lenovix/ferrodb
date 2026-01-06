package engine

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"ferrodb/internal/config"
	"ferrodb/internal/parser"
	"ferrodb/internal/persistence"
	"ferrodb/internal/storage"
)

type Engine struct {
	store     *storage.MemoryStore
	aof       *persistence.AOF
	startTime time.Time
}

func New(cfg *config.Config) *Engine {
	store := storage.NewMemoryStore(cfg.Engine.CleanupIntervalSec)

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

		expireAt := time.Now().Unix() + seconds

		ok := e.store.ExpireAt(cmd.Args[0], expireAt)
		if !ok {
			return "(nil)"
		}

		if persist {
			e.aof.Write(
				fmt.Sprintf("EXPIREAT %s %d", cmd.Args[0], expireAt),
			)
		}

		return "OK"

	case "EXPIREAT":
		if len(cmd.Args) < 2 {
			return "ERR EXPIREAT requires key and timestamp"
		}

		timestamp, err := strconv.ParseInt(cmd.Args[1], 10, 64)
		if err != nil {
			return "ERR invalid timestamp"
		}

		ok := e.store.ExpireAt(cmd.Args[0], timestamp)
		if !ok {
			return "(nil)"
		}

		return "OK"

	case "BGREWRITEAOF":
		go e.RewriteAOF()
		return "OK"

	case "INFO":
		return e.Info()

	case "HELP":
		return strings.Join([]string{
			"SET key value",
			"GET key",
			"DEL key",
			"EXPIRE key seconds",
			"BGREWRITEAOF",
			"INFO",
			"EXIT",
		}, "\n")

	default:
		return "ERR unknown command"
	}
}

func (e *Engine) RewriteAOF() string {
	snapshot := e.store.Snapshot()

	if err := e.aof.Rewrite(snapshot); err != nil {
		return "ERR rewrite failed"
	}

	return "OK"
}

func (e *Engine) Shutdown() {
	if e.aof != nil {
		e.aof.Sync()
		e.aof.Close()
	}
}
