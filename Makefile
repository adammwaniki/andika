# Build command to simplify deployment process
build:
	go run github.com/a-h/templ/cmd/templ@latest generate
	go run worker_main.go
	mkdir -p dist
	cp -r public/* dist/
	cp -r static/* dist/