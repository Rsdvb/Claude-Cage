import assert from "node:assert/strict";
import test from "node:test";
import { reconcileTabs } from "./workspace-state.mjs";

test("polling reuses unchanged tabs and preserves user order", () => {
  const first = { id: "browser", title: "Browser" };
  const second = { id: "terminal", title: "Terminal" };
  const previous = [second, first];
  assert.strictEqual(reconcileTabs(previous, [{ ...first }, { ...second }]), previous);

  const changed = { ...second, title: "Updated terminal" };
  const added = { id: "new", title: "New" };
  const current = reconcileTabs(previous, [{ ...first }, changed, added]);
  assert.deepEqual(current, [changed, first, added]);
  assert.strictEqual(current[1], first);
  assert.deepEqual(reconcileTabs(current, [changed, added]), [changed, added]);
});
