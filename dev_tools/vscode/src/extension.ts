import * as vscode from 'vscode';
import { TengeLinter } from './linter';
import { TengeFormatter } from './formatter';
import { TengeCompiler } from './compiler';

export function activate(context: vscode.ExtensionContext) {
    console.log('Tenge Language Support extension is now active!');

    // Initialize services
    const linter = new TengeLinter();
    const formatter = new TengeFormatter();
    const compiler = new TengeCompiler();

    // Register commands
    const formatCommand = vscode.commands.registerCommand('tenge.format', () => {
        const editor = vscode.window.activeTextEditor;
        if (editor && editor.document.languageId === 'tenge') {
            formatter.formatDocument(editor);
        }
    });

    const lintCommand = vscode.commands.registerCommand('tenge.lint', () => {
        const editor = vscode.window.activeTextEditor;
        if (editor && editor.document.languageId === 'tenge') {
            linter.lintDocument(editor);
        }
    });

    const compileCommand = vscode.commands.registerCommand('tenge.compile', () => {
        const editor = vscode.window.activeTextEditor;
        if (editor && editor.document.languageId === 'tenge') {
            compiler.compileDocument(editor);
        }
    });

    // Register providers
    const diagnosticCollection = vscode.languages.createDiagnosticCollection('tenge');
    context.subscriptions.push(diagnosticCollection);

    // Linter provider
    const linterProvider = vscode.languages.registerCodeActionsProvider('tenge', {
        provideCodeActions(document, range, context, token) {
            const actions: vscode.CodeAction[] = [];
            
            // Add quick fixes for common issues
            context.diagnostics.forEach(diagnostic => {
                if (diagnostic.code === 'tenge.agglutinative') {
                    const action = new vscode.CodeAction(
                        'Fix agglutinative naming',
                        vscode.CodeActionKind.QuickFix
                    );
                    action.diagnostics = [diagnostic];
                    actions.push(action);
                }
            });

            return actions;
        }
    });

    // Hover provider
    const hoverProvider = vscode.languages.registerHoverProvider('tenge', {
        provideHover(document, position, token) {
            const range = document.getWordRangeAtPosition(position);
            const word = document.getText(range);
            
            // Provide hover information for Tenge keywords and functions
            if (word.match(/^[a-zA-Z_][a-zA-Z0-9_]*$/)) {
                const hoverText = getHoverText(word);
                if (hoverText) {
                    return new vscode.Hover(hoverText);
                }
            }
            
            return null;
        }
    });

    // Completion provider
    const completionProvider = vscode.languages.registerCompletionItemProvider('tenge', {
        provideCompletionItems(document, position, token, context) {
            const completions: vscode.CompletionItem[] = [];
            
            // Add Tenge keywords
            const keywords = [
                'atqar', 'jasau', 'eger', 'azirshe', 'qaytar', 'end',
                'import', 'export', 'public', 'private', 'static',
                'san', 'jol', 'aqıqat', 'JsonObject', 'WebServer',
                'Component', 'TemplateEngine', 'ArchetypeEngine'
            ];
            
            keywords.forEach(keyword => {
                const item = new vscode.CompletionItem(keyword, vscode.CompletionItemKind.Keyword);
                completions.push(item);
            });
            
            // Add function patterns
            const functionPatterns = [
                '_jasau', '_alu', '_qosu', '_zhangartu', '_zhoyu', '_tekseru', '_opt'
            ];
            
            functionPatterns.forEach(pattern => {
                const item = new vscode.CompletionItem(pattern, vscode.CompletionItemKind.Snippet);
                item.insertText = pattern;
                completions.push(item);
            });
            
            return completions;
        }
    });

    // Register all providers
    context.subscriptions.push(
        formatCommand,
        lintCommand,
        compileCommand,
        linterProvider,
        hoverProvider,
        completionProvider
    );

    // Auto-format on save
    context.subscriptions.push(
        vscode.workspace.onDidSaveTextDocument(document => {
            if (document.languageId === 'tenge') {
                const config = vscode.workspace.getConfiguration('tenge');
                if (config.get('formatter.autoFormatOnSave', false)) {
                    formatter.formatDocument(vscode.window.activeTextEditor);
                }
            }
        })
    );

    // Auto-lint on change
    context.subscriptions.push(
        vscode.workspace.onDidChangeTextDocument(event => {
            if (event.document.languageId === 'tenge') {
                const config = vscode.workspace.getConfiguration('tenge');
                if (config.get('linter.enabled', true)) {
                    linter.lintDocument(vscode.window.activeTextEditor);
                }
            }
        })
    );
}

function getHoverText(word: string): string | null {
    const hoverInfo: { [key: string]: string } = {
        'atqar': 'Function declaration keyword in Tenge language',
        'jasau': 'Variable declaration keyword in Tenge language',
        'eger': 'If statement keyword in Tenge language',
        'azirshe': 'While loop keyword in Tenge language',
        'qaytar': 'Return statement keyword in Tenge language',
        'san': 'Integer number type in Tenge language',
        'jol': 'String type in Tenge language',
        'aqıqat': 'Boolean type in Tenge language',
        'JsonObject': 'JSON object type in Tenge language',
        'WebServer': 'Web server type in Tenge language',
        'Component': 'UI component type in Tenge language',
        'TemplateEngine': 'Template qozgaltqys type in Tenge language',
        'ArchetypeEngine': 'Archetype qozgaltqys type in Tenge language'
    };
    
    return hoverInfo[word] || null;
}

export function deactivate() {
    console.log('Tenge Language Support extension is now deactivated!');
}
