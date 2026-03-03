package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/ireoluwa12345/slot/internal/resp"
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
		reader := resp.NewReader(conn)
		value, err := reader.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := resp.NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(resp.Value{Typ: "string", Str: ""})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
