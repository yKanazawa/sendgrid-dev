package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/yKanazawa/sendgrid-dev/route"
)

func main() {
	fmt.Println("SENDGRID_DEV_PORT", os.Getenv("SENDGRID_DEV_PORT"))
	fmt.Println("SENDGRID_DEV_APIKEY", os.Getenv("SENDGRID_DEV_APIKEY"))
	fmt.Println("SENDGRID_SERVER_PORT", os.Getenv("SENDGRID_SMTP_SERVER"))
	fmt.Println("SENDGRID_SMTP_PORT", os.Getenv("SENDGRID_SMTP_PORT"))

	port, err := strconv.Atoi(os.Getenv("SENDGRID_DEV_PORT"))
	if err != nil || port < 0 || port > 65535 {
		port = 3030
	}

	router := route.Init()
	router.Logger.Fatal(router.Start(":3030"))
}
