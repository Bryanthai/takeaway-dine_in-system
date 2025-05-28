-- +goose Up
create table adminTable(
    username varchar(50) not null primary key,
    id int not null,
    password varchar(255) not null
    );

-- +goose Down
DROP TABLE adminTable;