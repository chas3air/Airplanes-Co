CREATE DATABASE airplanesco;

\c airplanesco

CREATE TABLE IF NOT EXISTS Customers (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50),
    surname VARCHAR(100),
    name VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS Flights (
    id SERIAL PRIMARY KEY,
    from_where VARCHAR(50) NOT NULL,
    destination VARCHAR(50) NOT NULL,
    flight_time TIMESTAMP NOT NULL,
    flight_duration INT NOT NULL
);

-- CREATE TABLE IF NOT EXISTS Tickets (
--     id SERIAL PRIMARY KEY,
--     flightId INT NOT NULL,
--     ownerId INT NOT NULL,
--     ticketCost NUMERIC(10, 2) NOT NULL,
--     classOfService VARCHAR(50) NOT NULL,
--     CONSTRAINT fk_flight
--         FOREIGN KEY (flightId) 
--         REFERENCES Flights(id)
--         ON DELETE CASCADE,
--     CONSTRAINT fk_owner
--         FOREIGN KEY (ownerId) 
--         REFERENCES Customers(id)
--         ON DELETE CASCADE
-- );