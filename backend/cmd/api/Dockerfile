FROM golang:1.22.0-bullseye
#ENV PGDATA=/postgres-storage
ENV API_CONFIG_PATH=/backend/config/dev_api.yaml
ENV API_STORAGE_PATH="user=user password=password dbname=cloud_render_api sslmode=disable"
ENV BUFFER_CONFIG_PATH=/backend/config/dev_buffer.yaml
ENV FILES_INPUT_PATH=/backend/files/input
ENV FILES_OUTPUT_PATH=/backend/files/output
ENV JWT_SECRET_KEY=randomkekv123
WORKDIR /backend
COPY . .
RUN go mod download && go mod verify
RUN go mod tidy
RUN go test ./... -v
RUN go build -C ./cmd/api -o ./../serverAPI
RUN go build -C ./cmd/buffer -o ./../serverBuff
WORKDIR /backend/cmd/
RUN chmod +x ./startAll.sh
CMD ./startAll.sh

