## 2024-04-05 - Synthetic Traffic Payload Allocation
**Learning:** In a synthetic traffic generator that runs tasks on a frequent loop, repeatedly allocating large dummy payloads (like a 10MB byte slice for upload tests) inside the execution loop causes massive, unnecessary Garbage Collection (GC) pressure and CPU overhead, defeating the goal of a high-performance edge binary.
**Action:** Always pre-allocate static/dummy payloads as global read-only variables initialized once at startup, allowing concurrent tests to share the same memory safely using `bytes.NewReader`.
