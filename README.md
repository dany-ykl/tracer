## Start jaeger in docker
```
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 14286:14268/tcp \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest
```