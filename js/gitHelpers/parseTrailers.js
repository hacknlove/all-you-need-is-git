/**
 * Strict Git-trailer parser.
 *
 * Assumptions:
 * - `raw` contains ONLY a valid trailer block (one paragraph).
 * - No blank lines inside the block (a trailing newline at EOF is OK).
 * - Continuation lines (leading whitespace) continue the previous trailer value.
 * - Repeated tokens become arrays.
 *
 * Returns: { [token]: string | string[] }
 * Throws on invalid input.
 */
export function parseGitTrailersStrict(raw, opts = {}) {
  const {
    separators = [":", "="],
    lowerCaseKeys = false
  } = opts;

  if (typeof raw !== "string") {
    throw new TypeError("parseGitTrailersStrict: raw must be a string");
  }

  // Normalize newlines, allow trailing newline(s)
  const normalized = raw.replace(/\r\n/g, "\n").replace(/\r/g, "\n");
  const lines = normalized.split("\n");

  // Drop final empty lines (common trailing newline at EOF)
  while (lines.length && lines[lines.length - 1] === "") lines.pop();

  if (lines.length === 0) return {};

  const out = Object.create(null);

  let currentKey = null;

  function addValue(key, value) {
    if (lowerCaseKeys) key = key.toLowerCase();

    if (Object.prototype.hasOwnProperty.call(out, key)) {
      const cur = out[key];
      if (Array.isArray(cur)) cur.push(value);
      else out[key] = [cur, value];
    } else {
      out[key] = value;
    }
  }

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];

    // Strict: no blank lines inside trailer block
    if (line.trim() === "") {
      throw new Error(`Invalid trailers: blank line at ${i + 1}`);
    }

    // Continuation line
    if (/^\s+/.test(line)) {
      if (!currentKey) {
        throw new Error(`Invalid trailers: continuation without a preceding trailer at ${i + 1}`);
      }

      // Append as a newline + trimmed continuation (remove leading indentation)
      const cont = line.replace(/^\s+/, "");
      // Update the last value for currentKey
      const cur = out[lowerCaseKeys ? currentKey.toLowerCase() : currentKey];

      if (Array.isArray(cur)) {
        cur[cur.length - 1] = cur[cur.length - 1] + "\n" + cont;
      } else {
        out[lowerCaseKeys ? currentKey.toLowerCase() : currentKey] = cur + "\n" + cont;
      }
      continue;
    }

    // Trailer line: find earliest separator occurrence
    let sepIdx = -1;
    let sep = null;
    for (const s of separators) {
      const idx = line.indexOf(s);
      if (idx !== -1 && (sepIdx === -1 || idx < sepIdx)) {
        sepIdx = idx;
        sep = s;
      }
    }

    if (sepIdx <= 0) {
      throw new Error(`Invalid trailers: expected "token${separators.join("|")}value" at ${i + 1}: "${line}"`);
    }

    const key = line.slice(0, sepIdx).trim();
    if (!key) {
      throw new Error(`Invalid trailers: empty token at ${i + 1}`);
    }

    // Git allows optional whitespace after separator
    const value = line.slice(sepIdx + sep.length).replace(/^\s*/, "");

    currentKey = key;
    addValue(key, value);
  }

  return out;
}