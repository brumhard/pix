build-frontend:
	npm --prefix frontend run build

build-backend:
	go build -o ogframe main.go

build-backend-pi:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o ogframe_pi main.go

build: build-frontend build-backend

build-pi: build-frontend build-backend-pi

deploy-pi: build-pi
	scp ogframe_pi pi@raspberrypi:~/ogframe
	rm ogframe_pi

run: build-frontend
	go run main.go --images images/
