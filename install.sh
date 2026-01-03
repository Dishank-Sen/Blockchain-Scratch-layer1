#!/usr/bin/env bash
set -e

SERVICE_NAME=bloc
BIN_PATH=/usr/local/bin/bloc
WORKDIR=/var/lib/bloc
USER=bloc

echo "[1/7] Installing bloc daemon"

# 1. Check binary
if ! command -v bloc >/dev/null; then
  echo "Error: bloc binary not found"
  exit 1
fi

# 2. Create system user if not exists
if ! id "$USER" >/dev/null 2>&1; then
  echo "Creating system user: $USER"
  useradd --system --no-create-home --shell /usr/sbin/nologin "$USER"
fi

# 3. Create working directory
mkdir -p "$WORKDIR"
chown -R "$USER:$USER" "$WORKDIR"

# 4. Initialize bloc repository CORRECTLY
if [ ! -d "$WORKDIR/.bloc" ]; then
  echo "Initializing bloc repository"
  cd "$WORKDIR"
  sudo -u "$USER" bloc init
fi


# 5. Write systemd service
cat >/etc/systemd/system/$SERVICE_NAME.service <<EOF
[Unit]
Description=Bloc Bootstrap Nodes
After=network.target

[Service]
User=$USER
Group=$USER
WorkingDirectory=$WORKDIR
ExecStart=$BIN_PATH start
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 6. Reload systemd
systemctl daemon-reload

# 7. Enable + start
systemctl enable $SERVICE_NAME
systemctl restart $SERVICE_NAME

echo "✅ Bloc daemon installed and running"
