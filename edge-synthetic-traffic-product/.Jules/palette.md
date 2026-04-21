## 2024-05-21 - Destructive Actions Confirmation
**Learning:** Destructive actions like task deletion hidden behind a single click (especially on hover states) can lead to unintentional data loss and frustrate users.
**Action:** When implementing destructive UI actions (such as deleting tasks or configurations), always wrap the action in a confirmation dialog (e.g., `window.confirm`) to prevent accidental data loss and adhere to the project's UX standards.
