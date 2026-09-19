# API errors

Application API failures use `application/json` with three required string fields:

```json
{"code":"internal_error","message":"Something went wrong. Please try again.","request_id":"server-generated-uuid"}
```

`api/openapi.yaml` defines the shared schema. The HTTP status is authoritative;
`code` is a stable machine-readable identifier, while `message` is safe to display.
Clients must tolerate unknown codes. Current mappings are 400 `invalid_request`,
404 `not_found`, 405 `method_not_allowed`, and 500 `internal_error`.
405 responses preserve `Allow`.

Each request receives a new server-generated UUID in `X-Request-ID`, the request
context and the access log. Incoming request IDs are not trusted. Error responses
include the same ID in their body. Internal errors log their original cause with
that ID but never return it to clients. Health probes still skip access logging.
The error contract covers application API responses, not infrastructure/proxy
failures or responses whose headers/body have already been sent.

Use `internal/apierror.Write` for middleware and unexpected handler errors.
Generated server error hooks are configured in `cmd/server/main.go`; do not edit
generated code. Future expected domain errors should use the shared contract and
be documented in OpenAPI alongside their endpoint.

The frontend client rejects failed HTTP responses with `ApiError`. Valid API
errors display their public message and optional expandable request details.
Network failures have a connection message; malformed or non-JSON responses have
a generic message. Cancelled requests do not produce UI errors. No requests are
automatically retried. Field validation errors and localization are deferred
until their first concrete use case.
