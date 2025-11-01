# Quick Start Guide

## Запуск сервера

```bash
# Запуск из исходников
go run cmd/server/main.go

# Или собрать и запустить бинарник
go build -o bin/server cmd/server/main.go
./bin/server
```

Сервер запустится на `http://localhost:8080`

## Тестирование API

### Вариант 1: Использовать готовый скрипт
```bash
./test_api.sh
```

### Вариант 2: Вручную через curl

**Проверка здоровья:**
```bash
curl http://localhost:8080/health
```

**Получить доступные цветы:**
```bash
curl http://localhost:8080/api/v1/flowers
```

**Создать заказ:**
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "mark_box": "VVA",
    "customer_id": "test-customer-123",
    "items": [
      {
        "variety": "Red Naomi",
        "length": 70,
        "box_count": 10.5,
        "pack_rate": 20,
        "total_stems": 210,
        "farm_name": "KENYA FARM 1",
        "truck_name": "TRUCK A",
        "price": 4.07
      }
    ],
    "notes": "Test order"
  }'
```

**Получить заказ по ID:**
```bash
curl http://localhost:8080/api/v1/orders/{order_id}
```

## База данных

Проект использует SQLite. База данных `flowers.db` создается автоматически при первом запуске с тестовыми данными из таблицы VVA KENYA.

## Структура проекта

- `cmd/server/` - точка входа приложения
- `internal/app/` - инициализация приложения и роутинг
- `internal/domain/` - модели данных и интерфейсы
- `internal/handlers/` - HTTP handlers
- `internal/services/` - бизнес-логика
- `internal/repository/sqlite/` - работа с базой данных
- `internal/config/` - конфигурация
- `internal/logger/` - логирование

## Следующие шаги

1. Создать фронтенд для формирования заказов
2. Добавить функционал для специалистов компании
3. Реализовать генерацию Excel файлов для отправки на ферму
