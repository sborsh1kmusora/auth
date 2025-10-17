-- +goose Up
create table accesses (
  id serial primary key,
  endpoint_addresses text not null,
  role text not null,
  created_at timestamp not null default clock_timestamp(),
  updated_at timestamp
);

-- +goose Down
drop table accesses;
