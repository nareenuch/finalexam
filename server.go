package main

import (
	"github.com/nareenuch/myapi/customers"
)

func main() {
	customers.CreateTable()
	r := customers.Setup()
	r.Run(":2019")
}