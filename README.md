# CarVia 🌐 API Gateway

Центральна точка входу для всіх клієнтських запитів.

## 🛠 Технології

- Go
- net/http
- Reverse Proxy

## ⚙️ Функціонал

- Маршрутизація запитів до мікросервісів:
  - Lots Service
  - SSO Service
  - Storage Service
- Обробка CORS
- Логування запитів
- Агрегація відповідей (за потреби)

### ⚠️ Примітка

Gateway не містить бізнес-логіки — тільки проксування.

## 📡 Маршрути

- `/api/lots/*` → Lots Service
- `/api/sso/*` → SSO Service
- `/api/storage/*` → Storage Service

## 🚀 Запуск

```bash
go run main.go
```
