CREATE TYPE order_status AS ENUM ('created', 'in_progress', 'completed', 'canceled');

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(255) UNIQUE NOT NULL,
    status order_status NOT NULL,
    user_info JSON NOT NULL,
    origin_latitude DOUBLE PRECISION NOT NULL,
    origin_longitude DOUBLE PRECISION NOT NULL,
    destination_latitude DOUBLE PRECISION NOT NULL,
    destination_longitude DOUBLE PRECISION NOT NULL,
    time_frame_days INT NOT NULL,
    time_frame_from TIMESTAMP NOT NULL,
    time_frame_to TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL
);