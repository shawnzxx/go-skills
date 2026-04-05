## 发现 1：后台 goroutine 永不退出（高风险）

- **风险**：高  
- **位置**：`(*Processor).Start` 中 `go func(){ ... for { select { ... } } }`（整个 goroutine 循环）
- **问题**：
  - 该 goroutine **没有监听 `ctx.Done()`**。

- **建议修复**：

```go
for {
  select {
  case <-ctx.Done():
    return
  case job, ok := <-jobs:
    if !ok {
      return
    }
    p.handle(job)
  case <-ticker.C:
    p.flush()
  }
}
```

## 发现 2：从 `jobs` 接收未处理关闭语义，可能导致忙等/空转（高风险）
- **位置**：`case job := <-jobs:`
