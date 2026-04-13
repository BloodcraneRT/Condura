## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2024-05-18 - Streaming vs Allocating Large Mock Payloads
**Learning:** Pre-allocating large byte arrays (e.g., `make([]byte, 10 * 1024 * 1024)`) simply to represent dummy payload data for network tests significantly increases memory footprint and triggers excessive Garbage Collection (GC) pauses.
**Action:** When generating large payloads for `http.NewRequest` or similar test constructs, implement a custom `io.Reader` (e.g. streaming a constant character) to simulate the payload without ever allocating a large backing array. Remember to explicitly set `req.ContentLength` since it cannot be inferred automatically from custom streaming readers.
