CREATE TYPE delivery_status AS ENUM ('init', 'finding', 'found', 'not_found', 'delivered', 'canceled');

CREATE TABLE deliveries (
    id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    provider VARCHAR(255) NOT NULL,
    origin_latitude DOUBLE PRECISION NOT NULL,
    origin_longitude DOUBLE PRECISION NOT NULL,
    destination_latitude DOUBLE PRECISION NOT NULL,
    destination_longitude DOUBLE PRECISION NOT NULL,
    time_frame_start TIMESTAMP NOT NULL,
    time_frame_end TIMESTAMP NOT NULL,
    status delivery_status NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL
);