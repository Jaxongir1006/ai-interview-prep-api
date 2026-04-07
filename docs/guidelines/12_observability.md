# Observability

Most observability concerns — logging, tracing, alerting — are handled automatically by the framework. Developers rarely need to configure these manually. The framework provides middleware for HTTP request logging, error alerting, and distributed tracing out of the box.

## Logging

- ALWAYS use `github.com/rise-and-shine/pkg/observability/logger`
- DO NOT use other logger packages

### What the Framework Handles

- **HTTP request/response logging** — handled by the logger middleware
- **Error alerting** — handled by the alerting middleware
- **Trace context propagation** — handled by the tracing middleware

### When to Write Logs Manually

- **HTTP clients** — add debug logs for outgoing requests and responses in `pkg/` HTTP clients
- **Use case operations** — log significant business events with contextual data (e.g., `clean-expired-sessions` logging the deleted count)
- **App lifecycle messages** — log when controllers (workers, consumers) start running and when they stop

### Logger API

```go
// Context-aware logging (always when ctx is available)
logger.WithContext(ctx).Info("message")
logger.WithContext(ctx).With("key", value).Error("error occurred")

// Named logger (for identifying the component: {module}_{layer} or component name or http client name, etc.)
logger.WithContext(ctx).Named("auth_usecase").With("operation_id", uc.OperationID()).Info("done")
```

### Log Levels

| Level   | When to Use                                                         |
| ------- | ------------------------------------------------------------------- |
| `debug` | Detailed information for debugging (HTTP client requests/responses) |
| `info`  | General operational information                                     |
| `warn`  | Warning conditions                                                  |
| `error` | Error conditions                                                    |

### Configuration

Use `json` encoding in production, `pretty` in development.

```yaml
logger:
  level: info # debug | info | warn | error
  encoding: json # json | pretty
```

## Distributed Tracing

- Standard: OpenTelemetry
- Trace context propagated through `ctx`
- All cross-service calls include trace headers
- Initialized automatically at application startup

## Metrics

- Standard: OpenTelemetry Metrics
- Expose application metrics for monitoring

## Error Alerting

- Library: `github.com/rise-and-shine/pkg/observability/alert`
- Provider-based alerting on application errors (e.g., Telegram)
- Initialized automatically at application startup

```yaml
alert:
  provider: telegram
  telegram_bot_token: <bot-token>
  telegram_chat_ids:
    - <chat-id>
```
