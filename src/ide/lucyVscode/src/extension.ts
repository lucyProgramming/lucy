'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';

const GoDefinitionProvider = require("./goto_definition");
const GoReferenceProvider = require("./findusage");
const GoDocumentSymbolProvider = require("./alldefinition");

// this method is called when your extension is activated
// your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {
    
    let sel : vscode.DocumentSelector = { scheme: 'file', language: 'lucy' } ; 

    // Use the console to output diagnostic information (console.log) and errors (console.error)
    // This line of code will only be executed once when your extension is activated
    console.log('Congratulations, your extension "lucy" is now active!');
    context.subscriptions.push(
        vscode.languages.registerDefinitionProvider(
            sel, new GoDefinitionProvider()));
        context.subscriptions.push(
        vscode.languages.registerReferenceProvider(
            sel, new GoReferenceProvider()));
        context.subscriptions.push(
            vscode.languages.registerDocumentSymbolProvider(
                sel, new GoDocumentSymbolProvider()));
   
}


// this method is called when your extension is deactivated
export function deactivate() {
}

