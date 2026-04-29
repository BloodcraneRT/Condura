## 2024-04-29 - Prevent Accidental Deletion
**Learning:** Destructive actions like deleting tasks happen instantaneously without user confirmation, leading to accidental data loss because the delete button appears on hover in a dense list.
**Action:** Always wrap destructive UI actions (such as deleting tasks or configurations) in a confirmation dialog (e.g., `window.confirm`) to prevent accidental data loss and maintain user trust.
