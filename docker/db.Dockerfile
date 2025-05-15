FROM mysql/mysql-server:8.0

COPY ./custom.cnf /etc/mysql/conf.d/custom.cnf
