GOARCH=amd64 go build -o counter
GOARCH=amd64 go build -o counter-go-grpc ./plugin-go-grpc
GOARCH=amd64 go build -o multer-go-grpc ./plugin-go-grpc-new
