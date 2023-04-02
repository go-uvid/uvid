package main

import (
	"luvsic3/uvid/api"
)

const DSN = "uvid.db"

func main() {
	api.New(DSN).Start()
}
