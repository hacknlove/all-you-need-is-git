import test from 'node:test';
import assert from 'node:assert/strict';
import { parseGitTrailersStrict } from '../gitHelpers/parseTrailers.js';

test('throws when raw is not a string', () => {
  assert.throws(
    () => parseGitTrailersStrict(null),
    /parseGitTrailersStrict: raw must be a string/
  );
});

test('returns empty object for empty input and trailing newlines', () => {
  assert.deepEqual(parseGitTrailersStrict(''), {});
  assert.deepEqual(parseGitTrailersStrict('\n\n\n'), {});
});

test('returns a null-prototype object', () => {
  const parsed = parseGitTrailersStrict('Token: value');
  assert.equal(Object.getPrototypeOf(parsed), null);
});

test('parses a single trailer with optional whitespace around separator and value', () => {
  const parsed = parseGitTrailersStrict('Signed-off-by  :   Example User <x@example.com>');

  assert.deepEqual({ ...parsed }, {
    'Signed-off-by': 'Example User <x@example.com>'
  });
});

test('parses repeated tokens as arrays preserving insertion order', () => {
  const parsed = parseGitTrailersStrict([
    'Reviewed-by: A',
    'Reviewed-by: B',
    'Reviewed-by: C'
  ].join('\n'));

  assert.deepEqual({ ...parsed }, {
    'Reviewed-by': ['A', 'B', 'C']
  });
});

test('parses continuation lines by appending newline and trimmed continuation text', () => {
  const parsed = parseGitTrailersStrict([
    'Notes: first line',
    '  second line',
    '\tthird line'
  ].join('\n'));

  assert.deepEqual({ ...parsed }, {
    Notes: 'first line\nsecond line\nthird line'
  });
});

test('continuation lines update the most recent repeated token value', () => {
  const parsed = parseGitTrailersStrict([
    'Co-authored-by: One',
    'Co-authored-by: Two',
    '  with continuation'
  ].join('\n'));

  assert.deepEqual({ ...parsed }, {
    'Co-authored-by': ['One', 'Two\nwith continuation']
  });
});

test('normalizes CRLF and CR newlines', () => {
  const parsed = parseGitTrailersStrict('A: 1\r\nB: 2\rC: 3\r\n');

  assert.deepEqual({ ...parsed }, {
    A: '1',
    B: '2',
    C: '3'
  });
});

test('supports lowerCaseKeys option and merges differently cased repeated tokens', () => {
  const parsed = parseGitTrailersStrict([
    'Acked-By: Alice',
    'acked-by: Bob'
  ].join('\n'), { lowerCaseKeys: true });

  assert.deepEqual({ ...parsed }, {
    'acked-by': ['Alice', 'Bob']
  });
});

test('supports custom separators and uses the earliest separator occurrence in a line', () => {
  const parsed = parseGitTrailersStrict('Key=left:right', {
    separators: [':', '=']
  });

  assert.deepEqual({ ...parsed }, {
    Key: 'left:right'
  });
});

test('allows empty values after a valid trailer token', () => {
  const parsed = parseGitTrailersStrict('Token:   ');

  assert.deepEqual({ ...parsed }, {
    Token: ''
  });
});

test('throws on blank lines inside trailer block', () => {
  assert.throws(
    () => parseGitTrailersStrict('A: 1\n\nB: 2'),
    /Invalid trailers: blank line at 2/
  );
});

test('throws on continuation line without preceding trailer', () => {
  assert.throws(
    () => parseGitTrailersStrict('  dangling continuation'),
    /Invalid trailers: continuation without a preceding trailer at 1/
  );
});

test('throws on lines that do not contain a valid token separator value form', () => {
  assert.throws(
    () => parseGitTrailersStrict('missing separator'),
    /Invalid trailers: expected "token:|=value" at 1/
  );
});

test('throws when separator appears at start of line (empty token)', () => {
  assert.throws(
    () => parseGitTrailersStrict(': value'),
    /Invalid trailers: expected "token:|=value" at 1/
  );
});
