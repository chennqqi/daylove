# DayLove

A tiny blog software, written with golang.

No title required. just plain text , may insert html, link, image or video what ever

## Install

```
go get -u -v github.com/netroby/daylove
go build
./daylove
```
Or you can install docker, then run docker container 

```
# to install docker on your platform
# wget -qO- get.docker.com | sudo sh
git clone https://github.com/netroby/daylove.git
./up.sh
```
Make sure your docker verision 1.9.1+
```
$ docker version
Client:
 Version:      1.9.1
 API version:  1.21
 Go version:   go1.4.2
 Git commit:   a34a1d5
 Built:        Fri Nov 20 13:20:08 UTC 2015
 OS/Arch:      linux/amd64

Server:
 Version:      1.9.1
 API version:  1.21
 Go version:   go1.4.2
 Git commit:   a34a1d5
 Built:        Fri Nov 20 13:20:08 UTC 2015
 OS/Arch:      linux/amd64

```

Once you docker up and running, you may access demo via http://127.0.0.1:8080
To login, you need visit http://127.0.0.1:8080/admin/login  (The password will be found in config.toml file)
To create blog , you can visit http://127.0.0.1:8080/admin/addblog



## Graceful restart 

And if you want reload daylove, just  run following command

```
docker kill -s HUP daylove
```
Or may you want to rebuild binary and graceful reload ?

```
git pull --rebase
./graceful-restart.sh
```

## License

MIT License

## Donate me please

### Bitcoin donate

```
136MYemy5QmmBPLBLr1GHZfkES7CsoG4Qh
```
### Alipay donate
![Scan QRCode donate me via Alipay](https://www.netroby.com/assets/images/alipayme.jpg)

**Scan QRCode donate me via Alipay**

