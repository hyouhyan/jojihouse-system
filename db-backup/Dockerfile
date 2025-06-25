# Debian系の軽量イメージをベースにする
FROM debian:bullseye-slim

# タイムゾーン設定
ENV TZ=Asia/Tokyo

# PostgreSQLクライアントとMongoDBクライアントツールをインストール
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    postgresql-client \
    mongodb-database-tools \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*