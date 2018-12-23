'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';
// const cp = require('child_process');

const GoDefinitionProvider = require("./goto_definition");
const GoReferenceProvider = require("./findusage");
const GoWorkSpaceSymbolProvider = require("./dir_definitions");
const GoCompletionItemProvider = require("./auto_completion");
const GoDocumentSymbolProvider = require("./outline");
const GoHoverProvider = require("./hovers");

// let diagnosticCollection: vscode.DiagnosticCollection;



// this method is called when your extension is activated
// your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {
    
    let lucySelector : vscode.DocumentSelector = { scheme: 'file', language: 'lucy' } ; 
    
    // Use the console to output diagnostic information (console.log) and errors (console.error)
    // This line of code will only be executed once when your extension is activated
    console.log('Congratulations, your extension "lucy" is now active!');
    context.subscriptions.push(
        vscode.languages.registerDefinitionProvider(
            lucySelector, new GoDefinitionProvider()));
        context.subscriptions.push(
        vscode.languages.registerReferenceProvider(
            lucySelector, new GoReferenceProvider()));
        context.subscriptions.push(
            vscode.languages.registerWorkspaceSymbolProvider(
                new GoWorkSpaceSymbolProvider()));
        context.subscriptions.push(
            vscode.languages.registerCompletionItemProvider(
                lucySelector, new GoCompletionItemProvider(), '.' , '\"')); 
        context.subscriptions.push(
            vscode.languages.registerDocumentSymbolProvider(
                lucySelector, new GoDocumentSymbolProvider()));
        context.subscriptions.push(
            vscode.languages.registerHoverProvider(
                lucySelector, new GoHoverProvider()));
        // const collection = vscode.languages.createDiagnosticCollection('lucy');``
        // context.subscriptions.push(collection);
        // context.subscriptions.push(vscode.window.onDidChangeActiveTextEditor(e => updateDiagnostics(e, collection)));
}


function updateDiagnostics(document: vscode.TextEditor | undefined , collection: vscode.DiagnosticCollection): void {
	// if (document && path.basename(document.uri.fsPath) === 'sample-demo.rs') {
	// 	collection.set(document.uri, [{
	// 		code: '',
	// 		message: 'cannot assign twice to immutable variable `x`',
	// 		range: new vscode.Range(new vscode.Position(3, 4), new vscode.Position(3, 10)),
	// 		severity: vscode.DiagnosticSeverity.Error,
	// 		source: '',
	// 		relatedInformation: [
	// 			new vscode.DiagnosticRelatedInformation(new vscode.Location(document.uri, new vscode.Range(new vscode.Position(1, 8), new vscode.Position(1, 9))), 'first assignment to `x`')
	// 		]
	// 	}]);
	// } else {
	// 	collection.clear();
    // }
    // console.log("!!!!!!!!!!!!!!!!!!!!!!!!");
}




// this method is called when your extension is deactivated
export function deactivate() {

}






