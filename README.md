
# Web Application with Gin & Docker-client
* A web application with golang.
* Fetchs running containers info through docker daemon
* Keeps the information at cache instead of fetching through docker socket for each
  request

# Usage

## GO
```
go build main.go
./main
```

## DOCKER
```
docker run --rm -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock uzumlukek/docker-client
```
OR
```
docker-compose up --build
```

# Options
| Environment Variable | Default |
|----------------------|---------|
| APP_PORT             | 8080    |  
| LOG_LEVEL            | INFO    |
| LOG_FILE          | logfile |
| CACHE_DEFAULT_EXPIRATION_TIME          | 10      |
| CACHE_CLEANUP_INTERVAL_TIME          | 10      |

# Usecases
* **List All Containers:** Listing whole containers with a small amount of metadata
```
curl --location --request GET 'http://localhost:8080/containers/'
```
* **Get Container By ID:** Get details for the specified container
```
curl --location --request GET 'http://localhost:8080/containers/:containerId'
```
* **Run container:** Run a container as a daemon (for instance nginx)

```
curl --location --request POST 'http://localhost:8080/containers' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "test-uzumlukek-test",
    "image": "uzumlukek/local-firestore:testing",
    "exposePort": "3000",
    "publishPort": "3000"
}'
```
* **Delete Container By ID:** Delete the specified container
```
curl --location --request DELETE 'http://localhost:8080/containers/:containerId'
```

# Resources

### Docker-client
* https://docs.docker.com/engine/api/sdk/examples/
* https://pkg.go.dev/github.com/docker/docker/client#APIClient
* https://medium.com/trendyol-tech/golang-docker-client-ile-container-i%CC%87%C5%9Flemleri-6417884f4dbb

### Gin - Web Service
* https://go.dev/doc/tutorial/web-service-gin

### Caching
* https://github.com/patrickmn/go-cache

### Automapper - DTO & Model
* https://github.com/stroiman/go-automapper

### Logging
* https://pkg.go.dev/go.uber.org/zap

### Tracing - Opentelemetry & Gin & Jaeger Exporter
* https://opentelemetry.uptrace.dev/instrumentations/go-gin.html#installation
* https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/github.com/gin-gonic/gin/otelgin
* https://pkg.go.dev/go.opentelemetry.io/otel/exporters
* https://pkg.go.dev/go.opentelemetry.io/otel/exporters/jaeger
* https://github.com/open-telemetry/opentelemetry-go/blob/main/example/jaeger/main.go
#### Manual Instrumentation
```
tr := otel.Tracer("gin-gonic")
_, span := tr.Start(c.Request.Context(), "controller")
span.SetAttributes(attribute.Key("limit").String(limit))
span.AddEvent("Container Controller")
defer span.End()
```
### Pre-commit
* https://pre-commit.com/
* https://github.com/dnephin/pre-commit-golang

### Testing
* https://github.com/golang/mock
* https://github.com/gin-gonic/gin#testing
```
mockgen -source=service_layer/containerServiceInterface.go -destination=mock/mock_container_service.go -package=mock
```

### Regex
* https://yourbasic.org/golang/regexp-cheat-sheet/