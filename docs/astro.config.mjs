import { defineConfig } from 'astro/config';
import sitemap from '@astrojs/sitemap';
import starlight from '@astrojs/starlight';

export default defineConfig({
  site: 'https://aynig.org',
  integrations: [
    starlight({
      title: 'All You Need Is Git',
      description: 'AYNIG runner and CLI documentation.',
      head: [],
      customCss: ['./src/styles/custom.css'],
      sidebar: [
        {
          label: 'Getting Started',
          items: [
            { label: 'Overview', link: '/' },
            { label: 'Quick Start', link: '/getting-started/quick-start/' }
          ]
        },
        {
          label: 'Installation',
          items: [
            { label: 'Install CLI', link: '/installation/cli/' }
          ]
        },
        {
          label: 'Repository Setup',
          items: [
            { label: 'Initialize', link: '/repository/init/' },
            { label: 'Repository Layout', link: '/repository/layout/' },
            { label: 'Install Workflow Packs', link: '/repository/install-workflows/' }
          ]
        },
        {
          label: 'Run and Operate',
          items: [
            { label: 'Run AYNIG', link: '/run/run/' },
            { label: 'Worktrees and Isolation', link: '/run/worktrees/' },
            { label: 'Leases and Liveness', link: '/run/leases/' }
          ]
        },
        {
          label: 'Operate AYNIG',
          items: [
            { label: 'Runbooks', link: '/operate/runbooks/' },
            { label: 'Retries', link: '/operate/retries/' },
            { label: 'History Hygiene', link: '/operate/history-hygiene/' },
            { label: 'Concurrency Patterns', link: '/operate/concurrency-patterns/' },
            { label: 'Trailer Conventions', link: '/operate/trailer-conventions/' }
          ]
        },
        {
          label: 'Commands',
          items: [
            { label: 'Authoring Commands', link: '/commands/authoring/' },
            { label: 'Commit Protocol', link: '/commands/protocol/' },
            { label: 'Environment Variables', link: '/commands/environment/' }
          ]
        },
        {
          label: 'Workflow Design',
          items: [
            { label: 'Workflow Design', link: '/workflows/design/' }
          ]
        },
        {
          label: 'Why AYNIG',
          items: [
            { label: 'Why AYNIG', link: '/why/' },
            { label: 'Benefits', link: '/why/benefits/' },
            { label: 'Tradeoffs', link: '/why/tradeoffs/' }
          ]
        },
        {
          label: 'CLI Reference',
          items: [
            { label: 'init', link: '/cli/init/' },
            { label: 'install', link: '/cli/install/' },
            { label: 'run', link: '/cli/run/' },
            { label: 'set-working', link: '/cli/set-working/' },
            { label: 'set-state', link: '/cli/set-state/' }
          ]
        },
        {
          label: 'Kernel Contract',
          items: [
            { label: 'Kernel Contract', link: '/contract/' }
          ]
        }
      ]
    }),
    sitemap()
  ]
});
