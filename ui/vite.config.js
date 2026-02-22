import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [
		sveltekit()
	],
	server: {
		host: '0.0.0.0',
		port: 5173,
		allowedHosts: true,
		watch: {
			usePolling: true,
			interval: 100
		},
		hmr: {
			overlay: true
		}
	},
	preview: {
		host: '0.0.0.0',
		port: 5173
	},
});
