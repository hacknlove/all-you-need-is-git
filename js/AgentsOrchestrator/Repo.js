import { git } from '../gitHelpers/git.js';
import { Branch } from './Branch.js';
import { Logger, resolveLevel } from '../utils/logger.js';
import { parseCommitMessage } from '../gitHelpers/parseCommitMessage.js';

/**
 * Repo class to manage Git repositories.
 *
 * It:
 * 1. Initializes with a configuration object.
 * 2. Fetches branches (local by default, or remote if aynigRemote is set).
 * 3. Creates Branch instances for each branch.
 * 4. Runs each Branch instance.
 * 5. waits for all Branch instances to complete.
 *
 * It has no clue about what each Branch does internally.
 */
export class Repo {
    constructor(config) {
        this.config = config;
    }

    filterBranches(branches, current, mode) {
        switch (mode) {
            case 'skip':
                return current ? branches.filter(name => name !== current) : branches;
            case 'only':
                return current ? [current] : [];
            default:
                return branches;
        }
    }

    async resolveCurrentRemoteBranch(localCurrent) {
        if (!localCurrent) {
            return '';
        }
        try {
            const upstream = (await git.raw(['rev-parse', '--abbrev-ref', `${localCurrent}@{upstream}`])).trim();
            if (!upstream) {
                return '';
            }
            if (!upstream.startsWith(`${this.config.aynigRemote}/`)) {
                this.config.logger.warn('Current branch upstream %s does not belong to remote %s', upstream, this.config.aynigRemote);
                return '';
            }
            return upstream;
        } catch {
            return '';
        }
    }

    remoteFromTrailers(trailers = {}) {
        const raw = trailers['aynig-remote'];
        if (Array.isArray(raw)) {
            return String(raw[0] || '').trim();
        }
        return String(raw || '').trim();
    }

    async resolveAynigRemoteFromHead() {
        if (this.config.aynigRemote) {
            return this.config.aynigRemote;
        }

        try {
            const raw = await git.raw(['log', '-1', '--pretty=format:%s%x1f%b']);
            const parts = raw.split('\x1f');
            const message = (parts[0] || '').replace(/\n+$/, '');
            const body = (parts[1] || '').replace(/\n+$/, '');
            const parsed = await parseCommitMessage({ message, body });
            return this.remoteFromTrailers(parsed.trailers);
        } catch {
            return '';
        }
    }

    async run() {
        if (!this.config.logger) {
            const baseLevel = resolveLevel({
                cliLevel: this.config.logLevel,
                cliSet: this.config.logLevelSet,
                trailerLevel: '',
                envLevel: process.env.AYNIG_LOG_LEVEL,
                defaultLevel: this.config.logLevel
            });
            this.config.logger = new Logger(baseLevel);
        }

        // Resolve repo root once and share it via config
        this.config.repoRoot = await git.revparse(['--show-toplevel']);
        this.config.logger.info('Repository root: %s', this.config.repoRoot.trim());

        this.config.aynigRemote = await this.resolveAynigRemoteFromHead();
        if (this.config.aynigRemote) {
            this.config.logger.info('Fetching remote branches from %s', this.config.aynigRemote);
            await git.fetch();
        }

        // Get remote branches if aynigRemote is set, otherwise get local branches
        const branchesInfo = this.config.aynigRemote
            ? await git.branch(['-r'])
            : await git.branchLocal();

        const localInfo = await git.branchLocal();
        const localCurrent = localInfo.current || '';

        const filterCurrent = this.config.aynigRemote
            ? await this.resolveCurrentRemoteBranch(localCurrent)
            : localCurrent;

        if (this.config.aynigRemote && this.config.currentBranch === 'only' && !filterCurrent) {
            this.config.logger.warn('Current branch has no upstream for --current-branch=only in remote mode');
        }

        const branchNames = this.filterBranches(branchesInfo.all, filterCurrent, this.config.currentBranch || 'skip');
        this.config.logger.info('Running %d branches (current-branch=%s)', branchNames.length, this.config.currentBranch || 'skip');

        this.branches = branchNames.map(name => new Branch({
            config: this.config,
            branchName: name,
            isCurrentBranch: !this.config.aynigRemote && name === localCurrent,
        }));

        await Promise.all(this.branches.map(branch => branch.run()));
    }
}
