package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "6379"
	}

	l, err := net.Listen("tcp", ":"+port)

	fmt.Printf("listening to port %s", port)

	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := l.Accept()

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		resp := NewReader(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
