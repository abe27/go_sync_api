package controllers

import (
	"database/sql"
	"fmt"

	"github.com/abe27/syncapi/configs"
	_ "github.com/godror/godror"
	_ "gopkg.in/goracle.v2"
)

func FetchTest() {
	fmt.Println(configs.ORAC_DNS)
	db, err := sql.Open("goracle", configs.ORAC_DNS)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
