import simpleGit from 'simple-git';

const git = simpleGit();

async function main() {
  try {
    // Get current status
    const status = await git.status();
    console.log('Current branch:', status.current);
    console.log('Modified files:', status.modified);
    console.log('Untracked files:', status.not_added);

    // Example: Add files
    // await git.add('./*');

    // Example: Commit changes
    // await git.commit('Your commit message');

    // Example: Push to remote
    // await git.push('origin', 'master');

    // Example: Fetch from remote
    // await git.fetch();

    // Example: Pull changes
    // await git.pull();

    console.log('Git operations completed successfully');
  } catch (error) {
    console.error('Git error:', error);
  }
}

main();