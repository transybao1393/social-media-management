package main

import (
	"fmt"
	"log"
	"net/http"
	"tiktok_api/app/connector"

	"github.com/spf13/viper"

	"tiktok_api/app/config"
	"tiktok_api/app/logger"
)

// - VIPER lib for config.json reading when init moment
func init() {
	err := config.LoadConfiguration("config")

	//- Please uncomment this line for VSCode debugging
	// err := utils.LoadConfigurationForDebugging()
	if err != nil {
		log.Fatal(err, "Fatal loading config file: ")
	}
}

func main() {

	//- logger initialize
	log := logger.NewLogrusLogger()

	//- go-chi implementation
	r := connector.SetupRouter()

	// Start server
	log.Infof("Program is running. Access http://localhost:%d", viper.GetString("SERVER.PORT"))

	//- router
	fmt.Println("Server start at port " + ":" + viper.GetString("SERVER.PORT"))
	err := http.ListenAndServe(":"+viper.GetString("SERVER.PORT"), r)
	if err != nil {
		log.Fatal(err, "error on serve server")
	}
}
