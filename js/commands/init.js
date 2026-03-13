import fs from 'fs/promises';
import path from 'path';
import { fileURLToPath } from 'url';
import simpleGit from 'simple-git';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const ASSETS_DIR = path.join(__dirname, '..', 'assets');

/**
 * Init command - initializes AYNIG in the current repository
 */
async function action() {
  const git = simpleGit();

  // Check if we're in a git repository
  try {
    await git.status();
  } catch (error) {
    console.error('Error: Not a Git repository. Please run `git init` first.');
    process.exit(1);
  }

  // Create .dwp/ directory
  const dwpDir = '.dwp';
  let dwpCreated = false;
  try {
    await fs.mkdir(dwpDir, { recursive: false });
    console.log(`✓ Created ${dwpDir}/`);
    dwpCreated = true;
  } catch (error) {
    if (error.code === 'EEXIST') {
      console.log(`⊘ ${dwpDir}/ already exists, skipping`);
    } else {
      throw error;
    }
  }

  // Ensure command directory exists
  const commandDir = path.join(dwpDir, 'command');
  try {
    await fs.mkdir(commandDir, { recursive: false });
    console.log(`✓ Created ${commandDir}/`);
  } catch (error) {
    if (error.code === 'EEXIST') {
      console.log(`⊘ ${commandDir}/ already exists, skipping`);
    } else {
      throw error;
    }
  }

  // Only create files if directory was just created
  if (dwpCreated) {
    // Copy COMMANDS.md
    await fs.copyFile(
      path.join(ASSETS_DIR, 'COMMANDS.md'),
      path.join(dwpDir, 'COMMANDS.md')
    );
    console.log(`✓ Created COMMANDS.md`);
  }

  // Ensure CONTRACT.md exists
  const contractPath = path.join(dwpDir, 'CONTRACT.md');
  try {
    await fs.access(contractPath);
    console.log('⊘ CONTRACT.md already exists, skipping');
  } catch {
    await fs.copyFile(
      path.join(ASSETS_DIR, 'CONTRACT.md'),
      path.join(dwpDir, 'CONTRACT.md')
    );
    console.log('✓ Created CONTRACT.md');
  }

  // Ensure clean command exists
  const cleanPath = path.join(commandDir, 'clean');
  try {
    await fs.access(cleanPath);
    console.log('⊘ clean command already exists, skipping');
  } catch {
    await fs.copyFile(path.join(ASSETS_DIR, 'clean'), cleanPath);
    await fs.chmod(cleanPath, 0o755);
    console.log('✓ Created clean command');
  }

  // Create .worktrees/ directory
  const worktreesDir = '.worktrees';
  try {
    await fs.mkdir(worktreesDir, { recursive: false });
    console.log(`✓ Created ${worktreesDir}/`);
  } catch (error) {
    if (error.code === 'EEXIST') {
      console.log(`⊘ ${worktreesDir}/ already exists, skipping`);
    } else {
      throw error;
    }
  }

  // Add runtime directories to .gitignore
  const gitignorePath = '.gitignore';
  const gitignoreEntries = ['.worktrees/', '.dwp/logs/'];

  try {
    let gitignoreContent = '';
    try {
      gitignoreContent = await fs.readFile(gitignorePath, 'utf-8');
    } catch (error) {
      if (error.code !== 'ENOENT') throw error;
    }

    let newContent = gitignoreContent;
    for (const entry of gitignoreEntries) {
      if (newContent.split('\n').some(line => line.trim() === entry)) {
        console.log(`⊘ ${entry} already in .gitignore, skipping`);
        continue;
      }
      newContent = newContent
        ? `${newContent.trimEnd()}\n${entry}\n`
        : `${entry}\n`;
      console.log(`✓ Added ${entry} to .gitignore`);
    }
    await fs.writeFile(gitignorePath, newContent);
  } catch (error) {
    console.error(`Error updating .gitignore: ${error.message}`);
    process.exit(1);
  }

  console.log('\nDWP/AYNIG initialized successfully!');
}

export function registerInitCommand(program) {
  program
    .command('init')
    .description('Initialize DWP/AYNIG in the current repository')
    .action(action);
}
