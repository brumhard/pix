# Pix
>
> This is a personal toy project designed for my specific usecase.

A server that shows off your pics collection to build an electronic photo frame with a Raspberry Pi.

## Run

```shell
git clone github.com/brumhard/pix.git
cd pix
docker compose up --build
```

## Pi Setup

### Browser UI

<https://diyprojects.io/open-html-page-starting-raspberry-pi-os-chromium-browser-full-screen-kiosk-mode/#.YRwaW1MzbLA>
<https://2021.jackbarber.co.uk/blog/2017-02-16-hide-raspberry-pi-mouse-cursor-in-raspbian-kiosk>

### Disable screensaver

<https://www.raspberrypi.org/documentation/computers/configuration.html#configuring-screen-blanking>

```shell
sudo raspi-config
# > update
# > display options > screen blanking > disable
```

### Enable VNC server

```shell
sudo raspi-config
# interfacing > vnc > enable
```
