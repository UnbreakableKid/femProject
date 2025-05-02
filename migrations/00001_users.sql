-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT null,
    email VARCHAR(255) UNIQUE NOT null,
    password_hash VARCHAR(255) NOT null,
    bio text,
    created_at TIMESTAMP
    with
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    with
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;

-- +goose StatementEnd
