# Сервис ПВЗ
## Описание
Cервис для сотрудников ПВЗ, который позволяет:​
- Регистрировать новые ПВЗ в Москве, Санкт-Петербурге и Казани (доступно только модераторам).
- Инициировать приёмку товаров (доступно сотрудникам ПВЗ).
- Добавлять и удалять товары (LIFO) в рамках приёмки (доступно сотрудникам ПВЗ).
- Закрывать приёмку (сотрудник ПВЗ).
- Получать информацию о ПВЗ с фильтрацией по дате.

## Стек
Язык программирования: Go  
Протоколы запросов: HTTP и gRPC (для списка ПВЗ)  
Мониторинг и сбор метрик: Prometheus  
База данных: PostgreSQL  
Контейнеризация: Docker Compose (для postgres и prometheus)  

## Установка и запуск
1. Клонируйте репозиторий.
```
git clone github.com:ph-wild/trainee-pvz.git && cd trainee-pvz
```
2. Соберите и запустите сервис.  
Для настройки среды: `make up`  
Для создания БД необходим goose (`go install github.com/pressly/goose/v3/cmd/goose@latest`): `make migrate-up`  
Для запуска приложения: `make run`  
```
make run-all
```
Сервис будет доступен по следующим адресам:​  
HTTP API: `http://localhost:8080`  
gRPC API: `localhost:3000`  
Метрики Prometheus: `http://localhost:9000/metrics​`  
Swagger: `http://localhost:8080/swagger/index.html`  

3. Для запуска тестов с процентом покрытия выполните:
```
make test
```

## Структура проекта и его архитектура
```
|trainee-pvz/
├── api/                        # Swagger спецификации
├── cmd/                        # Точка входа (main.go)
├── config/                     # Конфигурация (парсинг config.yaml)
├── internal/
│   ├── auth/                   # JWT-авторизация
|   ├── database/               # Коннект к БД (sqlx)
│   ├── errors/                 # Общие ошибки (errors.New)
│   ├── grpc/                   # gRPC логика и proto-файлы
│   ├── handler/                # HTTP-обработчики (chi) и middlewares
│   ├── integration/            # Интеграционный тест
│   ├── metrics/                # Прометеус метрики
│   ├── models/                 # Внутренние структуры данных
│   ├── openapi/                # Сгенерированные структуры из OpenAPI (DTO)
│   ├── repository/             # Доступ к базе данных (sqlx)
│   └── service/                # Бизнес-логика (сервисы) и unit-tests
├── migrations/                 # SQL-миграции (goose)
├── prometheus/                 # Конфиг для прокидывания внуть контейнера в Prometheus
├── .gitignore                  # untracked files для Git
├── config.yaml                 # Файл конфигурации приложения (порты, auth и прочее)
├── docker-compose.yml          # Поднятие окружения (postgres, prometheus)
├── go.mod                      # Список версий зависимостей
├── go.sum                      # Контрольные суммы
├── Makefile                    # Файл с инструкциями
└── README.md                   # Документация
```

## Схема базы данных
```
+-------------------+       +-------------+       +--------------+
| pvz               |       | receptions  |       | products     |
+-------------------+       +-------------+       +--------------+
| id                |-------+ pvz_id      |       | id           |
| registration_date |       | datetime    |       | datetime     |
| city              |       | status      |       | type         |
+-------------------+       | id          |-------+ reception_id |
                            +-------------+       +--------------+
```
Cвязи:
- pvz.id -> receptions.pvz_id (1 ко многим)
- receptions.id -> products.reception_id (1 ко многим)


# Основные задания
## Swagger
Swagger API доступно по `http://localhost:8080/swagger/index.html`  
Все handlers с из назначением видны по ссылке выше.  
Доступные роли: moderator/employee.  
Используемый формат для фильтрации по времени в методе `GET /pvz` - `time.RFC3339`.  

## Авторизация
`/dummyLogin` возвращает заранее сгенерированный токен на основе выбранной роли пользователя (moderator/employee).  
В дополнительном задании реализована авторизация через `JWT`

## unit-tests 
Находятся в `internal/service`  
Запуск тестов: `make test`  
Тестовое покрытие составляет 97.0%  
```
ok      trainee-pvz/internal/service    3.616s  coverage: 97.0% of statements
```

## Интеграционный тест
- Создает новый ПВЗ (рандомный город из трех возможных)
- Добавляет новую приёмку заказов
- Добавляет 50 товаров (рандомных типов из возможных) в рамках текущей приёмки заказов
- Закрывает приёмку заказов  

Находится в папке `internal/intergation`, запускается вместе с unit-tests через `make test`


# Дополнительные задания
## 1. JWT-авторизация 
Методы: `/register` и `/login`.  
Возможные role: moderator или employee.
Пароль хэшируется через bcrypt.
В качкстве ответа возвращается JWT-токен.  

Реализовано в `internal/auth/jwt.go`

## 2. gRPC API
Сервис предоставляет следующие gRPC-методы:​
- GetPVZList — возвращает список всех зарегистрированных ПВЗ.​

Пример использования с grpcurl:​
```
grpcurl -plaintext -d '{}' localhost:3000 pvz.v1.PVZService/GetPVZList
```
Реализация находится в `internal/grpc`

## 3. Prometheus метрики
Сервис собирает и предоставляет следующие метрики:​  
Технические:
- Количество HTTP-запросов.
- Время ответа.  

Для реализации используется middleware в модуле handler.

Бизнесовые:
- Количество созданных ПВЗ.
- Количество созданных приёмок.
- Количество добавленных товаров​.  

Количество созданных сущностей пишется из слоя бизнес-логики в `internal/service`. Инкрементируем счетчик при каждом успешном сохранении соответствующего Entity  


| метрика                             | тип                        | лейблы                       | описание                                                                                                       |
|-------------------------------------|----------------------------|------------------------------|----------------------------------------------------------------------------------------------------------------|
| http_request_duration_summary       | Summary                    | code, path, method, quantile | Время обработки запроса в секундах. Распределение по квантилям. Разбиение по коду ответа, url и методу запроса |
| http_request_duration_summary_count | Counter (часть от Summary) | code, path, method           | Количество запросов. Разбиение по коду ответа, url и методу запроса                                            |
| created_entity_count                | Counter                    | entity                       | Количество созданных сущностей. С разбиением по типу сущности                                                  |

Метрики доступны по адресу: http://localhost:9000/metrics​

## 4. Логирование в проекте
- Осуществляется с помощью пакетов `log/slog`  
- Логируются HTTP запросы через middleware в `internal/handler/middleware.go`
- Добавлены трассировка через `errors.Wrap()` и созданы некоторые ошибки через `errors.New()` (перечень ошибок в `internal/errors`)

## 5. Генерация DTO и endpoint'ов из OpenAPI
- Использован `oapi-codegen` для генерации DTO. Сгенерированы только DTO (struct), как было указано в задании
- Все структуры лежат в `internal/openapi/types.gen.go`
- Swagger-файл: `api/swagger.yaml`  

# Дополнительный блок
## Graceful shutdown
Добавлен в main.go:
```
ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
defer cancelFunc()
...
<-ctx.Done()
slog.Info("Got shutdown signal, exit program")
```

## config.yaml
Необходимые для запуска параметры вынесены в файл `/config.yaml`  
В том числе, перечислены все необходимые по заданию порты.  
Считывание файла в структуры производится в `/config/config.go`

## Миграции 
Для миграций используется goose. Установка: `go install github.com/pressly/goose/v3/cmd/goose@latest`  
Запуск через `make migrate-up`  
Находятся в папке `/migrations`

## Dockerfile
Дополнительно добавлена ветка InfraDocker, PR https://github.com/ph-wild/trainee-pvz/pull/1/files в которую вошло развертывание самого приложения PVZ через docker compose (добавлен Dockerfile). Через make run-all приложение разворачивается в контейнере (миграции в базу проходят, swagger открывается и все отрабатывает, метрики от prometheus доступны, gRPC отрабатывает). Но не успеваю дотестировать работоспособность и привести в порядок README и Makefile, что может запутать при тестировании моего решения, поэтому доработка не вошла в main.