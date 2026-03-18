import { test, expect, vi } from 'vitest';
import { registerEventsCommand } from '../commands/events.js';
import { git } from '../gitHelpers/git.js';

test('events uses dwp-state in output', async () => {
  const program = {
    command: () => program,
    description: () => program,
    option: () => program,
    action: (fn) => {
      program._action = fn;
      return program;
    }
  };

  registerEventsCommand(program);

  const rawSpy = vi.spyOn(git, 'raw').mockResolvedValue(
    'deadbeef\x1f2025-01-01T00:00:00Z\x1fsubject\n\nbody\n\ndwp-state: build\n\x1e'
  );

  const logSpy = vi.spyOn(console, 'log').mockImplementation(() => {});

  try {
    await program._action({ history: false, limit: '1', json: false });
    const output = logSpy.mock.calls.map(call => call.join(' ')).join('\n');
    expect(output).toContain('state=build');
  } finally {
    rawSpy.mockRestore();
    logSpy.mockRestore();
  }
});
