/**
 * Install command - installs AYNIG workflows from another repository
 */
async function action(repo, ref, subfolder) {
  console.log('Hello from install command!');
  console.log(`Repo: ${repo}`);
  console.log(`Ref: ${ref}`);
  console.log(`Subfolder: ${subfolder || '(none)'}`);
  // TODO: Implement install logic
  // - clone/fetch the specified repository
  // - copy .aynig/ scripts from source to current repo
}

export function registerInstallCommand(program) {
  program
    .command('install <repo> <ref> [subfolder]')
    .description('Install AYNIG workflows from another repository')
    .action(action);
}
