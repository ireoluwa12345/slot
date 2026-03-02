package main

import "sync"

var slot = map[string]string{}
var hashedSlot = map[string]map[string]string{}
var slotMutex = sync.RWMutex{}
var hashedSlotMutex = sync.RWMutex{}

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	slotMutex.Lock()
	slot[key] = value
	slotMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	slotMutex.RLock()
	value, ok := slot[key]
	slotMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	hashedSlotMutex.Lock()
	if _, ok := hashedSlot[hash]; !ok {
		hashedSlot[hash] = map[string]string{}
	}
	hashedSlot[hash][key] = value
	hashedSlotMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	hashedSlotMutex.RLock()
	value, ok := hashedSlot[hash][key]
	hashedSlotMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	hashedSlotMutex.RLock()
	values, ok := hashedSlot[hash]
	hashedSlotMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	array := make([]Value, len(values))
	i := 0
	for k, v := range values {
		array[i] = Value{typ: "array", array: []Value{
			Value{typ: "bulk", bulk: k},
			Value{typ: "bulk", bulk: v},
		}}
		i++
	}

	return Value{typ: "array", array: array}
}
