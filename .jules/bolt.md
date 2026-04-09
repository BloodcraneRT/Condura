## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2024-05-18 - Avoid Large Pre-allocated Slices in Synthetic Payloads
**Learning:** In synthetic test runners (like the `runUploadTest` in `runner.go`), pre-allocating large byte arrays (e.g., 10MB) for dummy payloads creates significant memory overhead and garbage collection pressure, especially when many tests run concurrently.
**Action:** Use a custom `io.Reader` implementation to stream dummy bytes continuously without pre-allocating large slices. When using a custom `io.Reader` with `http.NewRequest`, you must manually set `req.ContentLength` since it cannot be inferred, otherwise the Go HTTP client will use chunked transfer encoding which may alter test validity.
