package main

import(
	"database/sql"
	"net/http"
	"encoding/json"
	"context"
    "time"
    "log"
    "strings"

	"github.com/Bryanthai/ordersystem/internal/database"
    "github.com/Bryanthai/ordersystem/internal/auth"
)

// LOGIN
func loginHandler(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    type LoginRequest struct {
        Username string `json:"username"`
        Password string `json:"password"`
        ExpiresInSeconds int `json:"expires_in_seconds"`
    }
    type LoginResponse struct {
        Success  bool   `json:"success"`
        Message  string `json:"message"`
        Token    string `json:"token"`
        ID       int32  `json:"id"`
        Username string `json:"username"`
    }

    log.Println("Login request received")

    var loginReq LoginRequest
    if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
        log.Println("Error decoding request body:", err)
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        log.Println("Error connecting to database:", err)
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)
    account, err := queries.GetAccount(context.Background(), loginReq.Username)
    if err != nil {
        log.Println("Error fetching account:", err)
        http.Error(writer, "Invalid username", http.StatusUnauthorized)
        return
    }
    err = auth.CheckPasswordHash(loginReq.Password, account.Password)
    if err != nil {
        log.Println("Invalid password for user:", loginReq.Username)
        http.Error(writer, "Invalid password", http.StatusUnauthorized)
        return
    }

    if loginReq.ExpiresInSeconds <= 0 {
        loginReq.ExpiresInSeconds = 3600
    } else if loginReq.ExpiresInSeconds > 3600 {
        loginReq.ExpiresInSeconds = 3600
    }
    token, err := auth.MakeJWT(account.ID, account.Username, time.Duration(loginReq.ExpiresInSeconds) * time.Second)
    if err != nil {
        log.Println("Error creating JWT:", err)
        http.Error(writer, "Error creating JWT", http.StatusInternalServerError)
        return
    }

    resp := LoginResponse{Success: true, Message: "Login successful", Token: token, ID: account.ID, Username: account.Username}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

// REGISTER
func registerHandler(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    type RegisterRequest struct {
        Username string `json:"username"`
        Password string `json:"password"`
        Email    string `json:"email"`
        Address  string `json:"address"`
        Phone    int64	`json:"phone"`
    }
    type RegisterResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    log.Println("Register request received")

    var regReq RegisterRequest
    if err := json.NewDecoder(req.Body).Decode(&regReq); err != nil {
        log.Println("Error decoding request body:", err)
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    hashedPassword, err := auth.HashPassword(regReq.Password)
    if err != nil {
        log.Println("Error hashing password:", err)
        http.Error(writer, "Error hashing password", http.StatusInternalServerError)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        log.Println("Error connecting to database:", err)
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

    _, err = queries.GetAccount(context.Background(), regReq.Username)
    if err == nil {
        log.Println("Username already exists:", regReq.Username)
        resp := RegisterResponse{Success: false, Message: "Username already exists"}
        writer.Header().Set("Content-Type", "application/json")
        json.NewEncoder(writer).Encode(resp)
        return
    }

    err = queries.CreateAccount(context.Background(), database.CreateAccountParams{
        Username: 			regReq.Username,
        Password: 			hashedPassword,
        Email:    			regReq.Email,
        Address:  			regReq.Address,
        UserPhoneNumber:	regReq.Phone,
    })
	
    if err != nil {
        log.Println("Error creating account:", err)
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
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    token, err := auth.GetBearerToken(req.Header)
    if err != nil {
        http.Error(writer, "Invalid or missing token", http.StatusUnauthorized)
        return
    }

    username, userID, err := auth.ValidateJWT(token, authKey)
    if err != nil {
        http.Error(writer, "Invalid token", http.StatusUnauthorized)
        return
    }

    log.Println("Alter account request received for user:", username)

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

    account, err := queries.GetAccount(context.Background(), username)

    if account.ID != userID || account.Username != username {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }
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

func removeWhiteSpace(s string) []string {
    var result []string
    for _, word := range strings.Fields(s) {
        result = append(result, word)
    }
    return result
}

func getUserTag(username string) []string {
    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        log.Println("Error connecting to database:", err)
        return nil
    }
    defer db.Close()

    queries := database.New(db)
    account, err := queries.GetAccount(context.Background(), username)
    if err != nil {
        log.Println("Error fetching account:", err)
        return nil
    }

    if !account.UserTag.Valid {
        return nil
    }

    tags := removeWhiteSpace(account.UserTag.String)
    return tags
}

// GET FOOD BY USER TAG
func getFoodByUserTag(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    token, err := auth.GetBearerToken(req.Header)
    if err != nil {
        http.Error(writer, "Invalid or missing token", http.StatusUnauthorized)
        return
    }

    username, _, err := auth.ValidateJWT(token, authKey)
    if err != nil {
        http.Error(writer, "Invalid token", http.StatusUnauthorized)
        return
    }

    tags := getUserTag(username)
    if tags == nil {
        http.Error(writer, "No user tags found", http.StatusNotFound)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var foods []database.Food
    queries := database.New(db)
    for _, tag := range tags {
        food, err := queries.GetFoodByTag(context.Background(), tag)
        if err != nil {
            http.Error(writer, "Failed to get food by tags", http.StatusInternalServerError)
            return
        }
        foods = append(foods, food...)
    }

    resp := struct {
        Success bool              `json:"success"`
        Foods   []database.Food   `json:"foods"`
        Message string            `json:"message"`
    }{
        Success: true,
        Foods:   foods,
        Message: "Food list retrieved successfully",
    }
    
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

// Make Payment
func MakePayment(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    token, err := auth.GetBearerToken(req.Header)
    if err != nil {
        http.Error(writer, "Invalid or missing token", http.StatusUnauthorized)
        return
    }

    username, userID, err := auth.ValidateJWT(token, authKey)
    if err != nil {
        http.Error(writer, "Invalid token", http.StatusUnauthorized)
        return
    }

    log.Println("Payment request received from user:", username)

    type PaymentRequest struct {
        OrderID int32 `json:"order_id"`
    }
    type PaymentResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var paymentReq PaymentRequest
    if err := json.NewDecoder(req.Body).Decode(&paymentReq); err != nil {
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

    order, err := queries.GetOrderById(context.Background(), paymentReq.OrderID)
    if err != nil || order.UserID != userID || order.IsPaid {
        http.Error(writer, "Invalid order ID or already paid", http.StatusBadRequest)
        return
    }

    // Check if user has enough balance
    account, err := queries.GetAccount(context.Background(), username)
    if err != nil {
        http.Error(writer, "Failed to retrieve account", http.StatusInternalServerError)
        return
    }
    // Get the total price of the order
    items, err := queries.GetOrderedItems(context.Background(), paymentReq.OrderID)
    if err != nil {
        http.Error(writer, "Failed to retrieve order items", http.StatusInternalServerError)
        return
    }

    totalPrice := 0.0
    for _, item := range items {
        food, err := queries.GetFoodById(context.Background(), item.FoodID)
        if err != nil {
            http.Error(writer, "Failed to retrieve food item", http.StatusInternalServerError)
            return
        }
        totalPrice += food.Price * float64(item.Quantity)
    }
    if account.Balance < totalPrice {
        http.Error(writer, "Insufficient balance", http.StatusPaymentRequired)
        return
    }

    // Deduct the total price from user's balance
    newBalance := account.Balance - totalPrice
    err = queries.UpdateOrderUser(context.Background(), database.UpdateOrderUserParams{
        Balance: newBalance,
        ID:      userID,
    })
    if err != nil {
        http.Error(writer, "Failed to update user balance", http.StatusInternalServerError)
        return
    }

    // Update order status to paid
    err = queries.UpdateOrderPayment(context.Background(), paymentReq.OrderID)
    
    if err != nil {
        http.Error(writer, "Failed to update order status", http.StatusInternalServerError)
        return
    }

    resp := PaymentResponse{Success: true, Message: "Payment successful"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

// Recharge Account
func RechargeAccount(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    token, err := auth.GetBearerToken(req.Header)
    if err != nil {
        http.Error(writer, "Invalid or missing token", http.StatusUnauthorized)
        return
    }

    username, userID, err := auth.ValidateJWT(token, authKey)
    if err != nil {
        http.Error(writer, "Invalid token", http.StatusUnauthorized)
        return
    }

    log.Println("Recharge request received from user:", username)

    type RechargeRequest struct {
        Amount float64 `json:"amount"`
    }
    type RechargeResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var rechargeReq RechargeRequest
    if err := json.NewDecoder(req.Body).Decode(&rechargeReq); err != nil {
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }

    if rechargeReq.Amount <= 0 {
        http.Error(writer, "Invalid recharge amount", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

    account, err := queries.GetAccount(context.Background(), username)
    
    if account.ID != userID || account.Username != username {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Update user's balance
    newBalance := account.Balance + rechargeReq.Amount
    err = queries.UpdateOrderUser(context.Background(), database.UpdateOrderUserParams{
        Balance: newBalance,
        ID:      userID,
    })
    if err != nil {
        http.Error(writer, "Failed to update user balance", http.StatusInternalServerError)
        return
    }

    resp := RechargeResponse{Success: true, Message: "Recharge successful"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

// ADMIN: GET ALL USERS
func GetAllUsers(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    token, err := auth.GetBearerToken(req.Header)
    if err != nil {
        http.Error(writer, "Invalid or missing token", http.StatusUnauthorized)
        return
    }

    username, userID, err := auth.ValidateJWT(token, authKey)
    if err != nil {
        http.Error(writer, "Invalid token", http.StatusUnauthorized)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

    admin, err := queries.GetAdminAccount(context.Background())
    if admin.ID != userID || admin.Username != username {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    users, err := queries.GetAllAccounts(context.Background())
    if err != nil {
        http.Error(writer, "Failed to retrieve users", http.StatusInternalServerError)
        return
    }

    resp := struct {
        Success bool              `json:"success"`
        Users   []database.Account `json:"users"`
        Message string            `json:"message"`
    }{
        Success: true,
        Users:   users,
        Message: "User list retrieved successfully",
    }
    
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
}

// ADMIN: GET AVERAGE SPENDING
func GetAverageSpendingAll(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    token, err := auth.GetBearerToken(req.Header)
    if err != nil {
        http.Error(writer, "Invalid or missing token", http.StatusUnauthorized)
        return
    }

    username, userID, err := auth.ValidateJWT(token, authKey)
    if err != nil {
        http.Error(writer, "Invalid token", http.StatusUnauthorized)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

    admin, err := queries.GetAdminAccount(context.Background())
    if admin.ID != userID || admin.Username != username {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    avgSpending, err := queries.GetAverageSpendingByAllUsers(context.Background())
    if err != nil {
        http.Error(writer, "Failed to retrieve average spending", http.StatusInternalServerError)
        return
    }

    resp := struct {
        Success bool   `json:"success"`
        Average float64 `json:"average"`
        Message string `json:"message"`
    }{
        Success: true,
        Average: avgSpending.(float64),
        Message: "Average spending retrieved successfully",
    }
    
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
}

// ADMIN: GET AVERAGE SPENDING BY USER
func GetAverageSpendingByUser(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    token, err := auth.GetBearerToken(req.Header)
    if err != nil {
        http.Error(writer, "Invalid or missing token", http.StatusUnauthorized)
        return
    }

    username, userID, err := auth.ValidateJWT(token, authKey)
    if err != nil {
        http.Error(writer, "Invalid token", http.StatusUnauthorized)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

    admin, err := queries.GetAdminAccount(context.Background())
    if admin.ID != userID || admin.Username != username {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    accounts, err := queries.GetAllAccounts(context.Background())
    if err != nil {
        http.Error(writer, "Failed to retrieve accounts", http.StatusInternalServerError)
        return
    }

    var avgSpendingByUser = make(map[string]float64)

    for _, account := range accounts {
        value, err := queries.GetAverageSpendingByUser(context.Background(), account.ID)
        avgSpendingByUser[account.Username] = value.(float64)
        if err != nil {
            http.Error(writer, "Failed to retrieve average spending by user", http.StatusInternalServerError)
            return
        }
    }

    resp := struct {
        Success bool   `json:"success"`
        Average map[string]float64 `json:"average"`
        Message string `json:"message"`
    }{
        Success: true,
        Average: avgSpendingByUser,
        Message: "Average spending by user retrieved successfully",
    }
    
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

func getCurrentAccount(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    token, err := auth.GetBearerToken(req.Header)
    if err != nil {
        http.Error(writer, "Invalid or missing token", http.StatusUnauthorized)
        return
    }

    username, userID, err := auth.ValidateJWT(token, authKey)
    if err != nil {
        http.Error(writer, "Invalid token", http.StatusUnauthorized)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

    account, err := queries.GetAccount(context.Background(), username)
    if err != nil || account.ID != userID {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(account)
}