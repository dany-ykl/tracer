## Installation:
```
go get github.com/dany-ykl/tracer@v1.0.2
```

## Example:
```go
package main

import (
	"context"
	"github.com/dany-ykl/tracer"
	"log"
)

func main() {
	cancelTracer, err := tracer.New(&tracer.Config{
		ServiceName:              "example.api",
		Host:                     "localhost",
		Port:                     "4318",
		Environment:              "dev",
		TraceRatioFraction:       1.0,
		OTELExporterOTLPEndpoint: "http://jaeger:4317",
	})
	if err != nil {
		log.Fatalln(err)
	}

	_, span := tracer.StartTrace(context.Background(), "main-example")
	span.End()
	
	cancelTracer(context.Background())
}

```

## Start jaeger in docker:
```
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 14286:14268/tcp \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest
```