setup:
	go get -u ./...

up_db:
	docker-compose up

run:
	DB_HOST=192.168.99.100 \
	DB_PORT=5432 \
	DB_USER=snapfile DB_PASSWORD=snapfile DB_NAME=snapfile \
	PORT=3000 go run main.go
