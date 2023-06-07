FROM mysql:latest

ENV MYSQL_ROOT_PASSWORD=123456789
ENV MYSQL_DATABASE=mydatabase


COPY init.sql /docker-entrypoint-initdb.d/

