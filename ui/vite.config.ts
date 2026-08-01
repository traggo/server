/// <reference types="vitest/config" />
import {defineConfig} from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
    plugins: [react()],
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
