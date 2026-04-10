## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2024-11-09 - Avoid Pre-allocating Large Static Byte Slices for Synthetic Payloads
**Learning:** Pre-allocating large byte arrays (e.g., `make([]byte, 10 * 1024 * 1024)`) on every execution of a synthetic test runner causes significant garbage collection (GC) pressure and unnecessary memory usage.
**Action:** Use a custom `io.Reader` implementation to stream synthetic dummy payloads on the fly instead. When passing this custom reader to `http.NewRequest`, you must explicitly set `req.ContentLength` since the Go standard library cannot automatically infer it.
