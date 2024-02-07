-- +goose Up
-- +goose StatementBegin
CREATE TABLE GaugeMetrics (
    key text unique not null primary key,
    val double precision

--     Alloc double precision,
--     BuckHashSys double precision,
--     Frees double precision,
--     GCCPUFraction double precision,
--     GCSys double precision,
--     HeapAlloc double precision,
--     HeapIdle double precision,
--     HeapInuse double precision,
--     HeapObjects double precision,
--     HeapReleased double precision,
--     HeapSys double precision,
--     LastGC double precision,
--     Lookups double precision,
--     MCacheInuse double precision,
--     MCacheSys double precision,
--     MSpanInuse double precision,
--     MSpanSys double precision,
--     Mallocs double precision,
--     NextGC double precision,
--     NumForcedGC double precision,
--     NumGC double precision,
--     OtherSys double precision,
--     PauseTotalNs double precision,
--
--     StackInuse double precision,
--     StackSys double precision,
--     Sys double precision,
--     TotalAlloc double precision
--     RandomValue double precision,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists GaugeMetrics;
-- +goose StatementEnd
