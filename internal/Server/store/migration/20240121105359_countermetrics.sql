-- +goose Up
-- +goose StatementBegin
CREATE TABLE CounterMetrics (
    key text unique not null primary key,
    val bigint
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists CounterMetrics;
-- +goose StatementEnd
