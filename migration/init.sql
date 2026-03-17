CREATE TABLE IF NOT EXISTS shipments (
    id UUID PRIMARY KEY,
    reference_number VARCHAR(255) NOT NULL,
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    current_status VARCHAR(50) NOT NULL,
    driver_details TEXT,
    unit_details TEXT,
    shipment_amount NUMERIC(10, 2),
    driver_revenue NUMERIC(10, 2),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS shipment_events (
    id UUID PRIMARY KEY,
    shipment_id UUID NOT NULL REFERENCES shipments(id),
    previous_status VARCHAR(50) NOT NULL,
    new_status VARCHAR(50) NOT NULL,
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);