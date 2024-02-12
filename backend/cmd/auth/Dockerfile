FROM golang:1.21.5-bullseye
ENV AUTH_CONFIG_PATH=/backend/config/dev_auth.yaml
ENV AUTH_STORAGE_PATH="user=user password=password dbname=cloud_render_auth sslmode=disable"
ENV JWT_SECRET_KEY=randomkekv123
WORKDIR /backend
COPY . .
RUN go mod download && go mod verify
RUN go mod tidy
RUN go build -C ./cmd/auth -o ./../serverAuth
WORKDIR /backend-api/cmd/
CMD ./serverAuth
