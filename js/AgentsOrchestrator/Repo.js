import simpleGit from 'simple-git';
import { Branch } from './Branch.js';
const defaults = {
    remote: 'origin',
}

/**
 * Repo class to manage Git repositories.
 * 
 * It:
 * 1. Initializes with a configuration object.
 * 2. Fetches branches from the remote repository.
 * 3. Creates Branch instances for each remote branch.
 * 4. Runs each Branch instance.
 * 5. waits for all Branch instances to complete.
 * 
 * It has no clue about what each Branch does internally.
 */
export class Repo {
    constructor(config) {
        this.config = {
            ...defaults,
            ...config
        };
        this.git = simpleGit({ baseDir: this.config.dir });
        this.hotBranches = [];
    }
    

    async run() {
        await this.git.fetch();
        const branchesName = await this.git.branch(['-r']);

        this.branches = branchesName.all.map(branchesName => new Branch({
            git: this.git,
            config: this.config,
            branchName: branchesName,
        }))

        await Promise.all(this.branches.map(branch => branch.run()));
    }
}

