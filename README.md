# CarVia 🚗 Lots Service

Сервіс для роботи з оголошеннями про продаж автомобілів.

## 🛠 Технології

- Go
- PostgreSQL
- REST API

## ⚙️ Функціонал

- CRUD операції з лотами
- Фільтрація та пошук
- Пагінація
- Лайки (обрані лоти)
- Перевірка авторизації

## 📡 Основні ендпоінти

- `/api/lots/filtered` - отримання лотів за параметрами
- `/api/lots/id/{lot_id}` - отримання лота по ID
- `/api/lots/brands` - отримання брендів
- `/api/lots/models` - отримання моделей за брендом
- `/api/lots/create_lot` - створення лота
- `/api/lots/update_lot/{lot_id}` - оновлення лота
- `/api/lots/delete_lot/{lot_id}` - видалення лота
- `/api/lots/likes/{lot_id}` - лайк / дизлайк
- `/api/lots/buy_lot/{lot_id}` Купівля лота

## 🗄 База даних

- Таблиці:
  - `sell_lots`
  - `cars`
  - `brands`
  - `models`
  - `liked_lots`

## 🚀 Запуск

```bash
go run main.go
```
