package main

import (
	"gochitest/database"
	"gochitest/routes"
	"net/http"
)

func main() {
	database.Connect()

	r := routes.Setup()

	http.ListenAndServe(":3000", r)
}
