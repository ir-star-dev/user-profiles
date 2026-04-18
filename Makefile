DB_URL=postgres://postgres:pass@localhost:5432/user_profiles?sslmode=disable

reset-db:
	migrate -path ./migrations -database "postgres://postgres:pass@localhost:5432/user_profiles?sslmode=disable" drop -f