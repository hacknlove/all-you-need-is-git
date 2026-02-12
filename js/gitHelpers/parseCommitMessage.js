import { spawn } from 'child_process';
import { parseGitTrailersStrict } from './parseTrailers.js'; 

function interpretTrailers(body) {
    return new Promise((resolve, reject) => {
        const child = spawn('git', ['interpret-trailers', '--parse', '--only-trailers']);
        let stdout = '';
        let stderr = '';

        child.stdout.on('data', chunk => {
            stdout += chunk;
        });

        child.stderr.on('data', chunk => {
            stderr += chunk;
        });

        child.on('error', reject);

        child.on('close', code => {
            if (code === 0) {
                resolve(stdout);
                return;
            }
            const suffix = stderr.trim() ? `: ${stderr.trim()}` : '';
            reject(new Error(`git interpret-trailers failed with code ${code}${suffix}`));
        });

        child.stdin.on('error', error => {
            if (error?.code !== 'EPIPE') {
                reject(error);
            }
        });

        child.stdin.end(body);
    });
}

export async function parseCommitMessage(log) {
    if (!log) {
        return {
            firstLine: '',
            body: '',
            trailers: {}
        };
    }

    const firstLine = log.message || '';

    if (!log.body) {
        return {
            firstLine,
            body: '',
            trailers: {}
        };
    }

    const trailersRaw = await interpretTrailers(log.body);

    const trailers = parseGitTrailersStrict(trailersRaw);

    return {
        firstLine,
        body: log.body,
        trailers
    };
}
