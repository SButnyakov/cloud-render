FROM golang:1.21.5-bullseye
#ENV CONFIG_PATH=config/
ENV PGDATA=/postgres-storage
ENV API_CONFIG_PATH=/backend-api/config/localapi.yaml
ENV API_STORAGE_PATH="user=user password=password dbname=cloud_render_api sslmode=disable"
ENV BUFFER_CONFIG_PATH=/backend-api/config/localbuffer.yaml
ENV FILES_INPUT_PATH=/backend-api/files/input
ENV FILES_OUTPUT_PATH=/backend-api/files/output
ENV JWT_SECRET_KEY=randomkekv123
WORKDIR /backend-api
COPY . .
RUN go mod download && go mod verify
RUN go mod tidy
RUN go build -C ./cmd/api -o ./../serverAPI
RUN go build -C ./cmd/buffer -o ./../serverBuff
#RUN wget https://download.redis.io/redis-stable.tar.gz
#RUN tar -xzvf redis-stable.tar.gz
#WORKDIR /backend-api/redis-stable
#RUN make
#WORKDIR /backend-api/redis-stable/src
#RUN cp ./redis-server /backend-api/cmd/
WORKDIR /backend-api/cmd/
RUN chmod +x ./startAll.sh
#, "&&", "./serverBuff", "&&", "./redis-server"]
#CMD ["./serverAPI", "&&", "./serverAuth"]
CMD ./startAll.sh

