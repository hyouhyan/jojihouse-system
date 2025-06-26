FROM alpine:3.22

# タイムゾーン設定
ENV TZ=Asia/Tokyo
RUN apk --no-cache add tzdata

# PostgreSQLクライアントとMongoDBクライアントツールをインストール
RUN apk add --no-cache \
    postgresql-client \
    mongodb-tools