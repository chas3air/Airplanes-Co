CREATE DATABASE airplanesco;

\c airplanesco

CREATE TABLE IF NOT EXISTS Customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role INT NOT NULL,
    surname VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    FOREIGN KEY (role) REFERENCES Role (id)
);

CREATE TABLE IF NOT EXISTS Role (
    id SERIAL PRIMARY KEY,
    roleName VARCHAR(15) NOT NULL UNIQUE
)

CREATE TABLE IF NOT EXISTS City (
    id SERIAL PRIMARY KEY,
    title VARCHAR(15) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Plains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    capacity INT NOT NULL
);

CREATE TABLE IF NOT EXISTS Flights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fromWhere INT NOT NULL,
    destination INT NOT NULL,
    flightTime TIMESTAMP NOT NULL,
    plainId UUID NOT NULL,
    FOREIGN KEY (fromWhere) REFERENCES City(id) ON DELETE CASCADE,
    FOREIGN KEY (destination) REFERENCES City(id) ON DELETE CASCADE,
    FOREIGN KEY (plainId) REFERENCES Plains(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ClassNames (
    id SERIAL PRIMARY KEY,
    title VARCHAR(15) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS CostTable (
    flightId UUID NOT NULL,
    classId INT NOT NULL,
    isLeft INT NOT NULL,
    Cost INT NOT NULL,
    PRIMARY KEY (flightId, classId),
    FOREIGN KEY (flightId) REFERENCES Flights(id) ON DELETE CASCADE,
    FOREIGN KEY (classId) REFERENCES ClassNames(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flightId UUID NOT NULL,
    ownerId UUID NOT NULL,
    classId INT NOT NULL,
    FOREIGN KEY (flightId) REFERENCES Flights(id) ON DELETE CASCADE,
    FOREIGN KEY (ownerId) REFERENCES Customers(id) ON DELETE CASCADE,
    FOREIGN KEY (flightId, classId) REFERENCES CostTable(flightId, classId) ON DELETE CASCADE
