-- +goose Up
create table tags(
    tag varchar(100) not null,
    food_name varchar(255) not null,
    foreign key (food_name) references food(food_name) on delete cascade
);

-- +goose Down
DROP TABLE tags;