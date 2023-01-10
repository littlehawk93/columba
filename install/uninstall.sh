#!/bin/bash

if [[ $EUID > 0 ]]; then
    echo "Must run as root"
    exit
fi

read -r -p "Are you sure you want to uninstall Columba? [Y/N] " input

case $input in
    [yY][eE][sS]|[Yy])
        ;;
    [nN][oO]|[Nn])
        exit
        ;;
    *)
        exit
        ;;
esac

echo "Uninstalling..."

service columba stop

rm -rf /var/www/columba
rm -rf /var/log/columba
rm -rf /var/data/columba
rm -rf /etc/columba
rm -rf /etc/systemd/system/columba.service
rm -rf /usr/bin/columba

userdel columba

echo "Done!"