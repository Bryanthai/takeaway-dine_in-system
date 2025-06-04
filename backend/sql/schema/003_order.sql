-- +goose Up
create table orders(
	order_id int auto_increment primary key,
    user_id int not null,
    foreign key (user_id) references accounts(id) on delete cascade,
    order_info varchar(255) not null,
    feedback text default null,
    order_time timestamp default current_timestamp,
    estimated_time timestamp not null default current_timestamp,
    is_done bool default false not null,
    is_ranged bool default false not null,
    delivery_address varchar(255) default null,
    deleted bool default false not null,
    is_paid bool default false not null
    );

-- +goose Down
DROP TABLE orders;