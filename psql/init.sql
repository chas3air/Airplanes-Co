CREATE DATABASE airplanesco;

\c airplanesco

CREATE TABLE IF NOT EXISTS Customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS Flights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fromWhere VARCHAR(50) NOT NULL,
    destination VARCHAR(50) NOT NULL,
    flightTime TIMESTAMP NOT NULL,
    flightDuration INT NOT NULL,
    flightSeatsCost INTEGER[] NOT NULL
);

CREATE TABLE IF NOT EXISTS Tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flightId UUID NOT NULL,
    ownerId UUID NOT NULL,
    ticketCost NUMERIC(10, 2) NOT NULL,
    classOfService VARCHAR(50) NOT NULL
);
