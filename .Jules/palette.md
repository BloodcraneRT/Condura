## 2024-05-18 - Hover-Only Element Keyboard Accessibility
**Learning:** Hover-only elements (like the task 'Delete' button) with `opacity-0 group-hover:opacity-100` are completely inaccessible to keyboard-only users as they cannot visually see where their focus is.
**Action:** Always pair `opacity-0` hover effects on interactive elements with `focus-visible:opacity-100` and clear focus rings (e.g., `focus-visible:ring-2`) to ensure keyboard users can discover and use the element.
## 2026-04-07 - Accessible Forms & Keyboard Navigation

**Learning:** Forms using labels without explicit `htmlFor` attributes to associate with their target `id` create a barrier for screen readers. Further, actions hidden behind `opacity-0` hover states completely fail for keyboard navigation users who rely on tabbing, making core functionality unreachable.

**Action:** Always pair `<label htmlFor="id">` with `<input id="id">`. Ensure that any interaction hidden behind a mouse hover also includes a `focus-visible` override (e.g., `focus-visible:opacity-100`) combined with visible focus rings (`focus-visible:ring-2`) and keyboard outlines disabled (`focus-visible:outline-none`) to safely expose functionality to keyboard-only users.
## 2024-05-19 - Actionable Empty States
**Learning:** Plain text empty states ("No tasks currently scheduled") are unhelpful to users and miss an opportunity to guide them. Adding visual weight (dashed borders, icons) and a call-to-action button that programmatically focuses the relevant input form (`document.getElementById("id").focus()`) significantly reduces friction for first-time users and improves accessibility.
**Action:** When designing empty states for lists or tables, always include an actionable button that directs focus to the element needed to populate that list, and ensure the icon used is decorative with `aria-hidden="true"`.

## 2024-04-16 - Prevent extraneous UI dependencies during verifications
**Learning:** Adding Playwright directly to the project's `package.json` for a quick automated check inflates the project lockfile and adds heavy dev dependencies unnecessarily, especially since tests weren't formally added to the repo.
**Action:** When installing temporary tooling like Playwright (`pnpm add -D playwright`) to verify a UI change via script, immediately discard the changes to `package.json` and `pnpm-lock.yaml` using `git restore --staged` and `git checkout` after verification is complete.

## 2024-04-16 - Required indicators must be paired with HTML attributes
**Learning:** Adding visual required indicators (e.g. `<span aria-hidden="true">*</span>`) to form labels improves visual accessibility, but screen readers and native form validation will fail to enforce it unless the corresponding input also has the `required` HTML attribute.
**Action:** Always ensure any `<input>` or `<select>` element matched to a "required" label explicitly includes the `required` attribute.
