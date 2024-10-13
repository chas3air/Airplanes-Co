# Airplanes&Co

Airplanes & Co — это сервис для покупки дешевых авиабилетов, разработанный с использованием современных технологий на языке Go.

## Архитектура

Сервис состоит из нескольких микросервисов:

1. **CLI**: интерфейс командной строки для взаимодействия с приложением.
2. **BLL**: API, представляющее бизнес-логику приложения и взаимодействующее с DAL.
3. **DAL**: API, обеспечивающее работу с базами данных через паттерн репозиторий.
4. **PostgreSQL (PSQL)**: реляционная база данных для хранения данных.
5. **MongoDB**: NoSQL база данных для хранения неструктурированных данных.
6. **Redis**: система кэширования для повышения производительности.

Все микросервисы запускаются и собираются в Docker-контейнерах с помощью Docker Compose.

## Установка

### Предварительные требования

- Docker
- Docker Compose

### Клонирование репозитория

```bash
git clone https://github.com/chas3air/Airplanes-Co.git
cd Airplanes-Co
```

### Запуск приложения
```bash
docker compose up
```
