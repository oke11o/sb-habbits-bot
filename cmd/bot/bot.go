package main

import (
	"fmt"
	"os"

	"github.com/oke11o/sb_habbits_bot/internal/app"
)

func main() {
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
	fmt.Println("DONE")
}
