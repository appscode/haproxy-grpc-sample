## Docker Build

```
$ docker build -t diptadas/haproxy-grpc-sample-server:1.3 -f Dockerfile.server .
$ docker build -t diptadas/haproxy-grpc-sample-client:1.3 -f Dockerfile.client .
```

## Run Server/Client

```
$ go run server.go
2019/02/04 17:12:25
2019/02/04 17:12:25 Listening on address  :3000
```

```
$ go run client.go -crt=server-certs/server.crt -address=127.0.0.1:3000
2019/02/04 17:07:26 127.0.0.1:3000
2019/02/04 17:07:26 Codename result: name:"Daring Meitner"
```

## Configure TLS

### Create Certs

```
$ mkdir server-certs; cd server-certs
$ onessl create ca-cert
$ onessl create server-cert --domains ssl.appscode.test --ips 192.168.99.100,127.0.0.1
$ cat {server.crt,ca.crt} > bundle.crt
```

```
$ kubectl create secret tls server-secret -n demo --cert=server-certs/bundle.crt   --key=server-certs/server.key
secret/server-secret created
```

### Run Server/Client

```
$ go run server.go --key=server-certs/server.key --crt=server-certs/bundle.crt
2019/02/05 01:01:52 server-certs/bundle.crt server-certs/server.key
2019/02/05 01:01:53 Listening on address  :3000
```

```
$ go run client.go --address=127.0.0.1:3000 --crt=server-certs/bundle.crt
2019/02/04 17:07:26 127.0.0.1:3000 server.crt
2019/02/04 17:07:26 Codename result: name:"Daring Meitner"
```