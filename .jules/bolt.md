## 2024-05-18 - math/bits Provides Significant Performance Boost
**Learning:** Replacing manual bitwise counting loops with standard library functions from `math/bits` (like `bits.OnesCount8`) provides a significant performance boost for calculating bit errors in large payloads (~14x faster in my benchmarks) due to the use of optimized assembly instructions.
**Action:** When working with bit-level analysis or manipulations across large buffers, always prefer `math/bits` library functions over manual loop structures.

## 2026-04-08 - Streaming Payloads Reduces Memory Overhead in Synthetic Tests
**Learning:** Creating dummy payload bytes in memory (e.g., 10MB byte slices) for synthetic HTTP upload tests generates significant memory allocations and GC overhead when tests run frequently at the edge.
**Action:** Use a custom `io.Reader` implementation to stream synthetic dummy bytes incrementally directly into HTTP requests instead of pre-allocating large byte arrays.
