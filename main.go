package main

import (
	"fmt"
	"luvsic3/uvid/daos"
	"luvsic3/uvid/tools"
	"time"
)

const DSN = "seed.db"

func main() {
	tools.Seed(DSN)
	// TimeRange this week
	dao := daos.New(DSN)
	db := dao.TimeRange(time.Now().AddDate(0, 0, -7), time.Now())
	result := dao.FindPageViewInterval(db)
	fmt.Println(result)
}
