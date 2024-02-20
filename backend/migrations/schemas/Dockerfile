FROM postgres:16.0-bookworm
#ENV CONFIG_PATH=/render-service/backend-api/config/
#COPY ./api_db.sql /docker-entrypoint-initdb.d/
COPY ./db_init.sql /docker-entrypoint-initdb.d/
RUN chmod a+r /docker-entrypoint-initdb.d/*