## Finding 1
【Risk】High  
【Code And Explanation】  
`auditfixture/user_repository.go:8-21`  

`UserRepository` 把 `ctx context.Context` 存进了 struct 字段，并在构造函数里接收请求级 `ctx`。这会把**请求生命周期**“固化”到长生命周期对象里（例如被复用/缓存/跨请求共享），带来以下风险：  
- **请求取消/超时无法可靠生效**（后续调用可能用到早已结束的 ctx）。  
- **潜在泄漏**：ctx 可能携带 request-scoped 值/trace/span 等，长时间持有会延长其存活。  
- 设计上也会鼓励方法不显式接收 ctx，导致传播链变差、难审计。

【Suggested Fix】
```go
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) LoadUser(ctx context.Context, id string) error {
	return r.db.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1", id).Err()
}
```

## Finding 2
【Risk】High  
【Code And Explanation】  
`auditfixture/user_repository.go:23-26`  

`LoadUser` 里显式创建了 `queryCtx := context.Background()`，这会**丢弃上游请求 ctx**（取消、deadline、trace 等全部断开）。

【Suggested Fix】
```go
func (r *UserRepository) LoadUser(ctx context.Context, id string) error {
	return r.db.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1", id).Err()
}
```
