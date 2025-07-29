### ✅ Минимальный функционал
Создание заказа (POST /orders)
Получение всех заказов (GET /orders)
Получение заказа по ID (GET /orders/:id)
Обновление статуса (PATCH /orders/:id/status)
Отправка событий в Kafka (order_created, order_status_updated)
producer для записи событий в brocker

**Order:**
- id
- userid (кто заказал)
- courierId (кто везет, default: 0) 
- products []Product
- status



