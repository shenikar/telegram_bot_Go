-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    requests (
        id SERIAL PRIMARY KEY,
        user_id BIGINT NOT NULL,
        hash TEXT NOT NULL,
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            result TEXT
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE requests;

-- +goose StatementEnd