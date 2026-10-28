import { fileURLToPath } from 'node:url';
import adapter from '@sveltejs/adapter-static';
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig, loadEnv } from 'vite';

const envDir = fileURLToPath(new URL('..', import.meta.url));

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, envDir, '');
	const serverPort = env.PORT || process.env.PORT || '8080';

	return {
		envDir,
		plugins: [
			sveltekit({
				compilerOptions: {
					runes: ({ filename }) =>
						filename.split(/[/\\]/).includes('node_modules') ? undefined : true,
				},
				adapter: adapter({ fallback: '200.html' }),
			}),
			tailwindcss(),
		],
		server: {
			proxy: {
				'/api/': `http://localhost:${serverPort}`,
			},
		},
	};
});
