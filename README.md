Chaos experiment for on-prem.


# How to use

## Add fault
```
./chaos-os fault create network loss -i bond0 --timeout 10 --dest-ip 11.167.254.40 --dest-port 8000 --percent 100
./chaos-os fault destroy --id aae1e0a5-ec39-469e-8248-cfe454c0bc1b
```



# Server

```bash
./chaos-os server start --port 12345 --background


curl -G -X PUT "localhost:12345/fault" --data-urlencode "cmd=network loss -i bond0 --timeout 100 --dest-ip 11.167.254.40 --dest-port 8000 --percent 100"
curl -X GET "localhost:12345/fault/status?id=46339f2f-ae1d-40b7-b9f4-ed08d8a7e71b"
curl -X DELETE "localhost:12345/fault?id=46339f2f-ae1d-40b7-b9f4-ed08d8a7e71b"
```
