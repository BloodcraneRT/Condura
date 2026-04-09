## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2024-05-20 - Streaming Dummy Payloads for Requests Reduces GC Pressure
**Learning:** When generating large dummy payloads for HTTP requests (like in `runUploadTest`), pre-allocating large byte arrays (e.g. `make([]byte, 10*1024*1024)`) causes significant memory allocation and garbage collection overhead. Using a custom `io.Reader` to stream the dummy payload eliminates this allocation entirely. Note that when passing a custom reader to `http.NewRequest`, you must explicitly set `req.ContentLength` since it cannot be automatically inferred.
**Action:** Always implement and use custom streaming `io.Reader` implementations for synthetic load generation instead of static array allocations to optimize memory usage and reduce GC pressure.
