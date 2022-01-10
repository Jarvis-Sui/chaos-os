Chaos experiment for physical machine.


# How to use

## Add fault
```
./chaos-os fault create network loss -i bond0 --timeout 10 --dest-ip 11.167.254.40 --dest-port 8000 --percent 100
./chaos-os fault create network corrupt -i bond0 --timeout 10 --dest-ip 11.167.254.40 --dest-port 8000 --percent 100 --correlation 10
./chaos-os fault create network duplicate -i bond0 --timeout 10 --dest-ip 11.167.254.40 --dest-port 8000 --percent 100 --correlation 10
./chaos-os fault create network reorder -i bond0 --timeout 10 --dest-ip 11.167.254.40 --dest-port 8000 --percent 100 --delay 10
./chaos-os fault create network delay -i bond0 --timeout 10 --dest-ip 11.167.254.40 --dest-port 8000 --delay 100 --jitter 10 --correlation 10

./chaos-os fault create process pause --pattern test_proc --timeout 10
./chaos-os fault create memory stress --timeout 10 --worker-num 1 --bytes 10%
./chaos-os fault create cpu stress --cpu 2 --load 50 --taskset 3-4 --timeout 10

./chaos-os fault destroy --id aae1e0a5-ec39-469e-8248-cfe454c0bc1b
```


# Server

```bash
cp chaosos.service /usr/lib/systemd/system/
systemctl daemon-reload
systemctl enable chaosos.service
systemctl start chaosos.service

curl -G -X PUT "localhost:12345/fault" --data-urlencode "cmd=network loss -i bond0 --timeout 100 --dest-ip 11.167.254.40 --dest-port 8000 --percent 100"
curl -X GET "localhost:12345/fault/status?id=46339f2f-ae1d-40b7-b9f4-ed08d8a7e71b"
curl -X DELETE "localhost:12345/fault?id=46339f2f-ae1d-40b7-b9f4-ed08d8a7e71b"
```


# dependencies

1. [stress-ng](https://github.com/ColinIanKing/stress-ng). for centos 7 by default. for other OS, please build it
