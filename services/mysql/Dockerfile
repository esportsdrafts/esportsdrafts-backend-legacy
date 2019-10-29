
FROM docker.pkg.github.com/barreyo/esportsdrafts/esportsdrafts-mysql:10-28-2019

ARG VERSION

LABEL Name="esportsdrafts-mysql"
LABEL Version=${VERSION}

COPY config/mysql.cnf /etc/mysql/conf.d

EXPOSE 3306
