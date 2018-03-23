build:
	docker build -t fake-mom-api .

run:
	docker run -it --rm -p 24213:24213 --name fake-mom-api fake-mom-api
