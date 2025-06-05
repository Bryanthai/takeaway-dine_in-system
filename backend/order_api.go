package main

import(
	"database/sql"
	"net/http"
	"encoding/json"
	"context"
    "log"
    "fmt"

	"github.com/Bryanthai/ordersystem/internal/database"
    "github.com/Bryanthai/ordersystem/internal/auth"
)

func updateUserTag(userID int32) error {
    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        return fmt.Errorf("database error: %w", err)
    }
    defer db.Close()
    queries := database.New(db)

    tags, err := queries.TopThreeTagByUser(context.Background(), userID)
    if err != nil {
        return fmt.Errorf("failed to get top tags: %w", err)
    }

    var NewUserTag string
    for _, tag := range tags {
        NewUserTag += tag.Tag + " "
    }

    fmt.Println("changing user tag to:", NewUserTag)

    err = queries.UpdateUserTagByID(context.Background(), database.UpdateUserTagByIDParams{
        UserTag: sql.NullString{
            String: NewUserTag,
            Valid:  NewUserTag != "",
        },
        ID:      userID,
    })
    if err != nil {
        return fmt.Errorf("failed to update user tag: %w", err)
    }
    log.Println("User tag updated successfully for user ID:", userID)
    return nil
}

// CREATE NEW ORDER
func createOrderHandler(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Create order request received from user:", username)

    type CreateOrderRequest struct {
        OrderInfo string `json:"order_info"`
        IsRanged  bool   `json:"is_ranged"`
        OrderItems []struct {
            FoodID   int32 `json:"food_id"`
            Quantity int32 `json:"quantity"`
        } `json:"order_items"`
        DeliveryAddress string `json:"delivery_address"`
    }
    type CreateOrderResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
        OrderID int32  `json:"order_id"`
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

    account, err := queries.GetAccount(context.Background(), username)
    if account.ID != userID || err != nil {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    err = queries.CreateOrder(context.Background(), database.CreateOrderParams{
        UserID:    userID,
        OrderInfo: orderReq.OrderInfo,
        IsRanged:  orderReq.IsRanged,
        DeliveryAddress: sql.NullString{
            String: orderReq.DeliveryAddress,
            Valid:  orderReq.DeliveryAddress != "",
        },
    })
    if err != nil {
        http.Error(writer, "Failed to create order", http.StatusInternalServerError)
        return
    }

    NewOrder, err := queries.GetLastInsertedOrder(context.Background())
    if err != nil {
        http.Error(writer, "Failed to retrieve last inserted order", http.StatusInternalServerError)
        return
    }

    for _, item := range orderReq.OrderItems {
        err = queries.CreateOrderedItem(context.Background(), database.CreateOrderedItemParams{
            OrderID:  NewOrder.OrderID,
            FoodID:   item.FoodID,
            Quantity: item.Quantity,
        })
        if err != nil {
            http.Error(writer, "Failed to create order item", http.StatusInternalServerError)
            return
        }
    }

    TimeNeeded, err := queries.GetLongestTimeNeededFoodInOrder(context.Background(), NewOrder.OrderID)
    if err != nil {
        http.Error(writer, "Failed to get time needed for order", http.StatusInternalServerError)
        return
    }

    estimatedTime := TimeNeeded.(int64)
    
    err = queries.UpdateEstimatedTime(context.Background(), database.UpdateEstimatedTimeParams{
        DATEADD:       estimatedTime,
        OrderID:       NewOrder.OrderID,
    })
    if err != nil {
        http.Error(writer, "Failed to update estimated time", http.StatusInternalServerError)
        return
    }

    err = updateUserTag(userID)
    if err != nil {
        http.Error(writer, "Failed to update user tag", http.StatusInternalServerError)
        return
    }

    resp := CreateOrderResponse{Success: true, Message: "Order created successfully", OrderID: NewOrder.OrderID}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
	return
}

// GET ORDER STATUS
func getOrderStatusHandler(writer http.ResponseWriter, req *http.Request) {
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

    UserID := userID

    log.Println("Get order status request received from user:", username)

    type GetOrderStatusResponse struct {
        Success bool        		`json:"success"`
        Orders  []database.Order 	`json:"orders"`
        Message string      		`json:"message"`
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

    account, err := queries.GetAccount(context.Background(), username)
    if account.ID != userID || err != nil {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }
    orders, err := queries.GetOrder(context.Background(), UserID)
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

// DELETE ORDER
func deleteOrderHandler(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Delete order request received from user:", username)

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
    
    admin, err := queries.GetAdminAccount(context.Background())
    if admin.ID != userID || admin.Username != username {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }
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

// GET ALL ORDERS
func getAllOrdersHandler(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Get all orders request received from user:", username)

    type GetAllOrdersResponse struct {
        Success bool              `json:"success"`
        Orders  []database.Order  `json:"orders"`
        Message string            `json:"message"`
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
    
    orders, err := queries.GetAllOrders(context.Background())
    if err != nil {
        http.Error(writer, "Failed to get orders", http.StatusInternalServerError)
        return
    }

    resp := GetAllOrdersResponse{
        Success: true,
        Orders:  orders,
        Message: "Orders retrieved successfully",
    }
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

// RATE ORDERED ITEMS
func rateOrderedItems(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Rate ordered items request received from user:", username)

    type RateItemsRequest struct {
        OrderID int32 `json:"order_id"`
        FoodID  int32 `json:"food_id"`
        Rating  int32 `json:"rating"`
    }
    type RateItemsResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var rateReq RateItemsRequest
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
    account, err := queries.GetAccount(context.Background(), username)
    if account.ID != userID || err != nil {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    err = queries.RateFood(context.Background(), database.RateFoodParams{
        OrderID: rateReq.OrderID,
        FoodID:  rateReq.FoodID,
        Rating:  sql.NullInt32{
            Int32: rateReq.Rating,
            Valid: rateReq.Rating > 0,
        },
    })
    if err != nil {
        http.Error(writer, "Failed to rate ordered items", http.StatusInternalServerError)
        return
    }
    resp := RateItemsResponse{Success: true, Message: "Items rated successfully"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

// Update Feedback
func UpdateFeedback(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Update feedback request received from user:", username)

    type UpdateFeedbackRequest struct {
        OrderID  int32          `json:"order_id"`
        Feedback string `json:"feedback"`
    }
    type UpdateFeedbackResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var feedbackReq UpdateFeedbackRequest
    if err := json.NewDecoder(req.Body).Decode(&feedbackReq); err != nil {
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
    if account.ID != userID || err != nil {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    err = queries.UpdateFeedback(context.Background(), database.UpdateFeedbackParams{
        OrderID:  feedbackReq.OrderID,
        Feedback: sql.NullString{
            String: feedbackReq.Feedback,
            Valid:  feedbackReq.Feedback != "",
        },
    })
    
    if err != nil {
        http.Error(writer, "Failed to update feedback", http.StatusInternalServerError)
        return
    }

    resp := UpdateFeedbackResponse{Success: true, Message: "Feedback updated successfully"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

// GET ORDERED ITEMS
func getOrderedItemsHandler(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Get ordered items request received from user:", username)

    type GetOrderedItemsResponse struct {
        Success bool              `json:"success"`
        Items   []database.Item   `json:"items"`
        Message string            `json:"message"`
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

    account, err := queries.GetAccount(context.Background(), username)
    if account.ID != userID || err != nil {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    orderIDStr := req.URL.Query().Get("order_id")
    if orderIDStr == "" {
        http.Error(writer, "Missing order_id query parameter", http.StatusBadRequest)
        return
    }
    var orderID int32
    if _, err := fmt.Sscanf(orderIDStr, "%d", &orderID); err != nil {
        http.Error(writer, "Invalid order_id", http.StatusBadRequest)
        return
    }

    items, err := queries.GetOrderedItems(context.Background(), orderID)
    if err != nil {
        http.Error(writer, "Failed to get ordered items", http.StatusInternalServerError)
        return
    }

    resp := GetOrderedItemsResponse{
        Success: true,
        Items:   items,
        Message: "Ordered items retrieved successfully",
    }
    
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

func finishOrder(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Finish order request received from user:", username)

    type FinishOrderRequest struct {
        OrderID int32 `json:"order_id"`
    }
    
    type FinishOrderResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    var finishReq FinishOrderRequest
    if err := json.NewDecoder(req.Body).Decode(&finishReq); err != nil {
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
    if account.ID != userID || err != nil {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    err = queries.UpdateOrderDoneStatus(context.Background(), finishReq.OrderID)
    if err != nil {
        http.Error(writer, "Failed to finish order", http.StatusInternalServerError)
        return
    }

    writer.WriteHeader(http.StatusOK)
    resp := FinishOrderResponse{Success: true, Message: "Order finished successfully"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

func getAllUndoneOrder(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Get all undone orders request received from user:", username)

    type GetAllUndoneOrdersResponse struct {
        Success bool              `json:"success"`
        Orders  []database.Order  `json:"orders"`
        Message string            `json:"message"`
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

    account, err := queries.GetAccount(context.Background(), username)
    if account.ID != userID || err != nil {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    orders, err := queries.GetAllOrdersNotDone(context.Background())
    if err != nil {
        http.Error(writer, "Failed to get undone orders", http.StatusInternalServerError)
        return
    }

    resp := GetAllUndoneOrdersResponse{
        Success: true,
        Orders:  orders,
        Message: "Undone orders retrieved successfully",
    }
    
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
}

func GetOrderByIdHandler(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Get order by ID request received from user:", username)

    type GetOrderByIdResponse struct {
        Success bool              `json:"success"`
        Order   database.Order    `json:"order"`
        Message string            `json:"message"`
    }

    orderIDStr := req.URL.Query().Get("order_id")
    if orderIDStr == "" {
        http.Error(writer, "Missing order_id query parameter", http.StatusBadRequest)
        return
    }
    var orderID int32
    if _, err := fmt.Sscanf(orderIDStr, "%d", &orderID); err != nil {
        http.Error(writer, "Invalid order_id", http.StatusBadRequest)
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
    if account.ID != userID || err != nil {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    order, err := queries.GetOrderById(context.Background(), orderID)
    if err != nil {
        http.Error(writer, "Failed to get order by ID", http.StatusInternalServerError)
        return
    }

    resp := GetOrderByIdResponse{
        Success: true,
        Order:   order,
        Message: "Order retrieved successfully",
    }
    
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
}

func GetOrderTotalPrice(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Get order total price request received from user:", username)

    type GetOrderTotalPriceResponse struct {
        Success bool   `json:"success"`
        Total   float64 `json:"total"`
        Message string `json:"message"`
    }

    orderIDStr := req.URL.Query().Get("order_id")
    if orderIDStr == "" {
        http.Error(writer, "Missing order_id query parameter", http.StatusBadRequest)
        return
    }
    var orderID int32
    if _, err := fmt.Sscanf(orderIDStr, "%d", &orderID); err != nil {
        http.Error(writer, "Invalid order_id", http.StatusBadRequest)
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
    if account.ID != userID || err != nil {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    totalPrice, err := queries.GetOrderTotalPrice(context.Background(), orderID)
    if err != nil {
        http.Error(writer, "Failed to get order total price", http.StatusInternalServerError)
        return
    }

    resp := GetOrderTotalPriceResponse{
        Success: true,
        Total:   totalPrice.(float64),
        Message: "Order total price retrieved successfully",
    }
    
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

func GetAllOrderedItemsHandler(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Get all ordered items request received from user:", username)

    type GetAllOrderedItemsResponse struct {
        Success bool              `json:"success"`
        Items   []database.Item   `json:"items"`
        Message string            `json:"message"`
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

    account, err := queries.GetAccount(context.Background(), username)
    if account.ID != userID || err != nil {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

    items, err := queries.GetAllOrderedItems(context.Background())
    if err != nil {
        http.Error(writer, "Failed to get ordered items", http.StatusInternalServerError)
        return
    }

    resp := GetAllOrderedItemsResponse{
        Success: true,
        Items:   items,
        Message: "Ordered items retrieved successfully",
    }
    
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}