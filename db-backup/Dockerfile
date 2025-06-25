FROM ubuntu:latest

# タイムゾーン設定
ENV TZ=Asia/Tokyo

# mongodbインストール用
RUN apt update && \
    apt-get install gnupg curl -y
# GPGキーを追加
RUN curl -fsSL https://www.mongodb.org/static/pgp/server-8.0.asc | \
    gpg -o /usr/share/keyrings/mongodb-server-8.0.gpg \
    --dearmor
# MongoDBのリポジトリを追加
RUN echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-8.0.gpg ] https://repo.mongodb.org/apt/ubuntu noble/mongodb-org/8.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-8.0.list

# PostgreSQLインストール用
RUN apt-get update && \
    apt install -y postgresql-common && \
    /usr/share/postgresql-common/pgdg/apt.postgresql.org.sh -y

# PostgreSQLクライアントとMongoDBクライアントツールをインストール
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    postgresql-client \
    mongodb-org-tools \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*