package main

import (
	"fmt"

	"github.com/kumekay/kuhomon-server/device_builder/builder"
	"github.com/kumekay/kuhomon-server/server/db"
)

func main() {
	database := db.SetupDB()
	defer database.Close()

	device, tokens := builder.BuildDevice(database)

	fmt.Println("New device successfully created and inserted to DB")
	fmt.Println("Device ID:", device.ID)
	fmt.Println("Read Token:", tokens.ReadToken)
	fmt.Println("Write Token:", tokens.WriteToken)
}
