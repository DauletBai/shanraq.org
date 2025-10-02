import * as vscode from 'vscode';

export class TengeLinter {
    private diagnosticCollection: vscode.DiagnosticCollection;

    constructor() {
        this.diagnosticCollection = vscode.languages.createDiagnosticCollection('tenge');
    }

    public lintDocument(editor: vscode.TextEditor | undefined): void {
        if (!editor || editor.document.languageId !== 'tenge') {
            return;
        }

        const document = editor.document;
        const diagnostics: vscode.Diagnostic[] = [];
        const text = document.getText();

        // Parse the document and find issues
        const lines = text.split('\n');
        
        lines.forEach((line, lineNumber) => {
            const lineDiagnostics = this.lintLine(line, lineNumber);
            diagnostics.push(...lineDiagnostics);
        });

        // Set diagnostics
        this.diagnosticCollection.set(document.uri, diagnostics);
    }

    private lintLine(line: string, lineNumber: number): vscode.Diagnostic[] {
        const diagnostics: vscode.Diagnostic[] = [];
        const trimmedLine = line.trim();

        // Check for agglutinative naming patterns
        const functionMatch = line.match(/\b[a-zA-Z_][a-zA-Z0-9_]*\s*\(/g);
        if (functionMatch) {
            functionMatch.forEach(match => {
                const functionName = match.replace(/\s*\(/, '');
                if (!this.isAgglutinativeName(functionName)) {
                    const diagnostic = new vscode.Diagnostic(
                        new vscode.Range(lineNumber, line.indexOf(functionName), lineNumber, line.indexOf(functionName) + functionName.length),
                        `Function name "${functionName}" should follow agglutinative patterns (e.g., use _jasau, _alu, _qosu suffixes)`,
                        vscode.DiagnosticSeverity.Warning
                    );
                    diagnostic.code = 'tenge.agglutinative';
                    diagnostics.push(diagnostic);
                }
            });
        }

        // Check for proper Tenge keywords
        if (trimmedLine.startsWith('function ') || trimmedLine.startsWith('var ') || trimmedLine.startsWith('let ')) {
            const diagnostic = new vscode.Diagnostic(
                new vscode.Range(lineNumber, 0, lineNumber, trimmedLine.length),
                'Use Tenge keywords: "atqar" instead of "function", "jasau" instead of "var/let"',
                vscode.DiagnosticSeverity.Error
            );
            diagnostic.code = 'tenge.keywords';
            diagnostics.push(diagnostic);
        }

        // Check for missing semicolons
        if (trimmedLine && !trimmedLine.startsWith('//') && !trimmedLine.startsWith('/*') && 
            !trimmedLine.endsWith(';') && !trimmedLine.endsWith('{') && !trimmedLine.endsWith('}') &&
            !trimmedLine.match(/^(atqar|jasau|eger|azirshe|qaytar|end|import|export)/)) {
            const diagnostic = new vscode.Diagnostic(
                new vscode.Range(lineNumber, line.length - 1, lineNumber, line.length),
                'Missing semicolon at end of statement',
                vscode.DiagnosticSeverity.Warning
            );
            diagnostic.code = 'tenge.semicolon';
            diagnostics.push(diagnostic);
        }

        // Check for proper type annotations
        const variableMatch = line.match(/jasau\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*:/);
        if (variableMatch) {
            const variableName = variableMatch[1];
            if (!line.includes(': san') && !line.includes(': jol') && !line.includes(': aqıqat') && 
                !line.includes(': JsonObject') && !line.includes(': WebServer') && 
                !line.includes(': Component') && !line.includes(': TemplateEngine')) {
                const diagnostic = new vscode.Diagnostic(
                    new vscode.Range(lineNumber, line.indexOf(variableName), lineNumber, line.indexOf(variableName) + variableName.length),
                    `Variable "${variableName}" should have a proper type annotation`,
                    vscode.DiagnosticSeverity.Warning
                );
                diagnostic.code = 'tenge.type';
                diagnostics.push(diagnostic);
            }
        }

        // Check for proper return statements
        if (line.includes('atqar ') && !line.includes('qaytar ')) {
            const diagnostic = new vscode.Diagnostic(
                new vscode.Range(lineNumber, 0, lineNumber, line.length),
                'Function should have a return statement',
                vscode.DiagnosticSeverity.Warning
            );
            diagnostic.code = 'tenge.return';
            diagnostics.push(diagnostic);
        }

        // Check for Kazakh language compliance
        const englishWords = line.match(/\b(create|get|add|update|delete|check|optimize|qozgaltqys|manager)\b/gi);
        if (englishWords) {
            englishWords.forEach(word => {
                const diagnostic = new vscode.Diagnostic(
                    new vscode.Range(lineNumber, line.indexOf(word), lineNumber, line.indexOf(word) + word.length),
                    `Use Kazakh equivalents: "${word}" should be replaced with agglutinative Kazakh morphemes`,
                    vscode.DiagnosticSeverity.Warning
                );
                diagnostic.code = 'tenge.kazakh';
                diagnostics.push(diagnostic);
            });
        }

        return diagnostics;
    }

    private isAgglutinativeName(name: string): boolean {
        // Check if the name follows agglutinative patterns
        const agglutinativeSuffixes = [
            '_jasau', '_alu', '_qosu', '_zhangartu', '_zhoyu', '_tekseru', '_opt',
            '_eng', '_man', '_negizgi', '_qoldanu', '_marshrut', '_qosu', '_baskaru'
        ];
        
        return agglutinativeSuffixes.some(suffix => name.endsWith(suffix));
    }

    public clearDiagnostics(): void {
        this.diagnosticCollection.clear();
    }
}
