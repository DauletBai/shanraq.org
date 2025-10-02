import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';

export class TengeCompiler {
    private outputChannel: vscode.OutputChannel;

    constructor() {
        this.outputChannel = vscode.window.createOutputChannel('Tenge Compiler');
    }

    public async compileDocument(editor: vscode.TextEditor | undefined): Promise<void> {
        if (!editor || editor.document.languageId !== 'tenge') {
            return;
        }

        const document = editor.document;
        const filePath = document.uri.fsPath;
        const fileName = path.basename(filePath, '.tng');
        const outputDir = path.dirname(filePath);
        const outputPath = path.join(outputDir, `${fileName}.c`);

        this.outputChannel.clear();
        this.outputChannel.show();
        this.outputChannel.appendLine('🔨 Compiling Tenge to C...');

        try {
            // Get Tenge compiler path from configuration
            const config = vscode.workspace.getConfiguration('tenge');
            const compilerPath = config.get('compiler.path', 'tenge');

            // Check if compiler exists
            if (!await this.checkCompilerExists(compilerPath)) {
                this.outputChannel.appendLine(`❌ Tenge compiler not found at: ${compilerPath}`);
                this.outputChannel.appendLine('Please install Tenge compiler or update the path in settings.');
                return;
            }

            // Compile the file
            const result = await this.runCompiler(compilerPath, filePath, outputPath);
            
            if (result.success) {
                this.outputChannel.appendLine('✅ Compilation successful!');
                this.outputChannel.appendLine(`📁 Output file: ${outputPath}`);
                
                // Show success message
                vscode.window.showInformationMessage(
                    `Tenge file compiled successfully! Output: ${path.basename(outputPath)}`
                );

                // Open the generated C file
                const cFileUri = vscode.Uri.file(outputPath);
                await vscode.window.showTextDocument(cFileUri);
            } else {
                this.outputChannel.appendLine('❌ Compilation failed!');
                this.outputChannel.appendLine(result.error || 'Unknown error');
                
                // Show error message
                vscode.window.showErrorMessage('Tenge compilation failed! Check the output for details.');
            }
        } catch (error) {
            this.outputChannel.appendLine(`❌ Compilation error: ${error}`);
            vscode.window.showErrorMessage(`Tenge compilation error: ${error}`);
        }
    }

    private async checkCompilerExists(compilerPath: string): Promise<boolean> {
        try {
            // Try to run the compiler with --version flag
            const { exec } = require('child_process');
            const { promisify } = require('util');
            const execAsync = promisify(exec);
            
            await execAsync(`${compilerPath} --version`);
            return true;
        } catch (error) {
            return false;
        }
    }

    private async runCompiler(compilerPath: string, inputPath: string, outputPath: string): Promise<{success: boolean, error?: string}> {
        return new Promise((resolve) => {
            const { exec } = require('child_process');
            
            const command = `${compilerPath} -c "${inputPath}" -o "${outputPath}"`;
            this.outputChannel.appendLine(`🔧 Running: ${command}`);

            exec(command, (error: any, stdout: string, stderr: string) => {
                if (error) {
                    this.outputChannel.appendLine(`❌ Compiler error: ${error.message}`);
                    this.outputChannel.appendLine(`stderr: ${stderr}`);
                    resolve({ success: false, error: error.message });
                } else {
                    this.outputChannel.appendLine(`stdout: ${stdout}`);
                    if (stderr) {
                        this.outputChannel.appendLine(`stderr: ${stderr}`);
                    }
                    resolve({ success: true });
                }
            });
        });
    }

    public async buildProject(): Promise<void> {
        const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
        if (!workspaceFolder) {
            vscode.window.showErrorMessage('No workspace folder found');
            return;
        }

        this.outputChannel.clear();
        this.outputChannel.show();
        this.outputChannel.appendLine('🔨 Building Tenge project...');

        try {
            // Find all .tng files in the workspace
            const tngFiles = await vscode.workspace.findFiles('**/*.tng');
            
            if (tngFiles.length === 0) {
                this.outputChannel.appendLine('❌ No Tenge files found in workspace');
                return;
            }

            this.outputChannel.appendLine(`📁 Found ${tngFiles.length} Tenge files`);

            // Compile each file
            let successCount = 0;
            let errorCount = 0;

            for (const file of tngFiles) {
                const fileName = path.basename(file.fsPath, '.tng');
                const outputPath = path.join(path.dirname(file.fsPath), `${fileName}.c`);
                
                this.outputChannel.appendLine(`🔧 Compiling ${path.basename(file.fsPath)}...`);
                
                const result = await this.runCompiler('tenge', file.fsPath, outputPath);
                
                if (result.success) {
                    successCount++;
                    this.outputChannel.appendLine(`✅ ${path.basename(file.fsPath)} compiled successfully`);
                } else {
                    errorCount++;
                    this.outputChannel.appendLine(`❌ ${path.basename(file.fsPath)} compilation failed: ${result.error}`);
                }
            }

            this.outputChannel.appendLine(`\n📊 Build Summary:`);
            this.outputChannel.appendLine(`✅ Successful: ${successCount}`);
            this.outputChannel.appendLine(`❌ Failed: ${errorCount}`);

            if (errorCount === 0) {
                vscode.window.showInformationMessage(`Tenge project built successfully! ${successCount} files compiled.`);
            } else {
                vscode.window.showWarningMessage(`Tenge project build completed with ${errorCount} errors.`);
            }
        } catch (error) {
            this.outputChannel.appendLine(`❌ Build error: ${error}`);
            vscode.window.showErrorMessage(`Tenge project build error: ${error}`);
        }
    }

    public async runProject(): Promise<void> {
        const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
        if (!workspaceFolder) {
            vscode.window.showErrorMessage('No workspace folder found');
            return;
        }

        this.outputChannel.clear();
        this.outputChannel.show();
        this.outputChannel.appendLine('🚀 Running Tenge project...');

        try {
            // Look for main.tng or index.tng
            const mainFiles = await vscode.workspace.findFiles('**/main.tng');
            const indexFiles = await vscode.workspace.findFiles('**/index.tng');
            
            const mainFile = mainFiles[0] || indexFiles[0];
            
            if (!mainFile) {
                this.outputChannel.appendLine('❌ No main Tenge file found (main.tng or index.tng)');
                return;
            }

            this.outputChannel.appendLine(`📁 Running: ${path.basename(mainFile.fsPath)}`);

            // Run the Tenge file
            const { exec } = require('child_process');
            
            exec(`tenge run "${mainFile.fsPath}"`, (error: any, stdout: string, stderr: string) => {
                if (error) {
                    this.outputChannel.appendLine(`❌ Runtime error: ${error.message}`);
                    this.outputChannel.appendLine(`stderr: ${stderr}`);
                } else {
                    this.outputChannel.appendLine(`stdout: ${stdout}`);
                    if (stderr) {
                        this.outputChannel.appendLine(`stderr: ${stderr}`);
                    }
                }
            });
        } catch (error) {
            this.outputChannel.appendLine(`❌ Runtime error: ${error}`);
            vscode.window.showErrorMessage(`Tenge runtime error: ${error}`);
        }
    }
}
