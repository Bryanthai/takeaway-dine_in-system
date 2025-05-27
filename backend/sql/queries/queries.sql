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
INSERT INTO food (food_name, food_tag, price, info, ingredients, time_needed)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
);

-- name: GetFood :one
SELECT * FROM food WHERE food_name = ?;

-- name: GetAllFood :many
SELECT * FROM food;

-- name: AlterFood :exec
UPDATE food
SET
    food_tag = ?,
    price = ?,
    info = ?,
    ingredients = ?,
    time_needed = ?
WHERE
    food_name = ?;

-- name: DeleteFood :exec
DELETE FROM food
WHERE food_name = ?;

-- name: CreateOrder :exec
INSERT INTO orders (user_id, order_info, is_ranged)
VALUES (
    ?,
    ?,
    ?
);

-- name: CreateOrderedItem :exec
INSERT INTO items (order_id, food_id, quantity)
VALUES (
    ?,
    ?,
    ?
);

-- name: GetOrder :many
SELECT * FROM orders WHERE user_id = ?;

-- name: GetOrderedItems :many
SELECT * FROM items WHERE order_id = ?;

-- name: RateOrder :exec
UPDATE orders
SET
    rating = ?,
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