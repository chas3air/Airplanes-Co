# Airplanes & Co

Airplanes & Co — это сервис для покупки дешевых авиабилетов, разработанный с использованием современных технологий на языке Go.

## Архитектура

Сервис состоит микросервисов:

1. **CLI**: интерфейс командной строки для взаимодействия с приложением.
2. **BLL**: API, представляющее бизнес-логику приложения и взаимодействующее с DAL.
3. **DAL**:
   - **Пользователи**: API для работы с данными пользователей, доступный на порту **12000**. 
   - **Рейсы**: API для управления рейсами, доступный на порту **12001**.
   - **Билеты**: API для работы с билетами, доступный на порту **12002**. 
4. **Cart**: API, которое представляет корзину заказов каждого пользователя, запускается на порту **12003**.

. **PostgreSQL (PSQL)**: реляционная база данных для хранения данных.
. **MongoDB**: NoSQL база данных для хранения неструктурированных данных.
. **Redis**: система кэширования для повышения производительности.

Все микросервисы запускаются и собираются в Docker-контейнерах с помощью Docker Compose.

## Установка

### Предварительные требования

- Docker
- Docker Compose
- Git

### Клонирование репозитория

```bash
git clone https://github.com/chas3air/Airplanes-Co.git
```

```bash
cd Airplanes-Co
docker compose up
```