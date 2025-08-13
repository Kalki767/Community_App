#Base image
FROM golang:1.23-alpine

#Set workdir inside the container
WORKDIR /app

#copy the dependencies for efficiency
COPY go.mod go.sum ./

#install the dependencies
RUN go mod download

#copy the full code
COPY . .

#build the binary
RUN go build -o auth_service ./cmd/main.go

#Expose the port the service runs on
EXPOSE 8080

#Command to run the service
CMD ["./auth_service"]