# List Manipulations

Guidelines for implementing list (collection) endpoints: filtering, pagination, sorting, search, joins, and response structure.

## Response Structure

All list responses **must** be wrapped in a struct with a `content` field — even if pagination is not used yet. This ensures backward compatibility when pagination or metadata fields are added later.

```json
{
  "content": [...]
}
```

For paginated lists, the response includes pagination metadata:

```json
{
  "page_number": 1,
  "page_size": 20,
  "count": 150,
  "content": [...]
}
```

**Never return a bare array** from a list endpoint. Always wrap in `{ "content": [] }` — including non-paginated lists (e.g., internal lookups, small bounded sets). When pagination is not needed, use a simple wrapper struct:

```go
type Response struct {
    Content []SomeDTO `json:"content"`
}
```

### Don't Over-Optimize Response Fields

Prefer returning full domain entities rather than creating dedicated response DTOs with cherry-picked fields. Even if the frontend doesn't need every field, the cost of extra bytes is far lower than the cost of maintaining separate response structs and mapping logic in the use case layer.

Create a dedicated response struct only when there is a real reason — e.g., joining data from multiple sources, computed fields, etc.

## Pagination

Pagination is handled at the **use case layer** using `rise-and-shine/pkg/pagination`.

### Use Case Request

Embed `pagination.Request` in the use case request struct:

```go
type Request struct {
    pagination.Request

    Name   *string `query:"name"`
    Status *string `query:"status"`
}
```

### Use Case Execute

Call `in.Normalize()` at the start of the use case to apply defaults (`page_number=1`, `page_size=20`) and constraints (`page_size` max 100):

```go
func (uc *usecase) Execute(ctx context.Context, in *Request) (*pagination.Response[SomeDTO], error) {
    // Normalize pagination params
    in.Normalize()

    // Build filter with pagination
    filter := entity.Filter{
        Name:   in.Name,
        Status: in.Status,
        Limit:  lo.ToPtr(in.PageSize()),
        Offset: lo.ToPtr(in.Offset()),
    }

    // Get items and total count
    items, count, err := uc.domainContainer.SomeRepo().ListWithCount(ctx, filter)
    if err != nil {
        return nil, errx.Wrap(err)
    }

    // Return paginated response
    return pagination.NewResponse(items, count, in.Request), nil
}
```

### Key Points

- `in.PageSize()` and `in.Offset()` compute limit/offset from page_number and page_size
- `pagination.NewResponse` builds the standard response with `page_number`, `page_size`, `count`, and `content`
- The Filter struct in the domain layer carries `Limit` and `Offset` as `*int` — the repo applies them if non-nil
- Use `ListWithCount` (provided by `repogen.Repo`) to get both items and total count in one call — it applies `Limit`/`Offset` for the items but counts without them

## Sorting

Sorting is handled using `rise-and-shine/pkg/sorter`. There are two distinct approaches for different purposes.

### Fixed-Field Sorting (Internal)

For use case internal logic where the sort column is hardcoded — not exposed to the API consumer:

```go
// Domain Filter — typed field for a specific column
type Filter struct {
    // ...
    OrderByLastUsedAt *sorter.SortDirection
}
```

```go
// Use case — hardcoded sort for internal logic
sessions, err := repo.List(ctx, session.Filter{
    ActorID:           &actorID,
    OrderByLastUsedAt: lo.ToPtr(sorter.Asc),
})
```

```go
// Filter function in infra
if f.OrderByLastUsedAt != nil {
    q = q.Order("last_used_at " + cast.ToString(f.OrderByLastUsedAt))
}
```

### Dynamic Sorting (User-Facing)

For list endpoints where the API consumer controls sorting. The request carries a raw `Sort` string, and the use case parses it with `sorter.MakeFromStr` using an allowlist of permitted fields.

**Request struct** — define `Sort` as a `string` query parameter:

```go
type Request struct {
    pagination.Request

    Name   *string `query:"name"`
    Status *string `query:"status"`
    Sort   string  `query:"sort"`
}
```

The API consumer sends sort as comma-separated `field:direction` pairs:

```
GET /api/v1/catalog/list-products?sort=name:asc,created_at:desc
```

**Use case** — parse and validate with `sorter.MakeFromStr`:

```go
func (uc *usecase) Execute(ctx context.Context, in *Request) (*pagination.Response[ProductDTO], error) {
    in.Normalize()

    filter := product.Filter{
        Name:     in.Name,
        Status:   in.Status,
        Limit:    lo.ToPtr(in.PageSize()),
        Offset:   lo.ToPtr(in.Offset()),
        SortOpts: sorter.MakeFromStr(in.Sort, "name", "created_at", "status"),
    }
    // ...
}
```

`MakeFromStr` handles everything:

- Parses the `field:direction` format
- Validates each field against the allowlist — unknown fields are silently dropped
- Validates direction is `asc` or `desc` — invalid directions are silently dropped
- Returns `nil` for empty string (no sorting applied)

**Domain Filter** — always include `SortOpts`:

```go
type Filter struct {
    // ...
    SortOpts sorter.SortOpts
}
```

**Filter function in infra** — apply `SortOpts` at the end:

```go
for _, o := range f.SortOpts {
    q = q.Order(o.ToSQL()) // produces "name asc", "created_at desc", etc.
}
```

### Sorting Rules

- Every Filter struct **must** include `SortOpts sorter.SortOpts`
- Fixed-field `OrderBy*` fields are checked **before** `SortOpts` in the filter function
- Always pass an allowlist to `MakeFromStr` — it is the validation mechanism, no additional validation needed
- Fixed-field sorting is for **internal use case logic**; `SortOpts` is for **API consumers**

## Filter Struct

The Filter struct is defined in the **domain layer** and applied in the **infra layer**. It is the single mechanism for all query customization (filtering, pagination, sorting).

### Structure Convention

```go
type Filter struct {
    // Exact-match filters (pointer = optional)
    ID       *string
    Username *string
    IsActive *bool

    // Multi-value filters
    IDs []int64

    // Pagination (set by use case from pagination.Request)
    Limit  *int
    Offset *int

    // Fixed-field sorting (optional, use case specific)
    OrderByCreatedAt *sorter.SortDirection

    // Dynamic sorting (always present)
    SortOpts sorter.SortOpts
}
```

### Filter Rules

- All single-value filter fields are **pointers** — `nil` means "don't filter by this field"
- Multi-value filters (e.g., `IDs []int64`) use slices — `nil` means "don't filter", an empty slice means "match nothing" (intentional). Guard with `if f.IDs != nil` not `if len(f.IDs) > 0`
- `Limit` and `Offset` are always `*int`, never bare `int`
- Filter structs belong in the **domain layer** alongside the repo interface
- Filter functions belong in the **infra layer** (`infra/postgres/`)

### Filter Function Pattern

The filter function maps Filter fields to SQL conditions. Follow this order:

```go
func entityFilterFunc(q *bun.SelectQuery, f entity.Filter) *bun.SelectQuery {
    // 1. WHERE conditions
    if f.ID != nil {
        q = q.Where("id = ?", *f.ID)
    }
    if f.IsActive != nil {
        q = q.Where("is_active = ?", *f.IsActive)
    }
    if f.IDs != nil {
        q = q.Where("id IN (?)", bun.In(f.IDs))
    }

    // 2. Pagination
    if f.Limit != nil {
        q = q.Limit(*f.Limit)
    }
    if f.Offset != nil {
        q = q.Offset(*f.Offset)
    }

    // 3. Fixed-field sorting
    if f.OrderByCreatedAt != nil {
        q = q.Order("created_at " + cast.ToString(f.OrderByCreatedAt))
    }

    // 4. Dynamic sorting (always last)
    for _, o := range f.SortOpts {
        q = q.Order(o.ToSQL())
    }

    return q
}
```

## Joining Rules

### Prefer Use-Case-Layer Merging

Repositories should stay **as simple as possible** — thin wrappers around repogen. When you need data from multiple tables, the default approach is to query each repo separately and merge results in the use case layer:

```go
// Use case — fetch from two repos within the same module, merge in code
orders, count, err := uc.domainContainer.OrderRepo().ListWithCount(ctx, orderFilter)
...

customerIDs := uniqueIDs(orders, func(o Order) string { return o.CustomerID })

customers, err := uc.domainContainer.CustomerRepo().List(ctx, customer.Filter{
    IDs: customerIDs,
})
...

return buildOrderDTOs(orders, customers, count, in.Request), nil
```

This keeps repositories generic and reusable. Most list endpoints can be built this way.

### Within a Module — Joins Allowed (When Justified)

SQL joins between tables **within the same module/schema** are allowed as a last resort — when use-case-layer merging would be impractical (e.g., filtering or sorting by a joined column, aggregations, performance-critical queries). Use custom repository methods:

```go
// Domain interface — add custom method
type Repo interface {
    repogen.Repo[Entity, Filter]

    ListWithDetails(ctx context.Context, f Filter) ([]EntityWithDetails, error)
}
```

```go
// Infra implementation
func (r *entityRepo) ListWithDetails(ctx context.Context, f Filter) ([]EntityWithDetails, error) {
    var results []EntityWithDetails
    q := r.idb.NewSelect().
        TableExpr(schemaName+".entities AS e").
        ColumnExpr("e.*").
        ColumnExpr("d.name AS detail_name").
        Join("LEFT JOIN "+schemaName+".details AS d ON e.detail_id = d.id")

    q = entityFilterFunc(q, f)

    err := q.Scan(ctx, &results)
    return results, errx.Wrap(err)
}
```

### Across Modules/Schemas — NOT Allowed

**Never join tables from different modules or schemas.** Each module owns its data and its schema. Cross-module data access goes through **portals**.

```
WRONG:  JOIN auth.admins ON catalog.products.created_by = auth.admins.id
RIGHT:  Call the auth module's portal to get admin data separately
```

If you need data from another module in a list:

1. Query your module's repo for the list
2. Collect the foreign IDs from the results
3. Call the other module's portal with those IDs to get the related data
4. Merge the data in the use case layer

```go
// Use case — combining data from two modules
products, err := uc.domainContainer.ProductRepo().List(ctx, filter)
...

// Collect unique creator IDs
creatorIDs := uniqueIDs(products, func(p Product) string { return p.CreatedBy })

// Get creator names from auth module via portal
creators, err := uc.portalContainer.Auth().GetAdminsByIDs(ctx, creatorIDs)
...

// Merge in use case layer
return mergeProductsWithCreators(products, creators), nil
```

### Why No Cross-Module Joins

- Modules must be independently deployable — cross-schema joins create tight coupling
- Each module owns its schema and data — even in database there is foreign keys, application shouldn't need to know about it
- Portal interfaces enforce a clean contract between modules
- This enables future module extraction into separate services

## Search

Search is a specialization of filtering — it uses the same Filter struct and filter function, with text-matching conditions instead of exact-match.

### Request and Filter

Define `Search` as a plain `string` (not `*string`) in both the request and filter structs. An empty string means "no search" — there's no point searching for `""`:

```go
// Use case request
type Request struct {
    pagination.Request

    Status *string `query:"status"`
    Search string  `query:"search"`
    Sort   string  `query:"sort"`
}
```

```go
// Domain filter
type Filter struct {
    Status *string

    Search string

    Limit    *int
    Offset   *int
    SortOpts sorter.SortOpts
}
```

### Filter Function

Apply search only when non-empty. Use `ILIKE` with `WhereGroup` to search across multiple columns with `OR`:

```go
if f.Search != "" {
    q = q.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
        return sq.
            Where("name ILIKE ?", "%"+f.Search+"%").
            WhereOr("description ILIKE ?", "%"+f.Search+"%")
    })
}
```

### Search Rules

- `Search` is a `string`, not `*string` — empty string means no search
- Always wrap the value with `%` in the filter function — never let the caller control the LIKE pattern
- Use `WhereGroup` with `OR` to match across multiple columns
- Search and exact-match filters coexist in the same Filter — they combine with `AND`
