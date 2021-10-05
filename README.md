![HZPaste](https://github.com/ep4sh/hzpaste/actions/workflows/go.yml/badge.svg?branch=master)  

# HZPaste (REST API)
**HZPaste** is a simple REST API app, which serves source code snippets.  
**HZPaste** uses [gin](https://github.com/gin-gonic/gin) web framework.
By default HZPaste runs on port `:8888`

## Configuration

HZPaste server can be configured using the environment variables:
```
$ export HZPASTE_HOST=0.0.0.0
$ export HZPASTE_PORT=9090
```

## Run tests
```
make test
```

## Run debug verion
```
make run
```

## Run release verion
```
make release
```

## Generate swagger docs
```
make swag-init
```

## Endpoints and features

### UI  
There is an amazing UI for hzpaste backend:  
https://github.com/AlenaMaer/pastecodes

### Swagger
`HZPaste` provides swagger documentation:
```
http://localhost:8888/swagger/index.html
```

### PGC
HZPaste has PGC (Paste Garbage Collector), that sweeps and removes obsolete pastes from it's storage.

### Books
I'm in love with books. If you want to thank me, just help me to buy books from the list

[![buy-me-a-book](https://img.shields.io/badge/Amazon-Buy%20me%20a%20book-important)](https://www.amazon.com/hz/wishlist/ls/3NSSXQK5CTS8N?ref_=wl_share)
