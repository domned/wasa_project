import { fileURLToPath, URL } from 'node:url';

import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

// https://vitejs.dev/config/
export default defineConfig(({ command, mode, ssrBuild }) => {
	const ret = {
		plugins: [vue()],
		resolve: {
			alias: {
				'@': fileURLToPath(new URL('./src', import.meta.url)),
			},
		},
	};
	ret.define = {
		// Use different API URLs for development and production
		__API_URL__: JSON.stringify(
			mode === 'production'
				? 'http://localhost:3000'
				: 'http://localhost:3000'
		),
	};
	ret.server = {
		proxy: {
			'/users': 'http://localhost:3000',
		},
	};
	return ret;
});
