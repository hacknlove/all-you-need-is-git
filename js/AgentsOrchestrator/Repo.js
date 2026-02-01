import { git} from '../gitHelpers/git.js';
import { Branch } from './Branch.js';

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


    async run() {
        // Only fetch if using remote branches
        if (this.config.useRemote) {
            await git.fetch();
        }

        // Get remote branches if useRemote is set, otherwise get local branches
        const branchesName = this.config.useRemote
            ? await git.branch(['-r'])
            : await git.branchLocal();

        this.branches = branchesName.all.map(branchesName => new Branch({
            config: this.config,
            branchName: branchesName,
        }))

        await Promise.all(this.branches.map(branch => branch.run()));
    }
}

