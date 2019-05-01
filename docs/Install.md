# Install

## Quick Install
* Make the account
  `useradd --home-dir /srv/Raven --create-home --system --shell /bin/nologin raven`

* Copy the files into install directory
```
rsync -av src/templates /srv/Raven/
rsync -av src/huggin /srv/Raven/
rsync -av etc/raven.ini /srv/Raven/
rsync -av etc/raven.service /etc/systemd/system/
```
* Edit `/srv/Raven/raven.ini` 
  Add whatever you want to monitor 

* Start the daemon
```
systemctl daemon-reload
systemctl enable raven  
systemctl start raven  
systemctl status raven  
```

* Point your web browser at `http://localhost:8000/` assuming you are on the local machine

## Nginx 
I've included `etc/nginx-raven.conf` which is an nginx reverse proxy file if you want to use a real webserver to front it for access control.

## Troubleshooting
* Make sure you firewall has port 8000 open if you are not on the same machine


