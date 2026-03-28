import { readFileSync, writeFileSync } from 'node:fs';

const packageJsonPath = new URL('../package.json', import.meta.url);
const buildMetaPath = new URL('../build-meta.json', import.meta.url);
const pkg = JSON.parse(readFileSync(packageJsonPath, 'utf8'));
const timestamp = String(process.env.AYNIG_BUILD_TIMESTAMP_UNIX || Math.floor(Date.now() / 1000));

writeFileSync(buildMetaPath, `${JSON.stringify({ version: pkg.version, buildTimestampUnix: timestamp }, null, 2)}\n`);

process.stdout.write(`${timestamp}\n`);
