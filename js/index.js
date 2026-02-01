import config from './config.json' with { type: 'json' };
import { Repo } from './AgentsOrchestrator/Repo.js';

/**
 * Main entry point
 * 
 * It:
 * 1. instantiates Repo objects for each repository defined in the config
 * 2. runs them
 * 3. waits for all to complete
 * 
 * It has no clue about what each Repo does internally.
 */
async function main() {
  const promises = []
  for (const repo of config.repos) {
    const repoConfig = {
      ...config.global,
      ...repo,
    };

    const repository = new Repo(repoConfig);
    promises.push(repository.run());
  }

  await Promise.all(promises);
}

main();