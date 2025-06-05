package main

import(
	"database/sql"
	"net/http"
	"encoding/json"
	"context"
    "log"
    "fmt"
    "io"
    "os"

	"github.com/Bryanthai/ordersystem/internal/database"
    "github.com/Bryanthai/ordersystem/internal/auth"
)

// CREATE NEW FOOD
func createFoodHandler(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Create food request received from user:", username)

    type CreateFoodResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }

    err = req.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(writer, "Invalid form data", http.StatusBadRequest)
        return
    }

    foodName := req.FormValue("food_name")
    foodTag := req.FormValue("food_tag")
    priceStr := req.FormValue("price")
    info := req.FormValue("info")
    ingredients := req.FormValue("ingredients")
    timeNeededStr := req.FormValue("time_needed")
    description := req.FormValue("description")
    longRange := req.FormValue("long_range")

    var price float64
    var timeNeeded int32
    if priceStr != "" {
        _, err = fmt.Sscanf(priceStr, "%f", &price)
        if err != nil {
            http.Error(writer, "Invalid price", http.StatusBadRequest)
            return
        }
    }
    if timeNeededStr != "" {
        _, err = fmt.Sscanf(timeNeededStr, "%d", &timeNeeded)
        if err != nil {
            http.Error(writer, "Invalid time_needed", http.StatusBadRequest)
            return
        }
    }

    var imageFilename string
    file, handler, err := req.FormFile("image")
    if err == nil && handler != nil {
        imageFilename = "http://localhost:8080/images/" + handler.Filename
        defer file.Close()
        out, _ := os.Create("./static/images/food/" + handler.Filename)
        io.Copy(out, file)
        file.Close()
    } else {
        imageFilename = ""
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

    var infoStruct sql.NullString
    infoStruct.String = info
    infoStruct.Valid = info != ""

    log.Println(longRange)

    err = queries.CreateFood(context.Background(), database.CreateFoodParams{
        FoodName:    foodName,
        Price:       price,
        Picture:     sql.NullString{String: imageFilename, Valid: imageFilename != ""},
        Info:        infoStruct,
        Ingredients: ingredients,
        TimeNeeded:  timeNeeded,
        Description: description,
        LongRange:   longRange == "true",
    })

    tagSlice := removeWhiteSpace(foodTag)
    if len(tagSlice) == 0 {
        http.Error(writer, "Food tag cannot be empty", http.StatusBadRequest)
        return
    }
    for _, tag := range tagSlice {
        err = queries.NewTag(context.Background(), database.NewTagParams{
            Tag:      tag,
            FoodName: foodName,
        })
        if err != nil {
            http.Error(writer, "Failed to create food tag", http.StatusInternalServerError)
            return
        }
    }

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

    log.Println("Alter food request received from user:", username)

    type AlterFoodRequest struct {
        FoodName    string  `json:"food_name"`
        FoodTag     string  `json:"food_tag"`
        Price       float64 `json:"price"`
        Info        string  `json:"info"`
        Ingredients string  `json:"ingredients"`
        TimeNeeded  int32 `json:"time_needed"`
        LongRange   bool    `json:"long_range"`
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

    admin, err := queries.GetAdminAccount(context.Background())
    if admin.ID != userID || admin.Username != username {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }

	var infoStruct sql.NullString
	infoStruct.String = foodReq.Info
	infoStruct.Valid = foodReq.Info != ""

    err = queries.AlterFood(context.Background(), database.AlterFoodParams{
        Price:       foodReq.Price,
        Info:        infoStruct,
        Ingredients: foodReq.Ingredients,
        TimeNeeded:  foodReq.TimeNeeded,
        FoodName:    foodReq.FoodName,
        LongRange:   foodReq.LongRange,
    })
    if err != nil {
        http.Error(writer, "Failed to update food item", http.StatusInternalServerError)
        return
    }

    tagSlice := removeWhiteSpace(foodReq.FoodTag)
    if len(tagSlice) == 0 {
        http.Error(writer, "Food tag cannot be empty", http.StatusBadRequest)
        return
    }

    err = queries.DeleteTag(context.Background(), foodReq.FoodName)
    if err != nil {
        http.Error(writer, "Failed to delete existing food tags", http.StatusInternalServerError)
        return
    }

    for _, tag := range tagSlice {
        err = queries.NewTag(context.Background(), database.NewTagParams{
            Tag:      tag,
            FoodName: foodReq.FoodName,
        })

        if err != nil {
            http.Error(writer, "Failed to update food tag", http.StatusInternalServerError)
            return
        }
    }

    resp := AlterFoodResponse{Success: true, Message: "Food item updated successfully"}
    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(resp)
	return
}

// DELETE FOOD
func deleteFoodHandler(writer http.ResponseWriter, req *http.Request) {
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

    log.Println("Delete food request received from user:", username)

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
    admin, err := queries.GetAdminAccount(context.Background())
    if admin.ID != userID || admin.Username != username {
        http.Error(writer, "Unauthorized", http.StatusUnauthorized)
        return
    }
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

func getFoodByIdHandler(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    foodID := req.URL.Query().Get("food_id")
    if foodID == "" {
        http.Error(writer, "Missing food id", http.StatusBadRequest)
        return
    }
    var FoodID int32
    if _, err := fmt.Sscanf(foodID, "%d", &FoodID); err != nil {
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
    food, err := queries.GetFoodById(context.Background(), FoodID)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(writer, "Food not found", http.StatusNotFound)
        } else {
            http.Error(writer, "Failed to fetch food item", http.StatusInternalServerError)
        }
        return
    }

    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(food)
}

func GetFoodTagByFoodNameHandler(writer http.ResponseWriter, req *http.Request) {
    enableCORS(writer)
    if req.Method == "OPTIONS" {
        writer.WriteHeader(http.StatusOK)
        return
    }

    foodName := req.URL.Query().Get("food_name")
    if foodName == "" {
        http.Error(writer, "Missing food name", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    queries := database.New(db)
    tags, err := queries.GetFoodTagByFoodName(context.Background(), foodName)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(writer, "No tags found for this food item", http.StatusNotFound)
        } else {
            http.Error(writer, "Failed to fetch food tags", http.StatusInternalServerError)
        }
        return
    }

    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(tags)
}