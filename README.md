
## Autostart

## browser ui
https://diyprojects.io/open-html-page-starting-raspberry-pi-os-chromium-browser-full-screen-kiosk-mode/#.YRwaW1MzbLA

https://2021.jackbarber.co.uk/blog/2017-02-16-hide-raspberry-pi-mouse-cursor-in-raspbian-kiosk

### docker

> could also use docker technically

### systemd

service for server:
```ini
[Unit]
Description=ogframe server
# get the correct mount target with `systemctl list-units --type=mount` if needed and replace fritznas
After=network-online.target home-pi-fritznas.mount
Requires=network-online.target systemd-networkd-wait-online.service home-pi-fritznas.mount 

StartLimitIntervalSec=500
StartLimitBurst=5

[Service]
Restart=on-failure
RestartSec=5s

ExecStart=/home/pi/ogframe --images "/home/pi/fritznas/TOSHIBA_EXT/Media/Fotos n Vids/Kanada 2017"

[Install]
WantedBy=multi-user.target
```
save as `/etc/systemd/system/ogframe.service`
run with `sudo systemctl enable ogframe`
check status with `sudo systemctl status ogframe`

## disable screensaver

https://www.raspberrypi.org/documentation/computers/configuration.html#configuring-screen-blanking


```shell
sudo raspi-config
# > update
# > display options > screen blanking > disable
```


## enable vnc server

```shell
sudo raspi-config
# interfacing > vnc > enable
```