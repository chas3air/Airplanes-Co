CREATE DATABASE airplanesco;

\c airplanesco

-- CREATE TABLE IF NOT EXISTS Customers (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     login VARCHAR(50) NOT NULL UNIQUE,
--     password VARCHAR(255) NOT NULL,
--     role INT NOT NULL,
--     surname VARCHAR(100) NOT NULL,
--     name VARCHAR(100) NOT NULL,
--     FOREIGN KEY (role) REFERENCES Role (id)
-- );

-- CREATE TABLE IF NOT EXISTS Role (
--     id SERIAL PRIMARY KEY,
--     roleName VARCHAR(15) NOT NULL UNIQUE
-- )

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
    fromWhere INT NOT NULL,
    destination INT NOT NULL,
    flightTime TIMESTAMP NOT NULL,
    airplaneId UUID NOT NULL,
    FOREIGN KEY (fromWhere) REFERENCES City(id) ON DELETE CASCADE,
    FOREIGN KEY (destination) REFERENCES City(id) ON DELETE CASCADE,
    FOREIGN KEY (airplaneId) REFERENCES Airplanes(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS CostTable (
    id SERIAL PRIMARY KEY NOT NULL,
    flightId UUID NOT NULL,
    className VARCHAR(15) NOT NULL,
    Cost INT NOT NULL,
    UNIQUE (flightId, className),
    FOREIGN KEY (flightId) REFERENCES Flights(id) ON DELETE CASCADE
);

-- CREATE TABLE IF NOT EXISTS Tickets (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     flightId UUID NOT NULL,
--     ownerId UUID NOT NULL,
--     classId INT NOT NULL,
--     FOREIGN KEY (flightId) REFERENCES Flights(id) ON DELETE CASCADE,
--     FOREIGN KEY (ownerId) REFERENCES Customers(id) ON DELETE CASCADE,
--     FOREIGN KEY (flightId, classId) REFERENCES CostTable(flightId, classId) ON DELETE CASCADE
-- )