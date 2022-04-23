package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/api/option"
)

var FireStoreClient *firestore.Client
var ClientError error

var (
	flags = flag.NewFlagSet(
		`sample --userid=<name>`,
		flag.ExitOnError)
	name = flags.String("userid", "", "The name of the election")
)

func validateFlags() {

	if len(*name) == 0 {
		log.Fatalf("--election cannot be empty")
	}
}

func main() {
	name := os.Getenv("MY_POD_NAME")
	id := os.Getenv("MY_POD_ID")

	InitApp()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	FireStoreClient.Collection("NodeRegistry").Add(context.Background(), map[string]interface{}{
		"nodeId":   time.Now().Unix(),
		"isMaster": false,
		"name":     name,
		"id":       id,
	})
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func InitApp() {
	var opt = option.WithCredentialsFile("service-registry-29e77-firebase-adminsdk-i58d8-bee6cc6294.json")

	var config = &firebase.Config{ProjectID: "service-registry-29e77"}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	FireStoreClient, ClientError = app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	if err != nil {
		log.Fatalln(err)
	}
	if ClientError != nil {
		log.Fatalf("error creating client: %v\n", err)
	}

}
