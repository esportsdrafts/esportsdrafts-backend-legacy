
FROM docker.pkg.github.com/barreyo/efantasy/efantasy-mysql:27c1535

ARG VERSION

LABEL Name="eFantasy-mysql"
LABEL Version=${VERSION}

COPY config/mysql.cnf /etc/mysql/conf.d

EXPOSE 3306
