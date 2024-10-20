# Airplanes & Co

**Airplanes & Co** is a service for purchasing cheap flight tickets, developed using modern technologies in the Go programming language.

## Architecture

The service consists of several microservices:

1. **CLI**: Command-line interface for interacting with the application.
2. **BLL (Business Logic Layer)**: API representing the business logic of the application, which interacts with the DAL (Data Access Layer).
3. **DAL (Data Access Layer)**:
   - **Users**: API for managing user data at a lower level, available on port **12000**.
   - **Flights**: API for managing flight data at a lower level, available on port **12001**.
   - **Tickets**: API for handling tickets at a lower level, available on port **12002**.
4. **Cart**: API representing the order cart for each user, available on port **12003**.
5. **Flights Catalog**: API providing access to flight data through the DAL layer, available on port **12004**.
6. **Auth**: API used for system access and authentication, available on port **12005**.
7. **Flight Management**: API used for managing flights, available on port **12006**.
8. **Customer MAnagement**: **12007**

### Databases
- **PostgreSQL (PSQL)**: Relational database for storing structured data.
- **MongoDB**: NoSQL database for storing unstructured data.
- **Redis**: Caching system to enhance performance.

All **microservices** are built and run in Docker containers using Docker Compose.

## Installation

### Prerequisites

- Docker
- Docker Compose
- Git

### Cloning the Repository

To get started, clone the repository:

```bash
git clone https://github.com/chas3air/Airplanes-Co.git