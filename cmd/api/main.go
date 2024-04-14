package main

import (
	"fmt"
	"open-api-games/internal/transport/rest"
)

func main() {
	err := rest.Start()
	if err != nil {
		fmt.Println(err)
	}
}
