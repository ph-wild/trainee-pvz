# trainee-pvz

## Генерация DTO и endpoint'ов из OpenAPI
- Использован `oapi-codegen` для генерации DTO. Сгенерированы только DTO, как было указано в задании, без интерфейсов и тд
- Все структуры лежат в `internal/openapi/types.gen.go`
- Swagger-файл: `swagger.yaml`

## В main.go добавлен graceful shutdown
ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
defer cancelFunc()
...
<-ctx.Done()
slog.Info("Got shutdown signal, exit program")

## Добавлена конфигурация в yaml
/internal/config/config.go для изменения портов

## Миграции через goose
migrate-up:
	goose -dir ./migrations postgres "$(DB_DSN)" up

migrate-down:
	goose -dir ./migrations postgres "$(DB_DSN)" down

## JWT авторизация

## Swagger
http://localhost:8080/swagger/index.html

## логирование

## graceful shutdown

## make run запуск