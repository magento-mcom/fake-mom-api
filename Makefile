build:
	docker build -t fake-api .

run:
	docker run -it --rm -p 24213:24213 --name fake-api fake-api
