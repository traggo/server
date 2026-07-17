/// <reference types="vitest/config" />
import {defineConfig} from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
    // React is still on v16.12 until Stage 5's React 18 bump; v16.12 has no
    // automatic JSX runtime (react/jsx-runtime), so force the classic transform.
    plugins: [react({jsxRuntime: 'classic'})],
    // CRA/webpack polyfilled the Node `global` global for browser code; Vite doesn't.
    define: {
        global: 'globalThis',
    },
    base: './',
    server: {
        proxy: {
            '/graphql': 'http://localhost:3030',
        },
    },
    build: {
        outDir: 'build',
        sourcemap: false,
    },
    test: {
        environment: 'jsdom',
        globals: true,
        include: ['src/**/*.test.{ts,tsx}'],
    },
});
