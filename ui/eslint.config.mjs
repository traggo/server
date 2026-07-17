import js from '@eslint/js';
import tseslint from 'typescript-eslint';
import react from 'eslint-plugin-react';
import reactHooks from 'eslint-plugin-react-hooks';
import prettierConfig from 'eslint-config-prettier';
import globals from 'globals';

export default tseslint.config(
    {
        ignores: ['build/**', 'src/gql/__generated__/**', 'node_modules/**'],
    },
    js.configs.recommended,
    ...tseslint.configs.recommended,
    react.configs.flat.recommended,
    {
        // Only the two long-stable hook rules for now; eslint-plugin-react-hooks'
        // newer React-Compiler-oriented checks (immutability, refs, ...) are part
        // of the deliberate ruleset expansion + fixups planned for a later stage.
        plugins: {
            'react-hooks': reactHooks,
        },
        rules: {
            'react-hooks/rules-of-hooks': 'error',
            'react-hooks/exhaustive-deps': 'warn',
        },
    },
    prettierConfig,
    {
        languageOptions: {
            globals: {
                ...globals.browser,
                ...globals.node,
            },
        },
        settings: {
            react: {
                version: '16.12',
            },
        },
        rules: {
            'react/prop-types': 'off',
            'react/react-in-jsx-scope': 'off',
            'react/display-name': 'off',
            '@typescript-eslint/no-unused-vars': ['error', {argsIgnorePattern: '^_', ignoreRestSiblings: true}],
            'no-empty': ['error', {allowEmptyCatch: true}],
            '@typescript-eslint/no-explicit-any': 'off',
            '@typescript-eslint/no-empty-object-type': 'off',
        },
    }
);
