# webcam-monitor

The idea is to create a daemon that is able to
check the Mac logs and read if the webcam turned
on or off. Then I want to send a message to someone
(e.g. my wife) to alert her so she does not walk into the room.

## Run monitor

```bash
go build -o webcam_monitor
./webcam_monitor
```

## Run unit tests

```bash
go test -v
```

Tested on go version go1.23.3 darwin/arm64.
