-- name: CreateAccount :exec
INSERT INTO accounts (username, password, email, address, user_phone_number)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?
);

-- name: GetAccount :one
SELECT * FROM accounts WHERE username = ?;

-- name: AlterAccount :exec
UPDATE accounts
SET
    email = ?,
    address = ?,
    user_phone_number = ?
WHERE
    username = ?;

-- name: CreateFood :exec
INSERT INTO food (food_name, price, info, ingredients, time_needed, picture, description, long_range)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
);

-- name: GetFood :one
SELECT * FROM food WHERE food_name = ?;

-- name: GetFoodById :one
SELECT * FROM food WHERE food_id = ?;

-- name: GetAllFood :many
SELECT * FROM food;

-- name: GetAllFoodLongRange :many
SELECT * FROM food WHERE long_range = true;

-- name: AlterFood :exec
UPDATE food
SET
    price = ?,
    info = ?,
    ingredients = ?,
    time_needed = ?,
    description = ?,
    long_range = ?
WHERE
    food_name = ?;

-- name: DeleteFood :exec
DELETE FROM food
WHERE food_name = ?;

-- name: CreateOrder :exec
INSERT INTO orders (user_id, order_info, is_ranged, delivery_address)
VALUES (
    ?,
    ?,
    ?,
    ?
);

-- name: GetLastInsertedOrder :one
SELECT * FROM orders
WHERE order_id = LAST_INSERT_ID();

-- name: CreateOrderedItem :exec
INSERT INTO items (order_id, food_id, quantity)
VALUES (
    ?,
    ?,
    ?
);

-- name: GetOrder :many
SELECT * FROM orders WHERE user_id = ?;

-- name: GetOrderById :one
SELECT * FROM orders WHERE order_id = ?;

-- name: GetOrderedItems :many
SELECT * FROM items WHERE order_id = ?;

-- name: RateOrder :exec
UPDATE orders
SET
    feedback = ?,
    is_done = true
WHERE
    order_id = ?;

-- name: DeleteOrder :exec
UPDATE orders
SET
    deleted = true
WHERE
    order_id = ?;

-- name: GetAllOrders :many
SELECT * FROM orders WHERE deleted = false ORDER BY order_time DESC;

-- name: GetAllOrdersByUser :many
SELECT * FROM orders WHERE user_id = ? AND deleted = false;

-- name: GetAllDeletedOrdersByUser :many
SELECT * FROM orders WHERE user_id = ? AND deleted = true;

-- name: GetAllOrdersNotDone :many
SELECT * FROM orders WHERE is_done = false AND deleted = false ORDER BY order_time DESC;

-- name: UpdateOrderPayment :exec
UPDATE orders
SET
    is_paid = true
WHERE
    order_id = ?;

-- name: UpdateOrderUser :exec
UPDATE accounts
SET
    balance = ?
WHERE
    id = ?;

-- name: GetAllAccounts :many
SELECT * FROM accounts WHERE is_admin = false;

-- name: GetFoodByTag :many
SELECT DISTINCT food.*
FROM food
JOIN tags ON food.food_name = tags.food_name
WHERE tags.tag = ?;

-- name: GetMostOrderedFood :many
SELECT food.food_name, COUNT(items.food_id) AS order_count
FROM food
JOIN items ON food.food_id = items.food_id
GROUP BY food.food_name
ORDER BY order_count DESC;

-- name: GetAverageRating :one
SELECT AVG(rating) AS average_rating
FROM items
WHERE food_id = ? AND rating IS NOT NULL;

-- name: GetOrderedCountByFood :one
SELECT COUNT(*) AS ordered_count
FROM items
WHERE food_id = ?;

-- name: GetAverageSpendingByUser :one
SELECT AVG(total_price) AS average_spending
FROM (
    SELECT SUM(food.price * items.quantity) AS total_price
    FROM orders
    JOIN items ON orders.order_id = items.order_id
    JOIN food ON items.food_id = food.food_id
    WHERE orders.user_id = ?
    GROUP BY orders.order_id
) AS spending;

-- name: GetAverageSpendingByAllUsers :one
SELECT AVG(total_price) AS average_spending
FROM (
    SELECT SUM(food.price * items.quantity) AS total_price
    FROM orders
    JOIN items ON orders.order_id = items.order_id
    JOIN food ON items.food_id = food.food_id
    GROUP BY orders.user_id, orders.order_id
) AS spending;

-- name: RateFood :exec
UPDATE items
SET
    rating = ?
WHERE
    food_id = ? AND order_id = ?;

-- name: GetAllFoodTags :many
SELECT DISTINCT tag FROM tags;

-- name: GetLongestTimeNeededFoodInOrder :one
SELECT MAX(food.time_needed) AS longest_time_needed
FROM food
JOIN items ON food.food_id = items.food_id
WHERE items.order_id = ?;

-- name: UpdateEstimatedTime :exec
UPDATE orders
SET
    estimated_time = estimated_time + INTERVAL ? MINUTE
WHERE
    order_id = ?;

-- name: UpdateFeedback :exec
UPDATE orders
SET
    feedback = ?
WHERE
    order_id = ?;

-- name: GetAdminAccount :one
SELECT * FROM accounts WHERE is_admin = true;

-- name: UpdateOrderDoneStatus :exec
UPDATE orders
SET
    is_done = true
WHERE
    order_id = ?;

-- name: UpdateUserTagByID :exec
UPDATE accounts
SET
    user_tag = ?
WHERE
    id = ?;

-- name: NewTag :exec
INSERT INTO tags (tag, food_name)
VALUES (
    ?,
    ?
);

-- name: DeleteTag :exec
DELETE FROM tags
WHERE food_name = ?;

-- name: TopThreeTagByUser :many
SELECT tag, COUNT(*) AS count
FROM tags
WHERE food_name IN (
    SELECT food_name
    FROM items
    JOIN orders ON items.order_id = orders.order_id
    WHERE orders.user_id = ?
)
GROUP BY tag
ORDER BY count DESC
LIMIT 3;

-- name: GetOrderTotalPrice :one
SELECT SUM(food.price * items.quantity) AS total_price
FROM orders
JOIN items ON orders.order_id = items.order_id
JOIN food ON items.food_id = food.food_id
WHERE orders.order_id = ?;

-- name: GetMostOrderedTag :one
SELECT tag, COUNT(*) AS count
FROM tags
JOIN items ON tags.food_name = items.food_id
JOIN orders ON items.order_id = orders.order_id
WHERE orders.deleted = false
GROUP BY tag
ORDER BY count DESC
LIMIT 1;

-- name: GetFoodTagByFoodName :many
SELECT tag
FROM tags
WHERE food_name = ?;

-- name: GetAllOrderedItems :many
SELECT * FROM items;