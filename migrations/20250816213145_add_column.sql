-- +goose Up
-- +goose StatementBegin
ALTER TABLE headers
    ADD COLUMN input BOOLEAN DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE headers
    DROP COLUMN input;
-- +goose StatementEnd
