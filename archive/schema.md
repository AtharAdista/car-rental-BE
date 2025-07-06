-- Table: customers
CREATE TABLE customers_v1 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    nik VARCHAR(20) UNIQUE NOT NULL,
    phone_number VARCHAR(20) NOT NULL
);

-- Table: cars
CREATE TABLE cars_v1 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    stock INT NOT NULL,
    daily_rent NUMERIC(12, 2) NOT NULL
);

-- Table: bookings
CREATE TABLE bookings_v1 (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    cars_id INT NOT NULL,
    start_rent DATE NOT NULL,
    end_rent DATE NOT NULL,
    total_cost NUMERIC(12, 2) NOT NULL,
    finished BOOLEAN DEFAULT FALSE,

    CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers_v1(id) ON DELETE CASCADE,
    CONSTRAINT fk_car FOREIGN KEY (cars_id) REFERENCES cars_v1(id) ON DELETE CASCADE
);

-- V2

CREATE TABLE cars_v2 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    stock INT NOT NULL,
    daily_rent NUMERIC(12, 2) NOT NULL
);

CREATE TABLE drivers_v2 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    nik VARCHAR(50) NOT NULL UNIQUE,
    phone_number VARCHAR(20) NOT NULL,
    daily_cost DOUBLE PRECISION NOT NULL
);

CREATE TABLE memberships_v2 (
    id SERIAL PRIMARY KEY,
    membership_name VARCHAR(100) NOT NULL,
    discount INT NOT NULL
);

CREATE TABLE booking_types_v2 (
    id SERIAL PRIMARY KEY,
    booking_type VARCHAR(100) NOT NULL,
    description TEXT
);

CREATE TABLE customers_v2 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    nik VARCHAR(50) NOT NULL UNIQUE,
    phone_number VARCHAR(20) NOT NULL,
    membership_id INT,
    CONSTRAINT fk_membership
        FOREIGN KEY (membership_id)
        REFERENCES memberships_v2(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE bookings_v2 (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    cars_id INT NOT NULL,
    booking_type_id INT NOT NULL,
    driver_id INT NOT NULL,
    start_rent TIMESTAMP NOT NULL,
    end_rent TIMESTAMP NOT NULL,
    total_cost NUMERIC(12, 2) NOT NULL,
    finished BOOLEAN NOT NULL DEFAULT FALSE,
    discount NUMERIC(12, 2) NOT NULL,
    total_driver_cost NUMERIC(12, 2) NOT NULL,

    CONSTRAINT fk_customer
        FOREIGN KEY (customer_id)
        REFERENCES customers_v2(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_car
        FOREIGN KEY (cars_id)
        REFERENCES cars_v2(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_booking_type
        FOREIGN KEY (booking_type_id)
        REFERENCES booking_types_v2(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_driver
        FOREIGN KEY (driver_id)
        REFERENCES drivers_v2(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE drivers_incentives_v2 (
    id SERIAL PRIMARY KEY,
    booking_id INT NOT NULL,
    incentive NUMERIC(12, 2) NOT NULL,
    
    CONSTRAINT fk_booking
        FOREIGN KEY (booking_id)
        REFERENCES bookings_v2(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);