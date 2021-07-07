build-frontend:
	npm --prefix frontend run build

build-backend:
	go build -o ogframe main.go

build: build-frontend build-backend

run: build-frontend
	go run main.go