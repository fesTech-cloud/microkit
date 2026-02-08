# Design Philosophy

## Core Principles

### 1. Minimal Convenience Layer
microkit is **not a framework**. It's a toolkit that provides thin abstractions over common infrastructure patterns while staying close to the Go standard library.

**What this means:**
- Wraps stdlib types, doesn't replace them
- Minimal API surface
- Easy to understand by reading the source
- No magic, no code generation

### 2. Composable by Default
Every component should work independently and compose naturally with existing Go ecosystem tools.

**What this means:**
- No global state
- Explicit dependencies
- Standard interfaces (context.Context, io.Reader, etc.)
- Easy to wrap or extend

### 3. Opinionated Defaults, Flexible Overrides
Provide sensible defaults for common use cases, but allow users to override when needed.

**What this means:**
- Retry with exponential backoff + jitter by default
- Graceful shutdown patterns built-in
- But users can bring their own implementations

## What microkit IS

✅ Clean abstractions for messaging (Kafka, RabbitMQ)  
✅ Network client wrappers with retry/timeout  
✅ Common patterns (retry, circuit breaker, graceful shutdown)  
✅ Minimal boilerplate for microservice infrastructure  

## What microkit is NOT

❌ A web framework (use stdlib, chi, echo, etc.)  
❌ An observability platform (use OpenTelemetry)  
❌ A service mesh (use Istio, Linkerd, etc.)  
❌ A complete microservice framework (use go-kit, go-micro, etc.)  

## Extension Points

While microkit stays minimal, it provides extension points for common needs:

### Middleware/Hooks (Planned)
```go
// Optional middleware interface for HTTP clients
type Middleware func(next RoundTripper) RoundTripper

client := microhttp.NewClient(10*time.Second,
    microhttp.WithMiddleware(tracingMiddleware),
    microhttp.WithMiddleware(metricsMiddleware),
)
```

### Custom Retry Logic
```go
// Users can implement custom retry strategies
type RetryStrategy interface {
    ShouldRetry(attempt int, err error, resp *Response) bool
    NextDelay(attempt int) time.Duration
}
```

## Design Decisions

### Why panic → error in Call()?
**Problem:** HTTP client implements network.Client interface which has a Call() method for gRPC. HTTP client doesn't support gRPC, so it panicked.

**Why panic was wrong:** Libraries should never panic for expected conditions. Panic means "programmer error, crash the process"—too harsh for a library.

**Solution:** Return `network.ErrUnsupportedOperation` instead.

### Why jitter in retry?
**Problem:** Without jitter, all clients retry at the same intervals, causing thundering herd.

**Solution:** Add randomization (50-100% of delay) to spread out retry attempts.

### Why status-code-based retry?
**Problem:** Some failures (429 rate limit, 503 service unavailable) should be retried even though the HTTP request succeeded.

**Solution:** Add `WithRetryOnStatus()` option to retry specific status codes.

### Why import alias recommendation?
**Problem:** Package named `http` collides with `net/http`, forcing users to alias one of them.

**Solution:** Recommend aliasing microkit's package as `microhttp` in all docs/examples.

## Future Considerations

### Observability Hooks (Under Discussion)
- Optional hooks for tracing (OpenTelemetry compatible)
- Optional hooks for metrics (Prometheus compatible)
- **Key:** Must be opt-in, not required

### Structured Logging (Under Discussion)
- Should microkit log internally?
- If yes, use slog (stdlib) and allow injection
- If no, return errors and let users log

### Testing Utilities (Planned)
- Mock implementations of interfaces
- Test helpers for common scenarios
- In-memory message brokers for testing

## Contributing

When proposing new features, ask:
1. Is this a common pattern across microservices?
2. Does it reduce boilerplate significantly?
3. Can it be implemented as a thin wrapper over stdlib/existing tools?
4. Does it stay composable with the ecosystem?

If yes to all four, it's probably a good fit for microkit.

## Questions?

Open an issue or discussion to talk about design decisions, scope, or future direction.
