-- +goose Up
create table items(
	item_id int auto_increment primary key,
    order_id int not null,
    food_id int not null,
    quantity int not null default 1,
    foreign key (food_id) references food(food_id) on delete cascade,
    foreign key (order_id) references orders(order_id) on delete cascade
    );

-- +goose Down
DROP TABLE items;