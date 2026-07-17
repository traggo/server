import js from '@eslint/js';
import tseslint from 'typescript-eslint';
import react from 'eslint-plugin-react';
import reactHooks from 'eslint-plugin-react-hooks';
import sonarjs from 'eslint-plugin-sonarjs';
import prettierConfig from 'eslint-config-prettier';
import globals from 'globals';

export default tseslint.config(
    {
        ignores: ['build/**', 'src/gql/__generated__/**', 'node_modules/**'],
    },
    js.configs.recommended,
    ...tseslint.configs.recommended,
    react.configs.flat.recommended,
    sonarjs.configs.recommended,
    {
        // Only the two long-stable hook rules for now; eslint-plugin-react-hooks'
        // newer React-Compiler-oriented checks (immutability, refs, ...) are part
        // of a deliberate, separate ruleset expansion + code fixups, not bundled here.
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
                version: '18.3',
            },
        },
        rules: {
            'react/prop-types': 'off',
            'react/react-in-jsx-scope': 'off',
            'react/display-name': 'off',
            '@typescript-eslint/no-unused-vars': ['error', {argsIgnorePattern: '^_', ignoreRestSiblings: true}],
            'no-empty': ['error', {allowEmptyCatch: true}],
            '@typescript-eslint/no-empty-object-type': 'off',

            // Mapped from the old tslint.json + tslint-sonarts ruleset (see the tslint.json
            // that existed at this path before the tslint->ESLint migration, in git history).
            '@typescript-eslint/no-explicit-any': 'error',
            '@typescript-eslint/explicit-member-accessibility': 'error',
            '@typescript-eslint/naming-convention': [
                'error',
                {
                    selector: 'variableLike',
                    format: ['camelCase', 'PascalCase', 'UPPER_CASE'],
                    leadingUnderscore: 'allow',
                    trailingUnderscore: 'allow',
                },
            ],
            eqeqeq: ['error', 'always', {null: 'ignore'}],
            'no-var': 'error',
            curly: 'error',
            'guard-for-in': 'error',
            'default-case': 'error',
            'no-shadow': 'off',
            '@typescript-eslint/no-shadow': 'error',
            'object-shorthand': 'error',
            'no-duplicate-imports': 'error',
            radix: 'error',
            'use-isnan': 'error',
            'no-debugger': 'error',
            'no-cond-assign': 'error',
            '@typescript-eslint/no-inferrable-types': 'error',
            '@typescript-eslint/no-namespace': ['error', {allowDeclarations: true}],
            'no-redeclare': 'off',
            '@typescript-eslint/no-redeclare': 'error',
            'no-throw-literal': 'error',
            'no-unused-expressions': 'off',
            '@typescript-eslint/no-unused-expressions': 'error',
            complexity: ['error', 20],

            // sonarjs's recommended preset is tuned for large codebases and disables a few
            // rules by default (e.g. cognitive-complexity) that tslint-sonarts had enabled.
            'sonarjs/cognitive-complexity': ['error', 20],
        },
    }
);
