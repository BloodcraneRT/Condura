## 2024-03-24 - Batching Client Polling
**Learning:** The dashboard previously made 4 separate API calls every 5 seconds, creating unnecessary network overhead and sequential DB querying blocks.
**Action:** Always look for opportunities to aggregate data required by a single view into a single API endpoint that loads dependencies concurrently on the backend.
