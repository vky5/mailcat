package main

import (
	"fmt"

	"github.com/vky5/mailcat/internal/db"
)

func main() {
	db.InitDB()
	fmt.Print("This is main package")
}
