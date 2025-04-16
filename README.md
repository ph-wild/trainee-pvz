# trainee-pvz

## Генерация DTO и endpoint'ов из OpenAPI
- Использован `oapi-codegen` для генерации DTO и интерфейсов серверных хендлеров
- Все структуры и интерфейсы лежат в `internal/openapi/generated.go`
- Swagger-файл: `swagger.yaml`

## В main.go добавлен graceful shutdown
ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
defer cancelFunc()
...
<-ctx.Done()
slog.Info("Got shutdown signal, exit program")

## Добавлена конфигурация в yaml
/internal/config/config.go для изменения портов

## Добавлен Makefile
