# HZPaste (REST API)

`HZPaste` is a simple REST API app, that serve source code snippets.
`HZPaste` uses [gin](https://github.com/gin-gonic/gin) web framework.
By default HZPaste running on port `:8888`

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

### Swagger
`HZPaste` provides swagger documentation and UI:
```
http://localhost:8888/swagger/index.html
```

### PGC
HZPaste has PGC (Paste Garbage Collector), that sweeps and removes obsolete
pastes from it's storage.

### Books
I'm in love with books. If you want to thank me, just help me to buy books from the list

[![buy-me-a-book](https://img.shields.io/badge/Amazon-Buy%20me%20a%20book-important)](https://www.amazon.com/hz/wishlist/ls/3NSSXQK5CTS8N?ref_=wl_share)
