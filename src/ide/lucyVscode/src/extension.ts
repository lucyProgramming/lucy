'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';


// this method is called when your extension is activated
// your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {

    // Use the console to output diagnostic information (console.log) and errors (console.error)
    // This line of code will only be executed once when your extension is activated
    console.log('Congratulations, your extension "lucy" is now active!');
    context.subscriptions.push(
        vscode.languages.registerHoverProvider(
            "", new GoHoverProvider()));
   
}


// this method is called when your extension is deactivated
export function deactivate() {
}




class GoHoverProvider implements vscode.HoverProvider {
    constructor(){

    }

    public provideHover(
        document: vscode.TextDocument, position: vscode.Position, token: vscode.CancellationToken):
        Thenable<vscode.Hover> {
        console.log("call hover .... ,current now return nothing ");
        return new GoHoverProvider(); 
    }

    public then() : Thenable<GoHoverProvider> {
        return this ; 
    }
}

console.error("hello world");

[1][1] = 100; 
