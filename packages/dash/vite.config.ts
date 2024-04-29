import {defineConfig} from 'vite';
import react from '@vitejs/plugin-react';
import jotaiDebugLabel from 'jotai/babel/plugin-debug-label';
import jotaiReactRefresh from 'jotai/babel/plugin-react-refresh';
import {viteSingleFile} from 'vite-plugin-singlefile';

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [
		react({babel: {plugins: [jotaiDebugLabel, jotaiReactRefresh]}}),
		viteSingleFile(),
	],
	build: {
		sourcemap: false,
	},
});
