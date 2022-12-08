package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/abe27/syncapi/configs"
	"github.com/abe27/syncapi/controllers"
	_ "github.com/godror/godror"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}

	configs.ORAC_HOST = os.Getenv("ORAC_HOST")
	port, _ := strconv.ParseInt(os.Getenv("ORAC_PORT"), 10, 64)
	configs.ORAC_PORT = port
	configs.ORAC_SERVICE = os.Getenv("ORAC_SERVICE")
	configs.ORAC_USER = os.Getenv("ORAC_USER")
	configs.ORAC_PASSWORD = os.Getenv("ORAC_PASSWORD")
	configs.ORAC_DNS = fmt.Sprintf(`user="%s" password="%s" connectString="%s:%d/%s"`, configs.ORAC_USER, configs.ORAC_PASSWORD, configs.ORAC_HOST, configs.ORAC_PORT, configs.ORAC_SERVICE)
}

func main() {
	fmt.Println("STEP.1 :===> FetchTest")
	controllers.FetchTest()
}
