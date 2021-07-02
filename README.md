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

## Endpoint and features

### Swagger
`HZPaste` provides swagger documentation and UI:
```
http://localhost:8888/swagger/index.html
```

### PGC
HZPaste has PGC (Paste Garbage Collector), that sweeps and removes obsolete
pastes from it's storage.
