package main

import(
	"database/sql"
	"net/http"
	"encoding/json"
	"context"

	"github.com/Bryanthai/ordersystem/internal/database"
)

const dbURL = "hahant:123456@tcp(localhost:3306)/fooddb?parseTime=true&tls=false"

// LOGIN
func loginHandler(writer http.ResponseWriter, req *http.Request) {
    type LoginRequest struct {
        Username string `json:"username"`
    }
    type LoginResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var loginReq LoginRequest
    if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)
    _, err = queries.GetAccount(context.Background(), loginReq.Username)
    if err != nil {
        http.Error(writer, "Invalid username", http.StatusUnauthorized)
        return
    }

    resp := LoginResponse{Success: true, Message: "Login successful"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

// REGISTER
func registerHandler(writer http.ResponseWriter, req *http.Request) {
    type RegisterRequest struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Address  string `json:"address"`
        Phone    int64	`json:"phone"`
    }
    type RegisterResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var regReq RegisterRequest
    if err := json.NewDecoder(req.Body).Decode(&regReq); err != nil {
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

    _, err = queries.GetAccount(context.Background(), regReq.Username)
    if err == nil {
        resp := RegisterResponse{Success: false, Message: "Username already exists"}
        writer.Header().Set("Content-Type", "application/json")
        json.NewEncoder(writer).Encode(resp)
        return
    }

    err = queries.CreateAccount(context.Background(), database.CreateAccountParams{
        Username: 			regReq.Username,
        Email:    			regReq.Email,
        Address:  			regReq.Address,
        UserPhoneNumber:	regReq.Phone,
    })
	
    if err != nil {
        http.Error(writer, "Failed to create account", http.StatusInternalServerError)
        return
    }

    resp := RegisterResponse{Success: true, Message: "Registration successful"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)

	return
}

// ALTER ACCOUNT
func alterAccountHandler(writer http.ResponseWriter, req *http.Request) {
    type AlterAccountRequest struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Address  string `json:"address"`
        Phone    int64	`json:"phone"`
    }
    type AlterAccountResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var alterReq AlterAccountRequest
    if err := json.NewDecoder(req.Body).Decode(&alterReq); err != nil {
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)
    err = queries.AlterAccount(context.Background(), database.AlterAccountParams{
        Email:          alterReq.Email,
        Address:        alterReq.Address,
        UserPhoneNumber: alterReq.Phone,
        Username:       alterReq.Username,
    })
    if err != nil {
        http.Error(writer, "Failed to update account", http.StatusInternalServerError)
        return
    }

    resp := AlterAccountResponse{Success: true, Message: "Account updated successfully"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
	return
}

// CREATE NEW FOOD
func createFoodHandler(writer http.ResponseWriter, req *http.Request) {
    type CreateFoodRequest struct {
        FoodName    string  `json:"food_name"`
        FoodTag     string  `json:"food_tag"`
        Price       float64 `json:"price"`
        Info        string  `json:"info"`
        Ingredients string  `json:"ingredients"`
        TimeNeeded  string  `json:"time_needed"`
    }
    type CreateFoodResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var foodReq CreateFoodRequest
    if err := json.NewDecoder(req.Body).Decode(&foodReq); err != nil {
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

	var infoStruct sql.NullString
	infoStruct.String = foodReq.Info
	infoStruct.Valid = foodReq.Info != ""

	var timeNeededStruct sql.NullString
	timeNeededStruct.String = foodReq.TimeNeeded
	timeNeededStruct.Valid = foodReq.TimeNeeded != ""
    err = queries.CreateFood(context.Background(), database.CreateFoodParams{
        FoodName:    foodReq.FoodName,
        FoodTag:     foodReq.FoodTag,
        Price:       foodReq.Price,
        Info:        infoStruct,
        Ingredients: foodReq.Ingredients,
        TimeNeeded:  timeNeededStruct,
    })
    if err != nil {
        http.Error(writer, "Failed to create food item", http.StatusInternalServerError)
        return
    }

    resp := CreateFoodResponse{Success: true, Message: "Food item created successfully"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
	return
}


// ALTER FOOD DATA
func alterFoodHandler(writer http.ResponseWriter, req *http.Request) {
    type AlterFoodRequest struct {
        FoodName    string  `json:"food_name"`
        FoodTag     string  `json:"food_tag"`
        Price       float64 `json:"price"`
        Info        string  `json:"info"`
        Ingredients string  `json:"ingredients"`
        TimeNeeded  string  `json:"time_needed"`
    }
    type AlterFoodResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var foodReq AlterFoodRequest
    if err := json.NewDecoder(req.Body).Decode(&foodReq); err != nil {
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

	var infoStruct sql.NullString
	infoStruct.String = foodReq.Info
	infoStruct.Valid = foodReq.Info != ""

	var timeNeededStruct sql.NullString
	timeNeededStruct.String = foodReq.TimeNeeded
	timeNeededStruct.Valid = foodReq.TimeNeeded != ""
    err = queries.AlterFood(context.Background(), database.AlterFoodParams{
        FoodTag:     foodReq.FoodTag,
        Price:       foodReq.Price,
        Info:        infoStruct,
        Ingredients: foodReq.Ingredients,
        TimeNeeded:  timeNeededStruct,
        FoodName:    foodReq.FoodName,
    })
    if err != nil {
        http.Error(writer, "Failed to update food item", http.StatusInternalServerError)
        return
    }

    resp := AlterFoodResponse{Success: true, Message: "Food item updated successfully"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
	return
}

// DELETE FOOD
func deleteFoodHandler(writer http.ResponseWriter, req *http.Request) {
    type DeleteFoodRequest struct {
        FoodName string `json:"food_name"`
    }
    type DeleteFoodResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var delReq DeleteFoodRequest
    if err := json.NewDecoder(req.Body).Decode(&delReq); err != nil {
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)
    err = queries.DeleteFood(context.Background(), delReq.FoodName)
    if err != nil {
        http.Error(writer, "Failed to delete food item", http.StatusInternalServerError)
        return
    }

    resp := DeleteFoodResponse{Success: true, Message: "Food item deleted successfully"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
	return
}

// CREATE NEW ORDER
func createOrderHandler(writer http.ResponseWriter, req *http.Request) {
    type CreateOrderRequest struct {
        UserID    int32  `json:"user_id"`
        OrderInfo string `json:"order_info"`
        IsRanged  bool   `json:"is_ranged"`
    }
    type CreateOrderResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var orderReq CreateOrderRequest
    if err := json.NewDecoder(req.Body).Decode(&orderReq); err != nil {
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)
    err = queries.CreateOrder(context.Background(), database.CreateOrderParams{
        UserID:    orderReq.UserID,
        OrderInfo: orderReq.OrderInfo,
        IsRanged:  orderReq.IsRanged,
    })
    if err != nil {
        http.Error(writer, "Failed to create order", http.StatusInternalServerError)
        return
    }

    resp := CreateOrderResponse{Success: true, Message: "Order created successfully"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
	return
}

// GET ORDER STATUS
func getOrderStatusHandler(writer http.ResponseWriter, req *http.Request) {
    type GetOrderStatusRequest struct {
        UserID int32 `json:"user_id"`
    }
    type GetOrderStatusResponse struct {
        Success bool        		`json:"success"`
        Orders  []database.Order 	`json:"orders"`
        Message string      		`json:"message"`
    }

    var statusReq GetOrderStatusRequest
    if err := json.NewDecoder(req.Body).Decode(&statusReq); err != nil {
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)
    orders, err := queries.GetOrder(context.Background(), statusReq.UserID)
    if err != nil {
        http.Error(writer, "Failed to get orders", http.StatusInternalServerError)
        return
    }

    resp := GetOrderStatusResponse{
        Success: true,
        Orders:  orders,
        Message: "Orders retrieved successfully",
    }
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
}

// RATE ORDER
func rateOrderHandler(writer http.ResponseWriter, req *http.Request) {
    type RateOrderRequest struct {
        OrderID  int32  `json:"order_id"`
        Rating   int32  `json:"rating"`
        Feedback string `json:"feedback"`
    }
    type RateOrderResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var rateReq RateOrderRequest
    if err := json.NewDecoder(req.Body).Decode(&rateReq); err != nil {
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

	var feedbackStruct sql.NullString
	feedbackStruct.String = rateReq.Feedback
	feedbackStruct.Valid = rateReq.Feedback != ""

	var ratingStruct sql.NullInt32
	ratingStruct.Int32 = rateReq.Rating
	ratingStruct.Valid = rateReq.Rating != 0
    err = queries.RateOrder(context.Background(), database.RateOrderParams{
        Rating:   ratingStruct,
        Feedback: feedbackStruct,
        OrderID:  rateReq.OrderID,
    })
    if err != nil {
        http.Error(writer, "Failed to rate order", http.StatusInternalServerError)
        return
    }

    resp := RateOrderResponse{Success: true, Message: "Order rated successfully"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
	return
}

// DELETE ORDER
func deleteOrderHandler(writer http.ResponseWriter, req *http.Request) {
    type DeleteOrderRequest struct {
        OrderID int32 `json:"order_id"`
    }
    type DeleteOrderResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var delReq DeleteOrderRequest
    if err := json.NewDecoder(req.Body).Decode(&delReq); err != nil {
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)
    err = queries.DeleteOrder(context.Background(), delReq.OrderID)
    if err != nil {
        http.Error(writer, "Failed to delete order", http.StatusInternalServerError)
        return
    }

    resp := DeleteOrderResponse{Success: true, Message: "Order deleted successfully"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
	return
}

// SHOW MENU
func getAllFoodHandler(writer http.ResponseWriter, req *http.Request) {
    type GetAllFoodResponse struct {
        Success bool              `json:"success"`
        Foods   []database.Food   `json:"foods"`
        Message string            `json:"message"`
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)
    foods, err := queries.GetAllFood(context.Background())
    if err != nil {
        http.Error(writer, "Failed to get food list", http.StatusInternalServerError)
        return
    }

    resp := GetAllFoodResponse{
        Success: true,
        Foods:   foods,
        Message: "Food list retrieved successfully",
    }
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}
