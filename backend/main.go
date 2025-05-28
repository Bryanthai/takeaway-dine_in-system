package main

import _ "github.com/go-sql-driver/mysql"

import(
	"log"
	"net/http"

	"github.com/rs/cors"
)

func initServeMux(serveMux *http.ServeMux) {
	serveMux.HandleFunc("POST /users/login", loginHandler)
	serveMux.HandleFunc("POST /users/register", registerHandler)
	serveMux.HandleFunc("PUT /users/change-info", alterAccountHandler)
	serveMux.HandleFunc("GET /users/order", getOrderStatusHandler)
	serveMux.HandleFunc("POST /users/rate", rateOrderHandler)

	serveMux.HandleFunc("POST /foods", createFoodHandler)
	serveMux.HandleFunc("PUT /foods/change-info", alterFoodHandler)
	serveMux.HandleFunc("DELETE /foods", deleteFoodHandler)
	serveMux.HandleFunc("GET /foods", getAllFoodHandler)

	serveMux.HandleFunc("POST /orders", createOrderHandler)
	serveMux.HandleFunc("DELETE /orders", deleteOrderHandler)
	serveMux.HandleFunc("GET /orders", getAllOrdersHandler)
}

func main() {
	serveMux := http.NewServeMux()
	initServeMux(serveMux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(serveMux)
	server := http.Server{}
	server.Handler = handler
	server.Addr = ":8080"
	log.Fatal(server.ListenAndServe())
	return
}