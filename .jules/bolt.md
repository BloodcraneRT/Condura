## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.
## 2024-05-18 - Avoid Large Byte Allocations in Synthetic Payloads
**Learning:** Pre-allocating large byte arrays (e.g., `make([]byte, 10*1024*1024)`) repeatedly inside test runner functions like `runUploadTest` causes significant memory spikes and GC pressure. When using a custom `io.Reader` to stream these payloads to `http.NewRequest`, you must explicitly set `req.ContentLength` to avoid unintended chunked transfer encoding behavior.
**Action:** Use a custom streaming `io.Reader` implementation for large, repetitive synthetic payloads to reduce memory overhead, and always explicitly set `ContentLength` when passing it to HTTP requests.
