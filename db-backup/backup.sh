#!/bin/sh

# shなのでbashの独自機能は使わないように注意

set -e # スクリプトのいずれかのコマンドが失敗したら、スクリプト全体を終了させる

# バックアップ保存先ディレクトリ (コンテナ内のパス)
BACKUP_DIR="/backups"

# 日付フォーマット
YEAR=$(date "+%Y")
MONTH=$(date "+%m")
FULLDATE=$(date "+%Y%m%d")
TIME=$(date "+%H%M%S")

# 日付によってバックアップディレクトリを作成
BACKUP_DIR="${BACKUP_DIR}/${YEAR}-${MONTH}/${FULLDATE}"
mkdir -p "${BACKUP_DIR}"

FILENAME=${FULLDATE}_${TIME}

echo "--- Starting backup at ${FILENAME} ---"

# --- PostgreSQLのバックアップ ---
echo "Dumping PostgreSQL database..."
# postgresが生きてるか確認
if ! pg_isready -U ${POSTGRES_USER} -h ${POSTGRES_HOST}; then
    echo "PostgreSQL is not ready. Exiting."
    exit 1
fi
# pg_dumpallを使って全データベースのバックアップを取得
pg_dumpall -U ${POSTGRES_USER} -h ${POSTGRES_HOST}  > "${BACKUP_DIR}/postgres_backup_${FILENAME}.sql"
echo "PostgreSQL dump successful."

# --- MongoDBのバックアップ ---
echo "Dumping MongoDB database..."
# mongodumpを使ってMongoDBのバックアップを取得
mongodump --host ${MONGO_HOST} --username ${MONGO_INITDB_ROOT_USERNAME} --password ${MONGO_INITDB_ROOT_PASSWORD} --authenticationDatabase admin --archive > "${BACKUP_DIR}/mongodb_backup_${FILENAME}.archive"
echo "MongoDB dump successful."

echo "--- Backup completed successfully ---"