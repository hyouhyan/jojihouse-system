FROM alpine:3.22

# タイムゾーン設定
ENV TZ=Asia/Tokyo
RUN apk --no-cache add tzdata

# PostgreSQLクライアントとMongoDBクライアントツールをインストール
RUN apk add --no-cache \
    postgresql-client \
    mongodb-tools

# バックアップスクリプトをコピー
COPY backup.sh /usr/local/bin/backup.sh
# スクリプトに実行権限を付与
RUN chmod +x /usr/local/bin/backup.sh

CMD ["/usr/local/bin/backup.sh"]