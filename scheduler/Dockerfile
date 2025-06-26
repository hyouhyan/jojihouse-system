# Ofeliaの公式イメージをベースにする
FROM mcuadros/ofelia:latest

# Ofeliaコンテナ内でdockerコマンドを使えるようにインストールする
# OfeliaはAlpine Linuxベースなのでapkコマンドを使用
RUN apk add --no-cache docker-cli

CMD ["daemon", "--docker"]