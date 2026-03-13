import { test, expect } from 'vitest';
import { action } from '../commands/set-state.js';

test('set-state requires dwp-state', async () => {
  await expect(action({})).rejects.toThrow(/Missing required flag: --dwp-state/);
});

test('set-state rejects working', async () => {
  await expect(action({ dwpState: 'working' })).rejects.toThrow(/use aynig set-working/);
});
