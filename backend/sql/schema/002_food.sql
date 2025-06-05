-- +goose Up
create table food(
	food_id int auto_increment primary key,
    food_name varchar(255) unique NOT NULL,
    price double(5,2) NOT NULL,
    picture varchar(100),
    long_range boolean NOT NULL DEFAULT false,
    description varchar(600) NOT NULL,
    info varchar(100),
    ingredients varchar(600) NOT NULL,
    time_needed INT not null
    );

-- +goose Down
DROP TABLE food;