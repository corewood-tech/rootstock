import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapter({
			pages: '.svelte-build',
			assets: '.svelte-build',
			fallback: 'index.html',
			strict: false
		}),
		paths: {
			base: '/app'
		}
	}
};

export default config;
