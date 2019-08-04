build:
	protoc -I. --go_out=plugins=micro:.	\
		proto/product.proto
	sudo docker build -t product-api .

run:
	sudo docker run -p 50051:50051 \
		-e MICRO_SERVER_ADDRESS=:50051	\
		product-api