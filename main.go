package main

import (
	"log"

	"github.com/orlandokj/just/ui"
)

func main() {
    err := ui.RunUI()
    log.Fatal(err)
}
