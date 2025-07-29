create schema if not exists schema_name;

create table if not exists schema_name.orders
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    courier_id INTEGER NOT NULL DEFAULT -1,
    status VARCHAR(50) DEFAULT 'created'
);