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
```
📁 cmd/
--📁 user-profiles/
----📄 main.go
📁 configs/
--📄 config.go
📄 dump.go
📁 internal/
--📁 app/
----📄 app.go
----📁 auth/
------📄 cookie.go
------📄 dto.go
------📄 errors.go
------📄 handler.go
------📄 routes.go
------📄 service.go
----📁 profile/
------📄 admin-handler.go
------📄 dto.go
------📄 errors.go
------📄 handler.go
------📄 routes.go
------📄 service.go
--📁 domain/
----📁 token/
------📄 model.go
------📄 repository.go
----📁 user/
------📄 model.go
------📄 repository.go
--📁 http/
----📁 req/
------📄 decode.go
------📄 handler.go
------📄 validate.go
----📁 resp/
------📄 resp.go
--📁 middleware/
----📄 auth.go
----📄 ban.go
----📄 context.go
----📄 cors.go
----📄 role.go
--📁 security/
----📄 jwt.go
----📄 refresh.go
----📄 repository.go
--📁 storage/
----📁 postgres/
------📁 user/
--------📄 repository.go
📁 migrations/
--📄 0001_create_users_and_tokens_tables.down.sql
--📄 0001_create_users_and_tokens_tables.up.sql
--📄 0002_add_admin_to_users.up.sql
--📄 0003_add_more_columns_to_tokens.up.sql
--📄 0004_add_ban_column_to_users.up.sql
📁 pkg/
--📁 db/
----📄 postgres.go
```

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
