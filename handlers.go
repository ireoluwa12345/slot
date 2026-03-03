package main

import (
	"sync"

	"github.com/ireoluwa12345/slot/internal/resp"
)

var slot = map[string]string{}
var hashedSlot = map[string]map[string]string{}
var slotMutex = sync.RWMutex{}
var hashedSlotMutex = sync.RWMutex{}

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}

	return resp.Value{Typ: "string", Str: args[0].Bulk}
}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	slotMutex.Lock()
	slot[key] = value
	slotMutex.Unlock()

	return resp.Value{Typ: "string", Str: "OK"}
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	slotMutex.RLock()
	value, ok := slot[key]
	slotMutex.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: value}
}

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	hashedSlotMutex.Lock()
	if _, ok := hashedSlot[hash]; !ok {
		hashedSlot[hash] = map[string]string{}
	}
	hashedSlot[hash][key] = value
	hashedSlotMutex.Unlock()

	return resp.Value{Typ: "string", Str: "OK"}
}

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	hashedSlotMutex.RLock()
	value, ok := hashedSlot[hash][key]
	hashedSlotMutex.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: value}
}

func hgetall(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].Bulk

	hashedSlotMutex.RLock()
	values, ok := hashedSlot[hash]
	hashedSlotMutex.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	array := make([]resp.Value, len(values))
	i := 0
	for k, v := range values {
		array[i] = resp.Value{Typ: "array", Array: []resp.Value{
			resp.Value{Typ: "bulk", Bulk: k},
			resp.Value{Typ: "bulk", Bulk: v},
		}}
		i++
	}

	return resp.Value{Typ: "array", Array: array}
}
