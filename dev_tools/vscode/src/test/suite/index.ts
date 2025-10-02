import * as path from 'path';
import * as Mocha from 'mocha';

export function run(): Promise<void> {
    // Create the mocha test
    const mocha = new (Mocha as any)({
        ui: 'bdd',
        color: true
    });

    const testsRoot = path.resolve(__dirname, '..');
    const testFile = path.resolve(testsRoot, 'extension.test.js');

    return new Promise((c, e) => {
        // Add test file
        mocha.addFile(testFile);

        try {
            // Run the mocha test
            mocha.run((failures: number) => {
                if (failures > 0) {
                    e(new Error(`${failures} tests failed.`));
                } else {
                    c();
                }
            });
        } catch (err) {
            console.error(err);
            e(err);
        }
    });
}
