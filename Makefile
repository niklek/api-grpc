gen:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    places/places.proto

cert:
	cd cert; ./generate.sh; cd ..

.PHONY: gen cert 
