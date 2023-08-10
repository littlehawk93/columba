#!/bin/bash

if [[ $EUID > 0 ]]; then
    echo "Must run as root"
    exit
fi

apt-get install jq unzip chromium -y

if [[ $version == "" ]];
then
    version=1000000
fi

curl -s https://api.github.com/repos/littlehawk93/columba/releases > /tmp/releases.json

versions=()

while read line; 
do
	version=$(echo $line | grep -o "[0-9]*\.[0-9]*\.[0-9]*")
	parts=($(echo $version | tr "." "\n"))
	versionnum=$(( ${parts[0]} * 1000000 + ${parts[1]} * 1000 + ${parts[2]} ))
	versions+=($versionnum)
done < <(cat /tmp/releases.json | jq .[].name)

if [[ ${#versions[@]} == 0 ]];
then
    echo "Unable to retreive release data"
    rm -f /tmp/releases.json
    exit 1
fi

versionIdx=0

for (( i=1; i<${#versions[@]}; i++ ))
do

	if [[ ${versions[$idx]} < ${versions[$i]} ]];
	then
		versionIdx=$i
	fi
done

currentVersion=0

if [ -f "/etc/columba/version" ];
then
    currentVersion=$(cat /etc/columba/version)
fi

if ! id -u columba &>/dev/null;
then
    useradd -r -s /bin/false columba
fi

if ! [ -d /etc/columba ];
then
    mkdir -p /etc/columba
fi

if ! [ -d /var/log/columba ];
then
    mkdir -p /var/log/columba
fi

if ! [ -d /var/data/columba ];
then
    mkdir -p /var/data/columba
fi

if ! [ -d /var/www/columba ];
then
    mkdir -p /var/www/columba
fi

chown -R root:root /etc/columba
chown -R columba:columba  /var/log/columba
chown -R columba:columba /var/data/columba
chown -R columba:columba /var/www/columba

chmod -R 755 /etc/columba
chmod -R 755 /var/data/columba
chmod -R 755 /var/log/columba
chmod -R 755 /var/www/columba

if ! [ -f /etc/columba/config.yaml ];
then
cat > /etc/columba/config.yaml << EOL
database:
    database: /var/data/columba/db.sqlite
web:
    bind:
    port: 80
browser:
    useragent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.102 Safari/537.36
    timeout: 30
webroot: /var/www/columba
minrefreshtime: 1800
bgupdatetime: 3600
EOL

chmod 644 /etc/columba/config.yaml
fi

if ! [ -f /etc/systemd/system/columba.service ];
then
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
fi

if [[ ${versions[$versionIdx]} == $currentVersion ]];
then
    echo "Columba is already up to date. Exiting"
    exit 0
fi

assets=()

while read line;
do
	filename=$(sed -e 's/^"//' -e 's/"$//' <<< $line)
	filename=$(tr "[:upper:]" "[:lower:]" <<< $filename)
	assets+=($filename)
done < <(cat /tmp/releases.json | jq ".[$versionIdx].assets[].name")

assetIdx=-1

for (( i=0; i<${#assets[@]}; i++ ))
do
	if [[ ${assets[$i]} == "rpi.zip" ]];
	then
		assetIdx=$i
		break
	fi
done

if [[ $idx == -1 ]];
then
	echo "Unable to find RPi update package"
	rm /tmp/releases.json
	exit 1
fi

url=$(cat /tmp/releases.json | jq ".[$versionIdx].assets[$assetIdx].url")
url=$(sed -e 's/^"//' -e 's/"$//' <<< $url)

if [[ $url == "" ]];
then
	echo "No download URL found for RPi update package"
	rm /tmp/releases.json
	exit 2
fi

curl -s -L \
 -H "Accept: application/octet-stream" \
 -H "Authorization: Bearer $token" \
 "$url" \
> /tmp/rpi.zip

rm -f /tmp/releases.json

unzip -qq -d /tmp/rpi /tmp/rpi.zip 
rm -f /tmp/rpi.zip

mv /tmp/rpi/columba /usr/bin

chown root:root /usr/bin/columba
chmod 755 /usr/bin/columba
setcap CAP_NET_BIND_SERVICE=+eip /usr/bin/columba

mv /tmp/rpi/install.sh /etc/columba/update.sh
mv /tmp/rpi/uninstall.sh /etc/columba

chown root:root /etc/columba/update.sh
chown root:root /etc/columba/uninstall.sh
chmod 740 /etc/columba/update.sh
chmod 740 /etc/columba/uninstall.sh

unzip -qq -d /var/www/columba /tmp/rpi/web.zip

chown -R columba:columba /var/www/columba
chmod -R 754 /var/www/columba

rm -rf /tmp/rpi

echo "${versions[$versionIdx]}" > /etc/columba/version
