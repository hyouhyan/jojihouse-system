# ベースイメージとして公式のGolangイメージを使用
FROM golang:1.24-alpine

# タイムゾーンを設定
ENV TZ=Asia/Tokyo

# 作業ディレクトリを設定
WORKDIR /go/src

# Go Modulesの使用を許可
ENV GO111MODULE=on
# ローカルのモジュールキャッシュを最適化
COPY go.mod .
COPY go.sum .
RUN go mod download

# ソースコードをコンテナにコピー
COPY . .

# ポート8080を公開
EXPOSE 8080

# ビルド
RUN go build -o main cmd/main.go

# ログディレクトリを作成
RUN mkdir -p /var/log/jojihouse

# 実行可能ファイルをデフォルトのコマンドとして設定
CMD ["/bin/sh", "-c", "exec ./main >> /var/log/jojihouse/entrance-system_$(date +'%Y-%m-%d_%H%M%S').log 2>&1"]
