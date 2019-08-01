
FROM mariadb:10.3.16

ARG VERSION

LABEL Name="eFantasy-mysql"
LABEL Version=${VERSION}

COPY config/mysql.cnf /etc/mysql/conf.d

EXPOSE 3306
