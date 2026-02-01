export function parseCommitMessage(log) {
    if (!log) {
        return {
            firstLine: '',
            body: '',
            trailers: {}
        };
    }

    const firstLine = log.message || '';
    const body = log.body || '';

    // Parse trailers from the body
    const trailers = {};
    const lines = body.trim().split('\n');

    // Trailers are at the end of the commit message
    // They follow the format "Key: value" or "Key-Name: value"
    const trailerPattern = /^([A-Z][A-Za-z0-9-]*)\s*:\s*(.+)$/;

    // Start from the end and collect trailers
    let foundTrailers = false;
    for (let i = lines.length - 1; i >= 0; i--) {
        const line = lines[i].trim();

        if (!line) {
            // Empty line might separate trailers from body
            if (foundTrailers) {
                break;
            }
            continue;
        }

        const match = line.match(trailerPattern);
        if (match) {
            const [, key, value] = match;
            trailers[key] = value.trim();
            foundTrailers = true;
        } else if (foundTrailers) {
            // If we've found trailers and hit a non-trailer line, stop
            break;
        }
    }

    return {
        firstLine,
        body,
        trailers
    };
}
