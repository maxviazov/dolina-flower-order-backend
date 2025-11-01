# API Documentation

## Endpoints

### Health Check
```
GET /health
```
Проверка состояния сервера.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z",
  "service": "dolina-flower-order-backend"
}
```

### Ping
```
GET /api/v1/ping
```
Тестовый endpoint.

**Response:**
```json
{
  "message": "pong"
}
```

### Get Available Flowers
```
GET /api/v1/flowers
```
Получить список доступных цветов для заказа.

**Response:**
```json
{
  "flowers": [
    {
      "variety": "Red Naomi",
      "length": 70,
      "box_count": 10.5,
      "pack_rate": 20,
      "total_stems": 210,
      "farm_name": "KENYA FARM 1",
      "truck_name": "TRUCK A",
      "price": 0
    }
  ]
}
```

### Create Order
```
POST /api/v1/orders
```
Создать новый заказ.

**Request Body:**
```json
{
  "mark_box": "VVA",
  "customer_id": "customer-uuid",
  "items": [
    {
      "variety": "Red Naomi",
      "length": 70,
      "box_count": 10.5,
      "pack_rate": 20,
      "total_stems": 210,
      "farm_name": "KENYA FARM 1",
      "truck_name": "TRUCK A",
      "comments": "Optional comment",
      "price": 4.07
    }
  ],
  "notes": "Optional order notes"
}
```

**Response:**
```json
{
  "id": "order-uuid",
  "mark_box": "VVA",
  "customer_id": "customer-uuid",
  "items": [...],
  "status": "pending",
  "created_at": "2024-01-01T00:00:00Z",
  "notes": "Optional order notes",
  "total_amount": 856.7
}
```

### Get Order by ID
```
GET /api/v1/orders/:id
```
Получить заказ по ID.

**Response:**
```json
{
  "id": "order-uuid",
  "mark_box": "VVA",
  "customer_id": "customer-uuid",
  "items": [...],
  "status": "pending",
  "created_at": "2024-01-01T00:00:00Z",
  "total_amount": 856.7
}
```

## Order Status Values
- `pending` - Новый заказ
- `processing` - В обработке
- `farm_order` - Отправлен на ферму
- `completed` - Завершен
- `cancelled` - Отменен
