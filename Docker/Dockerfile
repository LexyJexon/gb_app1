# Используем официальный образ Golang
FROM golang:1.22

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /opt/kcbk

# Копируем содержимое текущей директории (где находится Dockerfile) внутрь контейнера в папку /go/src/app
COPY ../project /opt/kcbk

# Установка beego и bee
RUN cd /opt/kcbk
RUN go get -u github.com/beego/bee/v2
RUN go install github.com/beego/bee/v2
RUN apt-get update && \
    apt-get install -y default-mysql-client && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Пробрасываем порт 8080
EXPOSE 8080

# Устанавливаем точку входа
ENTRYPOINT ["/opt/kcbk/entrypoint.sh"]# Используем официальный образ Golang
FROM golang:1.22

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /opt/kcbk

# Копируем содержимое текущей директории (где находится Dockerfile) внутрь контейнера в папку /go/src/app
COPY ../project /opt/kcbk

# Установка beego и bee
RUN cd /opt/kcbk
RUN go get -u github.com/beego/bee/v2
RUN go install github.com/beego/bee/v2
RUN apt-get update && \
    apt-get install -y default-mysql-client && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Пробрасываем порт 8080
EXPOSE 8080

# Устанавливаем точку входа
ENTRYPOINT ["/opt/kcbk/entrypoint.sh"]
