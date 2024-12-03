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

CREATE TABLE IF NOT EXISTS Customers_Deleted (
    id UUID PRIMARY KEY,
    login VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    deleted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_by VARCHAR(50) NOT NULL
);

CREATE OR REPLACE FUNCTION move_customer_to_deleted() 
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO Customers_Deleted (id, login, password, role, surname, name, deleted_at, deleted_by)
    VALUES (OLD.id, OLD.login, OLD.password, OLD.role, OLD.surname, OLD.name, NOW(), current_user);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_customers_delete
AFTER DELETE ON Customers
FOR EACH ROW
EXECUTE FUNCTION move_customer_to_deleted();

-- тут другая бд

CREATE TABLE IF NOT EXISTS Flights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fromWhere VARCHAR(50) NOT NULL,
    destination VARCHAR(50) NOT NULL,
    flightTime TIMESTAMP NOT NULL,
    flightDuration INT NOT NULL,
    flightSeatsCost INTEGER[] NOT NULL
);

CREATE TABLE IF NOT EXISTS Flights_Deleted (
    id UUID PRIMARY KEY,
    fromWhere VARCHAR(50) NOT NULL,
    destination VARCHAR(50) NOT NULL,
    flightTime TIMESTAMP NOT NULL,
    flightDuration INT NOT NULL,
    flightSeatsCost INTEGER[] NOT NULL,
    deleted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_by VARCHAR(50) NOT NULL
);

CREATE OR REPLACE FUNCTION move_flight_to_deleted() 
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO Flights_Deleted (id, fromWhere, destination, flightTime, flightDuration, flightSeatsCost, deleted_at, deleted_by)
    VALUES (OLD.id, OLD.fromWhere, OLD.destination, OLD.flightTime, OLD.flightDuration, OLD.flightSeatsCost, NOW(), current_user);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_flights_delete
AFTER DELETE ON Flights
FOR EACH ROW
EXECUTE FUNCTION move_flight_to_deleted();

-- тут другая дб

CREATE TABLE IF NOT EXISTS Tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flightId UUID NOT NULL,
    ownerId UUID NOT NULL,
    ticketCost INT NOT NULL,
    classOfService VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS Tickets_Deleted (
    id UUID PRIMARY KEY,
    flightId UUID NOT NULL,
    ownerId UUID NOT NULL,
    ticketCost INT NOT NULL,
    classOfService VARCHAR(50) NOT NULL,
    deleted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_by VARCHAR(50) NOT NULL
);

CREATE OR REPLACE FUNCTION move_ticket_to_deleted() 
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO Tickets_Deleted (id, flightId, ownerId, ticketCost, classOfService, deleted_at, deleted_by)
    VALUES (OLD.id, OLD.flightId, OLD.ownerId, OLD.ticketCost, OLD.classOfService, NOW(), current_user);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_tickets_delete
AFTER DELETE ON Tickets
FOR EACH ROW
EXECUTE FUNCTION move_ticket_to_deleted();