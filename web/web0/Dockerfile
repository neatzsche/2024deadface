FROM node:18-buster

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y default-mysql-server && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./targetlist2 .

COPY ./init.sql /docker-entrypoint-initdb.d/init.sql

RUN npm install

EXPOSE 3001

ENV MYSQL_ROOT_PASSWORD=your_mysql_root_password
ENV MYSQL_DATABASE=my_app_db

CMD service mysql start && \
    sleep 5 && \
    echo "ALTER USER 'root'@'localhost' IDENTIFIED BY '${MYSQL_ROOT_PASSWORD}';" | mysql -u root && \
    echo "CREATE DATABASE IF NOT EXISTS ${MYSQL_DATABASE};" | mysql -u root -p${MYSQL_ROOT_PASSWORD} && \
    node init_db.js && \
    npm run build && \
    npm run start