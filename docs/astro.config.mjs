import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
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
            { label: 'Install', link: '/getting-started/install/' },
            { label: 'Run', link: '/getting-started/run/' }
          ]
        },
        {
          label: 'CLI',
          items: [
            { label: 'init', link: '/cli/init/' },
            { label: 'install', link: '/cli/install/' },
            { label: 'run', link: '/cli/run/' }
          ]
        }
      ]
    })
  ]
});
