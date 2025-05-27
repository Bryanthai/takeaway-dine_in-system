-- +goose Up
create table orders(
	order_id int auto_increment primary key,
    user_id int not null,
    foreign key (user_id) references accounts(id) on delete cascade,
    order_info varchar(255) not null,
    rating int default null check (rating >= 0 and rating <= 5),
    feedback text default null,
    order_time timestamp default current_timestamp,
    is_done bool default false not null,
    is_ranged bool default false not null,
    deleted bool default false not null
    );

-- +goose Down
DROP TABLE orders;