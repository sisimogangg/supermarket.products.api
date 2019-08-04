package main

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/micro/go-micro"
	"github.com/spf13/viper"

	discountProto "github.com/sisimogangg/supermarket.discount.api/proto"
	pb "github.com/sisimogangg/supermarket.products.api/proto"

	_repo "github.com/sisimogangg/supermarket.products.api/repository"
	"github.com/sisimogangg/supermarket.products.api/service"
	"github.com/sisimogangg/supermarket.products.api/utils"

	_ "github.com/micro/go-micro/registry/mdns"
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

	var rawProducts map[string]pb.Product
	err = client.NewRef("products").Get(ctx, &rawProducts)
	if err != nil {
		log.Fatal(err)
	}

	if len(rawProducts) == 0 {
		for _, p := range utils.Products {
			if err := client.NewRef(fmt.Sprintf("products/%s", p.Id)).Set(ctx, p); err != nil {
				log.Fatal(err)
			}
		}

		for _, d := range utils.Details {
			if err := client.NewRef(fmt.Sprintf("details/%s", d.Product.Id)).Set(ctx, d); err != nil {
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
	// create a context
	timeContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	srv := micro.NewService(
		micro.Name("supermarket.product"),
		micro.Version("latest"),
	)

	srv.Init()

	app := initializeFirebase()
	repo := _repo.NewFirebaseRepo(app)

	if viper.GetBool("debug") {
		seeding(app)
	}

	discountClient := discountProto.NewDiscountServiceClient("supermarket.discount", srv.Client()) //shippy.service.client
	productService := service.NewProductService(repo, discountClient, timeContext)

	pb.RegisterProductServiceHandler(srv.Server(), productService)

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
