package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/hugo0706/femProject/internal/app"
	"github.com/hugo0706/femProject/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Server port")
	flag.Parse()
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()
	
	app.Logger.Println("App running")
	
	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: r,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	
	err = server.ListenAndServe()
	
	if err != nil {
		app.Logger.Fatal(err)
	}
}

