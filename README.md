# Columba
Columba is a self-hosted package tracking system that was designed and tested to run on a Raspberry Pi. The name comes from the suborder of birds that contains pigeons and doves, most notably, the carrier pigeon.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Configuration Settings](#configuration-settings)
- [Updating Columba](#updating-columba)
- [Uninstalling Columba](#uninstalling-columba)
- [Roadmap](#roadmap)

## Features

*Columba* has a simple, mobile-friendly web interface for quickly viewing all your active packages. There is a selectable dark or light mode setting, and 3 different UI layouts for different users' preferences. Adding a package requires only the tracking number and shipping service provider, but an optional label can be provided if desired. Packages that have been delivered are highlighted in green, and can be removed from the view to clean up the list of packages if needed.

## Installation

*Columba* was designed to be easy to install on a Raspberry Pi.The easiest way to install *Columba* is to use `wget` or `curl` to pull the `install.sh` script file from the GitHub repository and execute it.

**cURL Example**
```
curl -s https://raw.githubusercontent.com/littlehawk93/columba/main/install/install.sh | /bin/bash
```
**wget Example**
```
wget -q https://raw.githubusercontent.com/littlehawk93/columba/main/install/install.sh | /bin/bash
```

The install script will pull the latest version of *Columba* from the GitHub releases and create the necessary service account, directories, and files for Columba to run as a service on your Pi. Once the install script finishes, you can enable the service with the command:
```
service columba enable
```
and start the webapp using:
```
service columba start
```

## Configuration Settings

Configuration settings are saved in YAML format in the file: `/etc/columba/config.yaml`. A typical installation will produce a configuration file that looks like this:
```
database:
    database: /var/data/columba/db.sqlite
web:
    bind:
    port: 80
webroot: /var/www/columba
minrefreshtime: 1800
bgupdatetime: 3600
```

Below are the configuration settings and their purposes

| Parent   | Property Name  | Description                                                                                                                                                                                                                                                                                        | Default |
|----------|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| database | database       | The path to the SQLite database file used to store application data. If the file does not exists, *Columba* will create a new empty database file when it launches                                                                                                                                 |         |
| web      | bind           | The IP address to bind the *Columba* web process to. Leave blank to bind to all IP addresses                                                                                                                                                                                                       |         |
| web      | port           | The TCP port to bind the *Columba* web process to                                                                                                                                                                                                                                                  |    80   |
|          | webroot        | The directory containing the web application static files                                                                                                                                                                                                                                          |         |
|          | minrefreshtime | The minimum time, in seconds, before *Columba* will attempt to scrape a tracking service for shipment updates. This is a safety tool to prevent tracking services from blocking your IP address due to spamming                                                                                    |         |
|          | bgupdatetime   | Time in seconds between background updates. Each time a background update triggers, *Columba* will scrape all active packages for updates and save any new tracking events, even if no clients are connected to the website. This setting ignores the minrefreshtime parameter, so use responsibly |         |

## Updating Columba

During installation, *Columba* saves a copy of the install script as the file `/etc/columba/update.sh`. If you run the script again, it will check the current installed version of Columba and update executable and web files if a newer version is found on GitHub. It will **not** overwrite your configuration settings or SQLite database files. 

## Uninstalling Columba

To uninstall *Columba*, you can execute the script file located at `/etc/columba/uninstall.sh`. This will remove all the default directories, web files, and executable and remove the service account from your Pi. If you have changed the SQLite database file location or web files locations from their default locations, the uninstall script will not remove those files.

## Roadmap

Support for other tracking services is planned, but ironically there is no estimated time frame for when new features will arrive. Currently, the next providers planned are:
- FedEx 
- AliExpress Shipping

If you have ideas or requests for other shipment trackers. Please feel free to make an issue on the github project page. 