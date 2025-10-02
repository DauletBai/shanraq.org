import * as vscode from 'vscode';

export class TengeFormatter {
    private indentationSize: number = 4;
    private useSpaces: boolean = true;

    constructor() {
        // Get configuration
        const config = vscode.workspace.getConfiguration('tenge');
        this.indentationSize = config.get('formatter.indentationSize', 4);
        this.useSpaces = config.get('formatter.useSpaces', true);
    }

    public formatDocument(editor: vscode.TextEditor | undefined): void {
        if (!editor || editor.document.languageId !== 'tenge') {
            return;
        }

        const document = editor.document;
        const text = document.getText();
        const formattedText = this.formatText(text);

        if (formattedText !== text) {
            editor.edit(editBuilder => {
                const fullRange = new vscode.Range(
                    document.positionAt(0),
                    document.positionAt(text.length)
                );
                editBuilder.replace(fullRange, formattedText);
            });
        }
    }

    private formatText(text: string): string {
        const lines = text.split('\n');
        const formattedLines: string[] = [];
        let indentLevel = 0;
        let inComment = false;
        let inString = false;
        let stringChar = '';

        for (let i = 0; i < lines.length; i++) {
            let line = lines[i];
            const originalLine = line;
            
            // Skip empty lines
            if (line.trim() === '') {
                formattedLines.push('');
                continue;
            }

            // Handle comments
            if (line.trim().startsWith('//')) {
                formattedLines.push(this.indentLine(line.trim(), indentLevel));
                continue;
            }

            // Handle block comments
            if (line.includes('/*')) {
                inComment = true;
            }
            if (inComment) {
                formattedLines.push(this.indentLine(line.trim(), indentLevel));
                if (line.includes('*/')) {
                    inComment = false;
                }
                continue;
            }

            // Handle strings
            if (!inString) {
                if (line.includes('"') || line.includes("'") || line.includes('`')) {
                    inString = true;
                    stringChar = line.includes('"') ? '"' : (line.includes("'") ? "'" : '`');
                }
            } else {
                if (line.includes(stringChar)) {
                    inString = false;
                }
            }

            if (inString) {
                formattedLines.push(this.indentLine(line.trim(), indentLevel));
                continue;
            }

            // Remove existing indentation
            line = line.trim();

            // Handle closing braces
            if (line.startsWith('}') || line.startsWith('end') || line.startsWith('endif') || 
                line.startsWith('endwhile') || line.startsWith('endfor')) {
                indentLevel = Math.max(0, indentLevel - 1);
            }

            // Format the line
            const formattedLine = this.formatLine(line);
            formattedLines.push(this.indentLine(formattedLine, indentLevel));

            // Handle opening braces and control structures
            if (line.includes('{') || line.match(/^(atqar|jasau|eger|azirshe|while|for|if|else|function|class|struct)\b/)) {
                indentLevel++;
            }
        }

        return formattedLines.join('\n');
    }

    private formatLine(line: string): string {
        // Add spaces around operators
        line = line.replace(/([=+\-*/%<>!&|^])(?!=)/g, ' $1 ');
        line = line.replace(/([=+\-*/%<>!&|^])(?!=)/g, ' $1 ');
        
        // Remove multiple spaces
        line = line.replace(/\s+/g, ' ');
        
        // Add spaces after commas
        line = line.replace(/,/g, ', ');
        
        // Add spaces around parentheses
        line = line.replace(/\(/g, ' (');
        line = line.replace(/\)/g, ') ');
        
        // Remove multiple spaces again
        line = line.replace(/\s+/g, ' ');
        
        // Trim the line
        line = line.trim();
        
        // Ensure semicolon at end of statements
        if (line && !line.endsWith(';') && !line.endsWith('{') && !line.endsWith('}') && 
            !line.startsWith('//') && !line.startsWith('/*') && 
            !line.match(/^(atqar|jasau|eger|azirshe|qaytar|end|import|export)/)) {
            line += ';';
        }

        return line;
    }

    private indentLine(line: string, indentLevel: number): string {
        if (line.trim() === '') {
            return '';
        }

        const indent = this.useSpaces 
            ? ' '.repeat(indentLevel * this.indentationSize)
            : '\t'.repeat(indentLevel);
        
        return indent + line;
    }

    public formatSelection(editor: vscode.TextEditor, selection: vscode.Selection): void {
        const document = editor.document;
        const selectedText = document.getText(selection);
        const formattedText = this.formatText(selectedText);

        if (formattedText !== selectedText) {
            editor.edit(editBuilder => {
                editBuilder.replace(selection, formattedText);
            });
        }
    }

    public formatOnType(event: vscode.TextDocumentChangeEvent): void {
        // Auto-format on typing certain characters
        const change = event.contentChanges[0];
        if (change && change.text === ';') {
            // Auto-format when semicolon is typed
            const editor = vscode.window.activeTextEditor;
            if (editor && editor.document.languageId === 'tenge') {
                const config = vscode.workspace.getConfiguration('tenge');
                if (config.get('formatter.autoFormatOnType', false)) {
                    this.formatDocument(editor);
                }
            }
        }
    }
}
