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

let diagnosticCollection: vscode.DiagnosticCollection;


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
        
        diagnosticCollection = vscode.languages.createDiagnosticCollection('go');
        context.subscriptions.push(diagnosticCollection);
}


// function onChange() {
//     let uri = document.uri;
//     check(uri.fsPath, goConfig).then(errors => {
//       diagnosticCollection.clear();
//       let diagnosticMap: Map<string, vscode.Diagnostic[]> = new Map();
//       errors.forEach(error => {
//         let canonicalFile = vscode.Uri.file(error.file).toString();
//         let range = new vscode.Range(error.line-1, error.startColumn, error.line-1, error.endColumn);
//         let diagnostics = diagnosticMap.get(canonicalFile);
//         if (!diagnostics) { diagnostics = []; }
//         diagnostics.push(new vscode.Diagnostic(range, error.msg, error.severity));
//         diagnosticMap.set(canonicalFile, diagnostics);
//       });
//       diagnosticMap.forEach((diags, file) => {
//         diagnosticCollection.set(vscode.Uri.parse(file), diags);
//       });
//     })
// }



// this method is called when your extension is deactivated
export function deactivate() {

}






