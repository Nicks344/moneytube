#!/bin/sh

# Install service
echo "[Unit]
Description=Moneytube service

[Service]
User=root
ExecStart=$(pwd)/moneytube
WorkingDirectory=$(pwd)
Restart=on-failure
Environment=\"GIN_MODE=release\"
Environment=\"MONEYTUBE_ENV=prod\"

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/moneytube.service
chmod -x /etc/systemd/system/moneytube.service
systemctl daemon-reload

# Enable autostart and start service
systemctl enable moneytube
service moneytube start
sleep 1
service moneytube status