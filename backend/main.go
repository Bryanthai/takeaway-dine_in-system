package main

import _ "github.com/go-sql-driver/mysql"

import(
	"log"
	"net/http"

	"github.com/rs/cors"
)


func main() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("POST /user/login", loginHandler)
	serveMux.HandleFunc("POST /user/register", registerHandler)
	serveMux.HandleFunc("PUT /user/alter", alterAccountHandler)
	serveMux.HandleFunc("POST /food/create", createFoodHandler)
	serveMux.HandleFunc("PUT /food/alter", alterFoodHandler)
	serveMux.HandleFunc("DELETE /food/item", deleteFoodHandler)
	serveMux.HandleFunc("POST /order/new", createOrderHandler)
	serveMux.HandleFunc("GET /user/order", getOrderStatusHandler)
	serveMux.HandleFunc("POST /user/rate", rateOrderHandler)
	serveMux.HandleFunc("DELETE /order/item", deleteOrderHandler)
	serveMux.HandleFunc("GET /food/foods", getAllFoodHandler)

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