import { test, expect } from 'vitest';
import { action } from '../commands/set-state.js';

test('set-state requires aynig-state', async () => {
  await expect(action({})).rejects.toThrow(/Missing required flag: --aynig-state/);
});

test('set-state rejects working', async () => {
  await expect(action({ aynigState: 'working' })).rejects.toThrow(/use aynig set-working/);
});
