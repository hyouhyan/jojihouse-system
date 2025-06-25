#!/bin/sh

# shなのでbashの独自機能は使わないように注意

set -e # スクリプトのいずれかのコマンドが失敗したら、スクリプト全体を終了させる

# バックアップ保存先ディレクトリ (コンテナ内のパス)
BACKUP_DIR="/backups"

# 日付フォーマット
DATE=$(date "+%Y%m%d_%H%M%S")

echo "--- Starting backup at ${DATE} ---"

# --- PostgreSQLのバックアップ ---
echo "Dumping PostgreSQL database..."
# 環境変数はOfeliaの定義から引き継ぐか、ここで指定します
pg_dumpall -U ${POSTGRES_USER} -h ${POSTGRES_HOST}  > "${BACKUP_DIR}/postgres_backup_${DATE}.sql"
echo "PostgreSQL dump successful."

# --- MongoDBのバックアップ ---
echo "Dumping MongoDB database..."
mongodump --host ${MONGO_HOST} --username ${MONGO_INITDB_ROOT_USERNAME} --password ${MONGO_INITDB_ROOT_PASSWORD} --authenticationDatabase admin --archive > "${BACKUP_DIR}/mongodb_backup_${DATE}.archive"
echo "MongoDB dump successful."

echo "--- Backup completed successfully ---"