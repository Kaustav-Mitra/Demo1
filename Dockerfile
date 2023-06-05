FROM mydatabase:Dockerfile

ENV MYSQL_ROOT_PASSWORD=root 

COPY ./database_mydatabase.sql /docker-entrypoint-initdb.d/



