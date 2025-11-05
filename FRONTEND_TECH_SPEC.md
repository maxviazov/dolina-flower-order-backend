# Техническое задание: Фронтенд для системы заказа цветов

## Обзор
Создать веб-приложение для заказа цветов на основе существующего бэкенда. Приложение должно позволять пользователям просматривать доступные цветы и создавать заказы.

## API Бэкенда
Бэкенд предоставляет следующие endpoints:

### Получение доступных цветов
- **URL:** `GET /api/v1/flowers`
- **Ответ:**
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
      "price": 0.0
    }
  ]
}
```

### Создание заказа
- **URL:** `POST /api/v1/orders`
- **Тело запроса:**
```json
{
  "mark_box": "VVA",
  "customer_id": "customer123",
  "items": [
    {
      "variety": "Red Naomi",
      "length": 70,
      "box_count": 10.5,
      "pack_rate": 20,
      "total_stems": 210,
      "farm_name": "KENYA FARM 1",
      "truck_name": "TRUCK A",
      "comments": "Optional comments",
      "price": 0.0
    }
  ],
  "notes": "Optional order notes"
}
```
- **Ответ:**
```json
{
  "id": 1,
  "mark_box": "VVA",
  "customer_id": "customer123",
  "status": "pending",
  "created_at": "2025-11-05T10:00:00Z",
  "items": ["..."],
  "notes": "Optional order notes"
}
```

### Получение списка заказов
- **URL:** `GET /api/v1/orders`
- **Query параметры:**
  - `limit` - количество заказов на странице (default: 50)
  - `offset` - смещение для пагинации (default: 0)
  - `status` - фильтр по статусу (optional)
- **Ответ:**
```json
{
  "orders": [
    {
      "id": 1,
      "mark_box": "VVA",
      "customer_id": "customer123",
      "status": "pending",
      "created_at": "2025-11-05T10:00:00Z",
      "total_items": 5
    }
  ],
  "total": 100,
  "limit": 50,
  "offset": 0
}
```

### Получение заказа по ID
- **URL:** `GET /api/v1/orders/{id}`
- **Ответ:** Полный объект заказа со всеми items

### Обновление заказа (для специалистов)
- **URL:** `PATCH /api/v1/orders/{id}`
- **Тело запроса:**
```json
{
  "status": "confirmed",
  "items": [
    {
      "variety": "Red Naomi",
      "length": 70,
      "price": 0.45
    }
  ]
}
```
- **Примечание:** Для MVP без аутентификации обновление статуса и цен доступно всем. В будущем будет добавлена проверка прав.
- **Ответ:** Обновленный объект заказа

### Коды ошибок
- `400 Bad Request` - неверные данные запроса
- `404 Not Found` - заказ/ресурс не найден
- `500 Internal Server Error` - ошибка сервера

Формат ошибки:
```json
{
  "error": "validation error: invalid customer_id"
}
```

## Функциональные требования

### 1. Страница списка цветов (`/flowers`)
- Отображение таблицы/карточек с доступными цветами
- Колонки: Variety, Length, Box Count, Pack Rate, Total Stems, Farm, Truck, Price
- Клиентская фильтрация по variety, farm, length (поиск в реальном времени)
- Клиентская сортировка по колонкам (по возрастанию/убыванию)
- Состояния:
  - **Loading**: показать скелетон/спиннер при загрузке
  - **Empty**: "Нет доступных цветов" если список пустой
  - **Error**: показать сообщение об ошибке с кнопкой "Повторить"
- Кнопка "Создать заказ" переводит на `/orders/new`

### 2. Страница создания заказа (`/orders/new`)
- Поля формы:
  - `mark_box` - выбор из dropdown (пока только "VVA", в будущем больше опций)
  - `customer_id` - текстовое поле (обязательное)
  - `notes` - textarea (необязательное)
- Секция добавления items:
  - Кнопка "Добавить позицию" открывает модал с выбором цветка
  - В модале: список цветов из `/api/v1/flowers` с поиском
  - После выбора цветка:
    - `variety`, `length`, `farm_name`, `truck_name` - автозаполняются
    - `box_count` - ввод числа (обязательное, > 0)
    - `pack_rate` - автозаполняется, но можно редактировать
    - `total_stems` - **автоматически рассчитывается** как `Math.floor(box_count * pack_rate)`
    - `price` - поле показывается как "0.00 USD" (будет установлено позже)
    - `comments` - текстовое поле (необязательное)
  - Каждая добавленная позиция отображается в списке с возможностью редактирования/удаления
- Валидация:
  - `customer_id` не пустой
  - Минимум 1 item в заказе
  - `box_count > 0` для каждого item
- Кнопки:
  - "Отмена" - возврат на `/flowers`
  - "Создать заказ" - отправка POST запроса, затем редирект на `/orders/{id}`
- Состояния:
  - Показывать спиннер на кнопке при отправке
  - Показывать toast с ошибкой при неудаче
  - Показывать toast "Заказ создан!" при успехе

### 3. Страница списка заказов (`/orders`)
- Таблица заказов с колонками: ID, Mark Box, Customer ID, Status, Дата создания, Кол-во позиций
- Пагинация (50 заказов на странице)
- Фильтр по статусу (dropdown: все/pending/confirmed/cancelled)
- Клик по строке ведет на `/orders/{id}`
- Цветовая индикация статуса:
  - `pending` - желтый
  - `confirmed` - зеленый
  - `cancelled` - красный
- Состояния: Loading, Empty, Error (аналогично странице цветов)

### 4. Страница детального просмотра заказа (`/orders/{id}`)
- Отображение информации:
  - ID, Mark Box, Customer ID, Status, Дата создания
  - Список всех items с деталями
  - Notes
  - **Total Amount**: сумма всех `item.price * item.total_stems` (если цены установлены)
- Для MVP: кнопка "Редактировать цены" открывает модал:
  - Список всех items
  - Для каждого item можно обновить `price`
  - Кнопка "Сохранить" отправляет PATCH запрос
- Кнопка "Обновить статус" с dropdown (pending/confirmed/cancelled)
- Кнопка "Назад к заказам" → `/orders`

### User Flow
```
/flowers (список цветов)
  ↓ [Создать заказ]
/orders/new (форма создания)
  ↓ [Создать заказ]
/orders/{id} (детали заказа)
  ↓ [Назад к заказам]
/orders (список заказов)
```

### Бизнес-логика цен
- При создании заказа клиентом все `price = 0.0`
- Специалист компании заходит на `/orders/{id}` и устанавливает цены после подтверждения от ферм
- После установки цен автоматически считается Total Amount
- В будущем: статус может меняться только специалистом (через модуль аутентификации)

## Технические требования

### Стек технологий
- **Framework**: React 18+ (с TypeScript)
- **HTTP клиент**: Axios для API запросов
- **CSS Framework**: Tailwind CSS
- **State management**: Context API + useReducer (для такого простого приложения достаточно)
- **Routing**: React Router v6
- **UI библиотека**: Shadcn UI или Headless UI (для модалов, dropdown)
- **Формы**: React Hook Form
- **Валидация**: Zod
- **Toast notifications**: React Hot Toast или Sonner

### Структура проекта
```
frontend/
  src/
    components/
      common/
        Button.tsx
        Input.tsx
        Modal.tsx
        Spinner.tsx
        Toast.tsx
      flowers/
        FlowerTable.tsx
        FlowerFilters.tsx
      orders/
        OrderForm.tsx
        OrderItemsList.tsx
        OrderStatusBadge.tsx
        OrderDetailsView.tsx
        PriceEditModal.tsx
    pages/
      FlowersPage.tsx
      OrdersListPage.tsx
      CreateOrderPage.tsx
      OrderDetailsPage.tsx
    services/
      api.ts          # Axios instance + API endpoints
      apiTypes.ts     # API request/response types
    types/
      flower.ts
      order.ts
    hooks/
      useFlowers.ts   # Fetch flowers logic
      useOrders.ts    # Orders CRUD operations
      useToast.ts     # Toast notifications
    utils/
      validation.ts   # Zod schemas
      calculations.ts # total_stems, total_amount calculations
    context/
      AppContext.tsx  # Global app state (if needed)
    App.tsx
    main.tsx
  public/
  package.json
  tsconfig.json
  tailwind.config.js
  vite.config.ts
```

### Дизайн
- Адаптивный дизайн для мобильных (375px+) и десктопных (1024px+) устройств
- Чистый, минималистичный интерфейс
- Использование цветов бренда:
  - Primary: зеленый (#10B981 / green-500)
  - Secondary: нейтральный серый
  - Status colors: желтый (pending), зеленый (confirmed), красный (cancelled)
- Шрифт: Inter или system font stack
- Spacing: использовать Tailwind spacing scale (4px increments)

### Безопасность
- Валидация всех входных данных на клиенте (Zod schemas)
- Обработка всех ошибок API с пользовательскими сообщениями
- Защита от XSS (React автоматически экранирует)
- Sanitize пользовательского ввода для comments/notes
- CORS настроен на бэкенде

### Производительность
- Code splitting: lazy loading страниц через React.lazy()
- Debounce для фильтров поиска (300ms)
- Мемоизация дорогих вычислений (useMemo для total_amount)
- Виртуализация длинных списков (если > 100 items) - react-window
- Оптимистичные UI обновления при создании/редактировании заказов

### Обработка ошибок
- **Network errors**: "Не удалось подключиться к серверу. Проверьте соединение."
- **400 Bad Request**: показать конкретное сообщение из response.error
- **404 Not Found**: "Заказ не найден"
- **500 Server Error**: "Ошибка сервера. Попробуйте позже."
- Все ошибки логируются в console для отладки
- Показывать retry кнопку при ошибках загрузки данных

## Этапы разработки

### Phase 1: Настройка (1-2 дня)
1. Создать Vite + React + TypeScript проект
2. Настроить Tailwind CSS
3. Настроить React Router
4. Создать базовую структуру папок
5. Настроить Axios instance с baseURL к бэкенду
6. Создать базовые типы (Flower, Order, OrderItem)

### Phase 2: Список цветов (2-3 дня)
1. Создать FlowersPage с таблицей
2. Реализовать useFlowers hook для fetch данных
3. Добавить фильтрацию и сортировку
4. Добавить состояния Loading/Empty/Error
5. Протестировать с реальным API

### Phase 3: Создание заказа (3-4 дня)
1. Создать CreateOrderPage с формой
2. Реализовать добавление items (модал выбора цветка)
3. Автоматический расчет total_stems
4. Валидация формы (Zod + React Hook Form)
5. Отправка POST запроса и редирект
6. Обработка ошибок и success уведомления

### Phase 4: Список и детали заказов (2-3 дня)
1. Создать OrdersListPage с таблицей и пагинацией
2. Реализовать фильтр по статусу
3. Создать OrderDetailsPage
4. Реализовать отображение Total Amount
5. Добавить возможность обновления статуса

### Phase 5: Редактирование цен (1-2 дня)
1. Создать PriceEditModal
2. Реализовать PATCH запрос для обновления цен
3. Обновление UI после сохранения
4. Валидация цен (должны быть >= 0)

### Phase 6: Полировка и тестирование (2-3 дня)
1. Проверить все user flows
2. Протестировать на мобильных устройствах
3. Добавить loading states везде где нужно
4. Оптимизировать производительность
5. Проверить обработку всех ошибок
6. Code review и рефакторинг

### Phase 7: Деплой (1 день)
1. Настроить production build
2. Настроить environment variables (API_BASE_URL)
3. Деплой на Vercel/Netlify
4. Проверить работу с production API

**Общее время: 12-18 дней**

## Критерии приемки

### Функциональность
- ✅ Все 4 страницы работают корректно
- ✅ Можно создать заказ с несколькими items
- ✅ Total stems рассчитывается автоматически
- ✅ Можно обновить цены и статус заказа
- ✅ Total amount отображается корректно
- ✅ Пагинация работает на странице заказов

### Валидация
- ✅ Нельзя создать заказ без customer_id
- ✅ Нельзя создать заказ без items
- ✅ Box count должен быть > 0
- ✅ Цены должны быть >= 0

### UX
- ✅ Все состояния Loading/Empty/Error реализованы
- ✅ Toast уведомления при успехе/ошибке
- ✅ Адаптивный дизайн работает на mobile и desktop
- ✅ Нет UI блокировок при асинхронных операциях

### Качество кода
- ✅ TypeScript без any types
- ✅ Компоненты разбиты логично
- ✅ Переиспользуемые компоненты вынесены в common/
- ✅ API вызовы централизованы в services/
- ✅ Валидация через Zod schemas
- ✅ Читаемый и поддерживаемый код

## Примечания для MVP
- Аутентификация будет добавлена в следующей фазе
- Пока любой может редактировать цены и статусы
- Валюта захардкожена как USD
- mark_box пока только "VVA"
- Поддержка только английского языка (можно добавить русский позже)
- Нет unit тестов в MVP (добавить в следующей итерации)</content>
<parameter name="filePath">/Users/maximviazov/Developer/Golang/GoLandWorkspace/dolina-flower-order-backend/FRONTEND_TECH_SPEC.md
