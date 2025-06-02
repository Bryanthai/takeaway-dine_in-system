package main

import(
	"database/sql"
	"net/http"
	"encoding/json"
	"context"
    "log"
    "fmt"

	"github.com/Bryanthai/ordersystem/internal/database"
)

const dbURL = "hahant:123456@tcp(localhost:3306)/fooddb?parseTime=true&tls=false"
const authKey = "kp7Nbll8hKieGO2L1EQyOphkwJDVH0xD/tOp+DssAJ0Ll1t+jFnqdxK2BOmrAlMzQehX/2v4nt1xdDyXuU0CQQ=="

func enableCORS(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// SHOW MENU
func getAllFoodHandler(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    type GetAllFoodResponse struct {
        Success bool              `json:"success"`
        Foods   []database.Food   `json:"foods"`
        Message string            `json:"message"`
    }

    log.Println("Get all food request received")

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

// GET FOOD RATING AND ORDERED TIMES
func getFoodRatingandOrderedTimesByFoodID (writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    type GetFoodRatingResponse struct {
        Success bool   `json:"success"`
        Rating  float64  `json:"rating"`
        Times   int64  `json:"times_ordered"`
        Message string `json:"message"`
    }
    
    foodIDStr := req.URL.Query().Get("food_id")
    if foodIDStr == "" {
        http.Error(writer, "Missing food_id query parameter", http.StatusBadRequest)
        return
    }
    var FoodID int32
    if _, err := fmt.Sscanf(foodIDStr, "%d", &FoodID); err != nil {
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
    ratingRaw, err := queries.GetAverageRating(context.Background(), FoodID)
    if err != nil {
        http.Error(writer, "Failed to get food rating", http.StatusInternalServerError)
        return
    }

    var rating float64
    if ratingRaw != nil {
        rating = ratingRaw.(float64)
        if rating < 0 {
            http.Error(writer, "Invalid food rating", http.StatusInternalServerError)
            return
        }
    } else {
        rating = 0.0
    }

    times, err := queries.GetOrderedCountByFood(context.Background(), FoodID)

    if err != nil {
        http.Error(writer, "Failed to get ordered times", http.StatusInternalServerError)
        return
    }

    resp := GetFoodRatingResponse{
        Success: true,
        Rating:  rating,
        Times:   times,
        Message: "Food rating and ordered times retrieved successfully",
    }
    
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}

// GET FOOD BY TYPE
func getFoodByTypeHandler(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    type FoodByType struct {
        Tag   string           `json:"tag"`
        Foods []database.Food  `json:"foods"`
    }
    type GetFoodByTypeResponse struct {
        Success bool         `json:"success"`
        Data    []FoodByType `json:"data"`
        Message string       `json:"message"`
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)

    // Get all unique tags
    tags, err := queries.GetAllFoodTags(context.Background())
    if err != nil {
        http.Error(writer, "Failed to get food tags", http.StatusInternalServerError)
        return
    }

    var result []FoodByType
    for _, tag := range tags {
        foods, err := queries.GetFoodByTag(context.Background(), tag)
        if err != nil {
            http.Error(writer, "Failed to get foods by tag", http.StatusInternalServerError)
            return
        }
        result = append(result, FoodByType{
            Tag:   tag,
            Foods: foods,
        })
    }

    resp := GetFoodByTypeResponse{
        Success: true,
        Data:    result,
        Message: "Foods grouped by type retrieved successfully",
    }
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
    return
}