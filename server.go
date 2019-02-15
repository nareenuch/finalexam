package main

import (
	"github.com/nareenuch/finalexam/customers"
)

func main() {
	customers.CreateTable()
	r := customers.Setup()
	r.Run(":2019")
}
