package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"

	"github.com/spf13/viper"

	"github.com/sisimogangg/supermarket.products.api/handler"
	"github.com/sisimogangg/supermarket.products.api/model"
	_repo "github.com/sisimogangg/supermarket.products.api/repository"
	_service "github.com/sisimogangg/supermarket.products.api/service"
	"github.com/sisimogangg/supermarket.products.api/utils"
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

func seeding(app *firebase.App) {
	ctx := context.Background()
	client, err := app.Database(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var rawProducts map[string]model.Product
	err = client.NewRef("products").Get(ctx, &rawProducts)
	if err != nil {
		log.Fatal(err)
	}

	if len(rawProducts) == 0 {
		for _, p := range utils.Products {
			if err := client.NewRef(fmt.Sprintf("products/%d", p.ID)).Set(ctx, p); err != nil {
				log.Fatal(err)
			}
		}

		for _, d := range utils.Details {
			if err := client.NewRef(fmt.Sprintf("details/%d", d.ID)).Set(ctx, d); err != nil {
				log.Fatal(err)
			}
		}
	}

}

func initializeFirebase() *firebase.App {
	opt := option.WithCredentialsFile("firebaseServiceAccount.json")

	ctx := context.Background()
	config := &firebase.Config{
		DatabaseURL: "https://supermarket-8aee3.firebaseio.com",
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatal(err)
	}
	return app
}

func main() {
	router := mux.NewRouter()

	app := initializeFirebase()
	repo := _repo.NewFirebaseRepo(app)

	if viper.GetBool("debug") {
		seeding(app)
	}

	timeContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	service := _service.NewProductService(repo, timeContext)

	handler.NewHandler(router, service)

	err := http.ListenAndServe(viper.GetString("server.address"), router)
	if err != nil {
		fmt.Print(err)
	}

}
