#!/usr/bin/env bash
set -euo pipefail

IMAGE_NAME="spider:latest"
CONTAINER_NAME="spider"

DATA_DIR="/home/laoin/spider-data"
LOGS_DIR="/home/laoin/spider-logs"
CONFIG_DIR="/home/laoin/spider-config"

# Ensure persistent directories exist
mkdir -p "$DATA_DIR" "$LOGS_DIR" "$CONFIG_DIR"

echo "[1/5] Building image: $IMAGE_NAME"
docker build -t "$IMAGE_NAME" .

echo "[2/5] Removing old container if exists: $CONTAINER_NAME"
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
  docker rm -f "$CONTAINER_NAME" >/dev/null
fi

echo "[3/5] Getting app uid/gid from image"
APP_UID="$(docker run --rm --entrypoint id "$IMAGE_NAME" -u app | sed -n 's/uid=\([0-9]\+\).*/\1/p')"
APP_GID="$(docker run --rm --entrypoint id "$IMAGE_NAME" -g app | tr -d '\r\n')"

if [[ -n "$APP_UID" && -n "$APP_GID" ]]; then
  echo "[4/5] Fixing permissions for mounted dirs: ${APP_UID}:${APP_GID}"
  chown -R "${APP_UID}:${APP_GID}" "$DATA_DIR" "$LOGS_DIR" "$CONFIG_DIR"
  chmod -R u+rwX "$DATA_DIR" "$LOGS_DIR" "$CONFIG_DIR"
else
  echo "[WARN] Could not detect app uid/gid, skipping chown/chmod"
fi

echo "[5/5] Starting container: $CONTAINER_NAME"
docker run -d --name "$CONTAINER_NAME" \
  --add-host=host.docker.internal:host-gateway \
  -v "$DATA_DIR:/spider/data" \
  -v "$LOGS_DIR:/spider/logs" \
  -v "$CONFIG_DIR:/config" \
  -e RUN_MODE=prod \
  "$IMAGE_NAME" ./spider run

echo "Done."
docker ps --filter "name=${CONTAINER_NAME}" --format 'table {{.Names}}\t{{.Status}}\t{{.Image}}'
echo "Use: docker logs -f ${CONTAINER_NAME}"