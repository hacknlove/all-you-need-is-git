/**
 * Converts a git commit hash to a UUID format.
 * Takes the first 128 bits (32 hex characters) of the commit hash
 * and formats it as a standard UUID.
 *
 * @param {string} hash - Git commit hash (SHA-1 or SHA-256)
 * @returns {string} UUID formatted string (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)
 *
 * Collision probability is negligible:
 * - 1 million commits: ~1 in 10^29 chance
 * - 1 billion commits: ~1 in 10^19 chance
 */
export function commitHashToUUID(hash) {
    if (!hash || hash.length < 32) {
        throw new Error('Invalid commit hash: must be at least 32 characters');
    }

    const hex = hash.substring(0, 32);
    return `${hex.substring(0, 8)}-${hex.substring(8, 12)}-${hex.substring(12, 16)}-${hex.substring(16, 20)}-${hex.substring(20, 32)}`;
}
