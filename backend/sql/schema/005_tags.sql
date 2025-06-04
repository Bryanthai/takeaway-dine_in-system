-- +goose Up
create table tags(
    tag varchar(50) not null,
    food_name varchar(30) not null,
    foreign key (food_name) references food(food_name) on delete cascade
);

-- +goose Down
DROP TABLE tags;