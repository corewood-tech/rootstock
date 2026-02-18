import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
	plugins: [
		tailwindcss(),
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
