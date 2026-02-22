import { git } from '../gitHelpers/git.js';
import { Branch } from './Branch.js';
import { Logger } from '../utils/logger.js';

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

    filterBranches(branchesInfo, mode) {
        switch (mode) {
            case 'skip':
                return branchesInfo.all.filter(name => name !== branchesInfo.current); 
            case 'only':
                return [branchesInfo.current];
            default:
                return branchesInfo.all; 
        }
    }

    async run() {
        if (!this.config.logger) {
            this.config.logger = new Logger(this.config.logLevel);
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

        const branchNames = this.filterBranches(branchesInfo, this.config.currentBranch || 'skip');
        this.config.logger.info('Running %d branches (current-branch=%s)', branchNames.length, this.config.currentBranch || 'skip');

        this.branches = branchNames.map(name => new Branch({
            config: this.config,
            branchName: name,
            isCurrentBranch: name === branchesInfo.current,
        }));

        await Promise.all(this.branches.map(branch => branch.run()));
    }
}
