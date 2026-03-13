import { test, expect } from 'vitest';
import { resolveDwpRemote } from '../commands/stateCommon.js';

test('resolveDwpRemote uses git: locator', () => {
  const trailers = { 'dwp-source': 'git:origin' };
  expect(resolveDwpRemote('', trailers)).toBe('origin');
});

test('resolveDwpRemote ignores non-git locators', () => {
  const trailers = { 'dwp-source': 'http:https://example.com' };
  expect(resolveDwpRemote('', trailers)).toBe('');
});
