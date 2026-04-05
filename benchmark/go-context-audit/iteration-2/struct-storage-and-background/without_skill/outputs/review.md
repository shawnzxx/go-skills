## 总览

该文件存在**明显的 context 误用**：把请求/调用方传入的 `ctx` **存进结构体**、并在查询时**丢弃**它，改用 `context.Background()`。

## 发现与风险评估

### 1) 在结构体中持久化 `context.Context`
- **风险等级**：高
- **位置**：`type UserRepository struct { ctx context.Context ... }`，`NewUserRepository(ctx, db)`

### 2) 查询时使用 `context.Background()` 丢弃调用方上下文
- **风险等级**：高
- **位置**：`LoadUser` 中 `queryCtx := context.Background()`

### 3) API 设计隐藏了 context 传播需求
- **风险等级**：中
- **位置**：`func (r *UserRepository) LoadUser(id string) error`
