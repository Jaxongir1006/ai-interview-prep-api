# Controllers

## Controller Rules

- **One-to-one mapping** — each use case has exactly one controller entry point
- **No business logic** — delegate everything to use cases
- **Always use forwarders** — use `forward.To*` and `worker.ForwardTo*` helpers instead of writing manual handlers

## Forwarder Pattern

Always wire use cases through forwarder functions. Forwarders handle request parsing, validation triggering, response formatting, and error mapping automatically.

### When Manual Handlers Are Acceptable

Manual handlers are only acceptable when the generic forwarder cannot handle the input format — for example, use cases with user_action type that accept **file uploads** or generates files in response. In such cases, write a custom handler acceptable.

## API Design

We use HTTP but **do not follow REST conventions**. Our API is intentionally simplified to two verbs:

- **GET** — operations that don't change state (queries, lookups, listings)
- **POST** — operations that change state (creates, updates, deletes, actions)

### URL Structure

```
{method} api/v1/{module}/{usecase-operation-id}
```

Examples:

```
GET  api/v1/auth/list-users?role=admin
POST api/v1/auth/create-user
POST api/v1/auth/deactivate-user
GET  api/v1/catalog/get-product-details?id=123
```

### Rules

- **No path parameters** — pass identifiers via query params (GET) or request body (POST)
- **GET inputs** — query parameters only
- **POST inputs** — JSON request body only (exception: file uploads)
- **Operation ID = use case name** — the URL segment maps directly to a use case

## Wiring Reference

| Type      | Wiring                                                                                           |
| --------- | ------------------------------------------------------------------------------------------------ |
| HTTP      | `v1.Post("/path", forward.ToUserAction(c.usecaseContainer.SomeUseCase()))`                       |
| Consumer  | `kafka.NewConsumer(brokerCfg, cfg, forward.ToEventSubscriber(c.usecaseContainer.SomeUseCase()))` |
| AsyncTask | `worker.ForwardToAsyncTask(c.worker, c.usecaseContainer.SomeUseCase())`                          |
| CLI       | `c.usecaseContainer.SomeUseCase().Execute(ctx, input)`                                           |
