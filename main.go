package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/spf13/viper"

	"github.com/sisimogangg/supermarket.products.api/controller"
	_repo "github.com/sisimogangg/supermarket.products.api/repository"
	_service "github.com/sisimogangg/supermarket.products.api/service"
)

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool("debug") {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	router := mux.NewRouter()

	repo := _repo.NewFirebaseRepo()

	timeContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	service := _service.NewProductService(repo, timeContext)

	controller.NewProductHandler(router, service)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}

}
