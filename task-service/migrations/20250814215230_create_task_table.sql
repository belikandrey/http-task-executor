-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS task
(
    id                   SERIAL PRIMARY KEY,
    url                  TEXT        NOT NULL,
    method               TEXT        NOT NULL,
    status               VARCHAR(40) NOT NULL,
    response_status_code SMALLINT,
    response_length      BIGINT
);

CREATE TABLE headers
(
    id      SERIAL PRIMARY KEY,
    name    TEXT NOT NULL,
    value   TEXT NOT NULL,
    task_id BIGINT REFERENCES task (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS headers;
DROP TABLE IF EXISTS task;
-- +goose StatementEnd
