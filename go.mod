module qlogger

go 1.20

require internal/logger v1.0.0

replace internal/logger => ./internal/logger

require github.com/go-sql-driver/mysql v1.7.1

require github.com/joho/godotenv v1.5.1 // indirect
