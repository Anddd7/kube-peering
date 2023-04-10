# Debug in container

## get the container ip

```zsh
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' kpeering_container

# 172.18.0.2
```

## configured in toml file

...

## start server/client in different container

```sh
# [kpeering] start server
make -C cmd/kpeering/ run args=start 

# [kpctl] start client
make -C cmd/kpctl/ run args=connect

# [kpctl] start fake app
make example_http_app

# [kpeering] curl or request from client
make example_http_client

```
