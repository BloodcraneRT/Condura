## 2024-03-24 - Pre-allocate dummy payloads for synthetic network tests
**Learning:** High-concurrency synthetic testing endpoints (like download endpoints or BERT tests) that generate massive dummy payloads on the fly create substantial GC pressure if the payload buffers are re-allocated per request.
**Action:** Always pre-allocate static dummy payloads into read-only global slices during init() to completely eliminate per-request memory allocations for these tests.
