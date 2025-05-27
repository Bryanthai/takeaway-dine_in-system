-- +goose Up
create table accounts(
	id int auto_increment primary key,
    username varchar(30) unique NOT NULL,
    password varchar(100) NOT NULL,
    email varchar(100) unique NOT NULL,
    address varchar(255) NOT NULL,
    user_tag varchar(100),
    user_phone_number BIGINT NOT NULL
    );

-- +goose Down
DROP TABLE accounts;