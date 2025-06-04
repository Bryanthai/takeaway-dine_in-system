package main

import _ "github.com/go-sql-driver/mysql"

import(
	"log"
	"net/http"
	"fmt"

	"github.com/rs/cors"
)

func initServeMux(serveMux *http.ServeMux) {
	serveMux.HandleFunc("POST /users/login", loginHandler) //done
	serveMux.HandleFunc("POST /users/register", registerHandler) //done
	serveMux.HandleFunc("PUT /users/change-info", alterAccountHandler) //done
	serveMux.HandleFunc("GET /users", getCurrentAccount) //done
	
	serveMux.HandleFunc("PUT /payment", MakePayment) //done
	serveMux.HandleFunc("PUT /recharge", RechargeAccount) //done

	serveMux.HandleFunc("POST /foods", createFoodHandler) //done
	serveMux.HandleFunc("PUT /foods/change-info", alterFoodHandler) //done
	serveMux.HandleFunc("DELETE /foods", deleteFoodHandler) //done
	serveMux.HandleFunc("GET /foods", getFoodByIdHandler) //done


	serveMux.HandleFunc("POST /orders", createOrderHandler) //done
	serveMux.HandleFunc("DELETE /orders", deleteOrderHandler) //done
	serveMux.HandleFunc("GET /orders/user", getOrderStatusHandler) //done
	serveMux.HandleFunc("GET /orders/info", GetOrderByIdHandler) //done
	serveMux.HandleFunc("PUT /orders/rate", rateOrderedItems) //done
	serveMux.HandleFunc("PUT /orders/feedback", UpdateFeedback) //done
	serveMux.HandleFunc("GET /orders/items", getOrderedItemsHandler) //done
	serveMux.HandleFunc("PUT /orders/finish", finishOrder) //done
	serveMux.HandleFunc("GET /orders/price", GetOrderTotalPrice) //done

	serveMux.HandleFunc("GET /menu", getAllFoodHandler) //done
	serveMux.HandleFunc("GET /menu/rating-times-info", getFoodRatingandOrderedTimesByFoodID) //done
	serveMux.HandleFunc("GET /menu/sort-type", getFoodByTypeHandler) //done
	serveMux.HandleFunc("GET /menu/sort-by-usertag", getFoodByUserTag) //done

	serveMux.HandleFunc("GET /admin/users", GetAllUsers) //done
	serveMux.HandleFunc("GET /admin/undone-orders", getAllUndoneOrder) //done
	serveMux.HandleFunc("GET /admin/total-average", GetAverageSpendingAll) //done
	serveMux.HandleFunc("GET /admin/total-average-by-user", GetAverageSpendingByUser) //done
	serveMux.HandleFunc("GET /admin/orders-all", getAllOrdersHandler) //done

	serveMux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./static/images/food"))))
}

func main() {
	serveMux := http.NewServeMux()
	initServeMux(serveMux)
	fmt.Println("Server is running on port 8080...")

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