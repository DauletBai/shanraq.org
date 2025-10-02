import * as assert from 'assert';
import * as vscode from 'vscode';

describe('Tenge Extension Test Suite', () => {
    it('Extension should be present', () => {
        assert.ok(vscode.extensions.getExtension('shanraq.tenge-language-support'));
    });

    it('Extension should activate', async () => {
        const extension = vscode.extensions.getExtension('shanraq.tenge-language-support');
        if (extension) {
            await extension.activate();
            assert.ok(extension.isActive);
        }
    });

    it('Tenge language should be registered', async () => {
        const languages = await vscode.languages.getLanguages();
        assert.ok(languages.includes('tenge'));
    });

    it('Tenge grammar should be registered', () => {
        const grammar = vscode.workspace.getConfiguration('tenge');
        assert.ok(grammar);
    });
});
