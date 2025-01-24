CREATE TABLE delivery_audits (
    id SERIAL PRIMARY KEY,
    delivery_id BIGINT NOT NULL,
    provider VARCHAR(255) NOT NULL,
    status delivery_status NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);