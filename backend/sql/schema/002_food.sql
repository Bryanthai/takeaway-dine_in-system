-- +goose Up
create table food(
	food_id int auto_increment primary key,
    food_name varchar(30) unique NOT NULL,
    price double(5,2) NOT NULL,
    food_type varchar(30) NOT NULL,
    picture varchar(100),
    long_range boolean NOT NULL DEFAULT false,
    description varchar(255) NOT NULL,
    info varchar(100),
    ingredients varchar(100) NOT NULL,
    time_needed INT not null
    );

-- +goose Down
DROP TABLE food;