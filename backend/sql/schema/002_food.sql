-- +goose Up
create table food(
	food_id int auto_increment primary key,
    food_name varchar(30) unique NOT NULL,
    food_tag varchar(30) NOT NULL,
    price double(5,2) NOT NULL,
    info varchar(100),
    ingredients varchar(100) NOT NULL,
    time_needed varchar(50)
    );

-- +goose Down
DROP TABLE food;