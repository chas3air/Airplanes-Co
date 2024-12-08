-- Создание базы данных
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

CREATE TABLE IF NOT EXISTS City (
    id SERIAL PRIMARY KEY,
    name VARCHAR(15) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Airplanes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS Flights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fromwhere INTEGER NOT NULL REFERENCES City(id) ON DELETE CASCADE,
    destination INTEGER NOT NULL REFERENCES City(id) ON DELETE CASCADE,
    flighttime TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    flightduration INTEGER NOT NULL,
    airplaneid UUID NOT NULL REFERENCES Airplanes(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS CostTable (
    id SERIAL PRIMARY KEY NOT NULL,
    flightId UUID NOT NULL,
    className VARCHAR(15) NOT NULL,
    Cost INT NOT NULL,
    UNIQUE (flightId, className),
    FOREIGN KEY (flightId) REFERENCES Flights(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flightId UUID NOT NULL REFERENCES Flights(id) ON DELETE CASCADE,
    ownerId UUID NOT NULL REFERENCES Customers(id) ON DELETE CASCADE,
    classOfService VARCHAR(50) NOT NULL
);