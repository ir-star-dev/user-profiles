# User profiles application

## Imported packages
 - github.com/go-chi/chi (v5)
 - github.com/joho/godotenv
 - github.com/go-playground/validator (v10)
 - github.com/golang-jwt/jwt (v5)
 - golang.org/x/crypto/bcrypt
 - github.com/google/uuid

 ### Database
  - postgres:16.4
  - github.com/lib/pq
  - github.com/golang-migrate/migrate (v4)
  - github.com/jmoiron/sqlx


=====================================================================
## Structure

### cmd/main.go - Entry point приложения
Роль:
- только запуск
- никаких бизнес-логик

### internal/app - Composition Root (DI контейнер)
Роль:
- собирает зависимости
- связывает слои

### internal/domain - Бизнес-сущности + интерфейсы (контракты)
Ключевая идея: домен не знает про postgres / jwt / http
user/
- модель User
- интерфейс Repository
token/
- модель RefreshToken
- интерфейс Repository

### internal/app/auth и profile - Application Layer (use cases)
Состоит из:
1. service.go
- бизнес-логика
- use cases
2. handler.go
- HTTP слой
- адаптер
3. dto.go
- входные/выходные структуры
4. routes.go
- регистрация роутов

### internal/infrastructure - Адаптеры (реализации интерфейсов)
db/postgres/user 
- реализация user.Repository
Смысл: 
- домен говорит: у меня есть репозиторий
- инфраструктура говорит: я реализую это через postgres

security
- JWT
- refresh токены
- token repository

### internal/middleware - Cross-cutting concerns
- JWT middleware
- CORS
- context
- logger

### pkg - Переиспользуемые утилиты
- req (decode + validate, body handler)
- resp (json response)
- db (connection)
- logger