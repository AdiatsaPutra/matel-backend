package main

import (
	config "motor/configs"
	"motor/route"
)

func main() {
	db := config.InitDB()
	route.CreateHandler(db)
}
