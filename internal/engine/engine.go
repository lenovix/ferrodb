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
	store := storage.NewMemoryStore(cfg.Engine.DBCount, cfg.Engine.CleanupIntervalSec)

	aof, err := persistence.OpenAOF("data/ferrodb.aof")
	if err != nil {
		panic(err)
	}

	engine := &Engine{
		store:     store,
		aof:       aof,
		startTime: time.Now(),
	}

	// 🔁 Replay AOF (default DB = 0)
	aof.Replay(func(line string) {
		engine.executeInternal(0, line, false)
	})

	return engine
}

func (e *Engine) Execute(db int, input string) string {
	return e.executeInternal(db, input, true)
}

func (e *Engine) executeInternal(db int, input string, persist bool) string {
	cmd := parser.Parse(input)

	switch cmd.Name {

	case "SET":
		if len(cmd.Args) < 2 {
			return "ERR SET requires key and value"
		}

		e.store.Set(db, cmd.Args[0], cmd.Args[1])

		if persist {
			e.aof.Write(fmt.Sprintf("SET %d %s %s", db, cmd.Args[0], cmd.Args[1]))
		}
		return "OK"

	case "GET":
		if len(cmd.Args) < 1 {
			return "ERR GET requires key"
		}

		val, ok := e.store.Get(db, cmd.Args[0])
		if !ok {
			return "(nil)"
		}
		return val

	case "DEL":
		if len(cmd.Args) < 1 {
			return "ERR DEL requires key"
		}

		e.store.Del(db, cmd.Args[0])

		if persist {
			e.aof.Write(fmt.Sprintf("DEL %d %s", db, cmd.Args[0]))
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
		ok := e.store.ExpireAt(db, cmd.Args[0], expireAt)
		if !ok {
			return "(nil)"
		}

		if persist {
			e.aof.Write(fmt.Sprintf("EXPIREAT %d %s %d", db, cmd.Args[0], expireAt))
		}
		return "OK"

	case "EXPIREAT":
		if len(cmd.Args) < 3 {
			return "ERR EXPIREAT requires db key timestamp"
		}

		dbi, _ := strconv.Atoi(cmd.Args[0])
		timestamp, err := strconv.ParseInt(cmd.Args[2], 10, 64)
		if err != nil {
			return "ERR invalid timestamp"
		}

		e.store.ExpireAt(dbi, cmd.Args[1], timestamp)
		return "OK"

	case "TTL":
		if len(cmd.Args) < 1 {
			return "ERR TTL requires key"
		}

		ttl := e.store.TTL(db, cmd.Args[0])
		return strconv.FormatInt(ttl, 10)

	case "PERSIST":
		if len(cmd.Args) < 1 {
			return "ERR PERSIST requires key"
		}

		ok := e.store.Persist(db, cmd.Args[0])
		if !ok {
			return "0"
		}

		if persist {
			e.aof.Write(fmt.Sprintf("PERSIST %d %s", db, cmd.Args[0]))
		}
		return "1"

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
			"TTL key",
			"PERSIST key",
			"BGREWRITEAOF",
			"INFO",
			"SELECT db",
			"ACL WHOAMI",
			"ACL LIST",
			"ACL CAT",
			"AUTH username password",
			"LOGOUT",
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
