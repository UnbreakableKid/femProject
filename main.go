package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/unbreakablekid/femProject/internal/app"
	"github.com/unbreakablekid/femProject/internal/routes"
)

func main() {

	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      r,
	}

	app.Logger.Printf("we are running our app %d\n", port)
	err = server.ListenAndServe()

	if err != nil {
		app.Logger.Fatal(err)
	}

}
