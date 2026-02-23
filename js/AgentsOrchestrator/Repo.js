import { git } from '../gitHelpers/git.js';
import { Branch } from './Branch.js';
import { Logger, resolveLevel } from '../utils/logger.js';

/**
 * Repo class to manage Git repositories.
 *
 * It:
 * 1. Initializes with a configuration object.
 * 2. Fetches branches (local by default, or remote if useRemote is set).
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
            if (!upstream.startsWith(`${this.config.useRemote}/`)) {
                this.config.logger.warn('Current branch upstream %s does not belong to remote %s', upstream, this.config.useRemote);
                return '';
            }
            return upstream;
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

        // Only fetch if using remote branches
        if (this.config.useRemote) {
            this.config.logger.info('Fetching remote branches from %s', this.config.useRemote);
            await git.fetch();
        }

        // Get remote branches if useRemote is set, otherwise get local branches
        const branchesInfo = this.config.useRemote
            ? await git.branch(['-r'])
            : await git.branchLocal();

        const localInfo = await git.branchLocal();
        const localCurrent = localInfo.current || '';

        const filterCurrent = this.config.useRemote
            ? await this.resolveCurrentRemoteBranch(localCurrent)
            : localCurrent;

        if (this.config.useRemote && this.config.currentBranch === 'only' && !filterCurrent) {
            this.config.logger.warn('Current branch has no upstream for --current-branch=only in remote mode');
        }

        const branchNames = this.filterBranches(branchesInfo.all, filterCurrent, this.config.currentBranch || 'skip');
        this.config.logger.info('Running %d branches (current-branch=%s)', branchNames.length, this.config.currentBranch || 'skip');

        this.branches = branchNames.map(name => new Branch({
            config: this.config,
            branchName: name,
            isCurrentBranch: !this.config.useRemote && name === localCurrent,
        }));

        await Promise.all(this.branches.map(branch => branch.run()));
    }
}
