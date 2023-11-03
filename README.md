# switch-watch

A small CLI utility to monitor ports on SNMPv2-enabled network devices


## Building
```shell
go build -trimpath -ldflags="-w -s" .
```

## Usage
```shell
./switch-watch <ip> <community>

# example usage
./switch-watch 192.128.88.1 public
```
