build:
	docker build -t example-otelgo-miniservice .

run:
	docker run --env-file .env -d -p 8080:8080 --name example-otelgo example-otelgo-miniservice