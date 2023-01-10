#!/bin/bash

if [[ $EUID > 0 ]]; then
    echo "Must run as root"
    exit
fi

useradd -r -s /bin/false columba

mkdir -p /etc/columba
mkdir -p /var/log/columba
mkdir -p /var/data/columba
mkdir -p /var/www/columba

chown -R root:root /etc/columba
chown -R columba:columba  /var/log/columba
chown -R columba:columba /var/www/columba
chown -R columba:columba /var/data/columba

chmod -R 754 /etc/columba
chmod -R 754 /var/www/columba
chmod -R 754 /var/data/columba
chmod -R 755 /var/log/columba

cat > /etc/columba/config.yaml << EOL
database:
  database: /var/data/columba/db.sqlite
web:
  bind:
  port: 80
webroot: /var/www/columba
minrefreshtime: 1800
bgupdatetime: 3600
EOL

cat > /etc/systemd/system/columba.service << EOL
[Unit]
Description=Columba Self-Hosted Package Tracking Manager
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=on-success
RestartSec=1
User=columba
StandardOutput=/var/log/columba/output.log
StandardError=/var/log/columba/error.log
ExecStart=/usr/bin/columba --config /etc/columba/config.yaml

[Install]
WantedBy=multi-user.target
EOL