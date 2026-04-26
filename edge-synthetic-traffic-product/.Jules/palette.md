## 2026-04-26 - Destructive Action Confirmations
**Learning:** Destructive actions like deleting a task were executing immediately upon button click, which could easily lead to accidental data loss.
**Action:** Always wrap destructive UI actions in a confirmation dialog (e.g., `window.confirm`) to introduce a safe friction point before executing irreversible actions.
