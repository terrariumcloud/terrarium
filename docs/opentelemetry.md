# OpenTelemetry
Terrarium is instrumented with OpenTelemetry to allow monitoring of the operations for diagnosing problems and collecting information on usage.
For more information on OpenTelemetry see [here](https://opentelemetry.io/).

The docker compose file includes, for development a [Jaeger](https://www.jaegertracing.io/) all-in-one container to allow you to monitor traces. 
To access the Jaeger UI goto http://localhost:16686

The jaeger container also exposes port 4317 (OTLP grpc) and 4318 (OTLP http) to the host to make it easy to push traces 
into Jaeger when you're developing/debugging outside the docker network.