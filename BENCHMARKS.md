# Policy System Benchmarks

## Running Benchmarks

### All Middleware Benchmarks
```bash
go test -bench=. -benchmem -benchtime=1s ./internal/gateway/middleware
```

### All Policy Benchmarks
```bash
go test -bench=. -benchmem -benchtime=1s ./internal/gateway/policies
```

### Specific Benchmarks
```bash
# Model extraction
go test -bench=BenchmarkModelExtraction -benchmem ./internal/gateway/middleware

# Token estimation
go test -bench=BenchmarkTokenEstimation -benchmem ./internal/gateway/middleware

# JSON unmarshal overhead
go test -bench=BenchmarkJSONUnmarshal -benchmem ./internal/gateway/middleware

# Context access patterns
go test -bench=BenchmarkContextAccess -benchmem ./internal/gateway/middleware

# Policy checks
go test -bench=BenchmarkTokenLimitPolicy -benchmem ./internal/gateway/policies
go test -bench=BenchmarkModelAllowlistPolicy -benchmem ./internal/gateway/policies
go test -bench=BenchmarkCELPolicy -benchmem ./internal/gateway/policies
```

### Compare Before/After Changes
```bash
# Save baseline
go test -bench=. -benchmem ./internal/gateway/middleware > old.txt

# Make changes...

# Compare
go test -bench=. -benchmem ./internal/gateway/middleware > new.txt
benchstat old.txt new.txt
```

Install benchstat:
```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

## Current Baseline (Before Optimizations)

### Model Extraction
```
BenchmarkModelExtraction/1KB-11     20733    5555 ns/op    288 B/op    8 allocs/op
BenchmarkModelExtraction/10KB-11     2308   51626 ns/op    288 B/op    8 allocs/op
BenchmarkModelExtraction/50KB-11      465  257145 ns/op    288 B/op    8 allocs/op
```

**Observations:**
- Time scales linearly with request size (JSON parsing overhead)
- ~5μs for 1KB, ~50μs for 10KB, ~250μs for 50KB
- 8 allocations per parse (struct + fields)

### Expected After Optimization

After implementing Critical Fix #1 (parse once, store in context):
- Model extraction: 0 ns/op (already parsed)
- Token estimation: Use pre-parsed data
- Total savings: 3x JSON unmarshal eliminated

## Performance Targets

**Per-Request Overhead:**
- Total policy overhead: < 1ms (P50), < 2ms (P99)
- Memory overhead: < 5KB per request
- Rate limiter: < 100μs per check

**Throughput:**
- Handle 1000 req/sec on single instance
- No memory growth over 24 hours
- Rate limiter accuracy: 99%+

## Profiling

### CPU Profile
```bash
go test -bench=BenchmarkPolicyEnforcement -cpuprofile=cpu.prof
go tool pprof -http=:8080 cpu.prof
```

### Memory Profile
```bash
go test -bench=BenchmarkPolicyEnforcement -memprofile=mem.prof
go tool pprof -http=:8080 mem.prof
```

### Trace
```bash
go test -bench=BenchmarkPolicyEnforcement -trace=trace.out
go tool trace trace.out
```

## Load Testing

For end-to-end load testing with real HTTP requests:

```bash
# Install k6
brew install k6

# Run load test
k6 run scripts/load-test.js
```

(Note: load test script to be created)
