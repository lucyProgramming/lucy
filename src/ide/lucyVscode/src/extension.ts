'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';


const GoDefinitionProvider = require("./goto_definition");
const GoReferenceProvider = require("./findusage");
// const GoDocumentSymbolProvider = require("./alldefinition");
const GoCompletionItemProvider = require("./auto_completion");


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
        // context.subscriptions.push(
        //     vscode.languages.registerDocumentSymbolProvider(
        //         lucySelector, new GoDocumentSymbolProvider()));
        // context.subscriptions.push(
        //     vscode.languages.registerCompletionItemProvider(
        //         lucySelector, new GoCompletionItemProvider(), 
        //         //'+' , '-' , '*' , '/' , '%', '=' , '&' , '|' ,
        //         //'>' , '<' ,'!' , '^' , '~',
        //         '.','q','w','e','r','t','y','u','i','o','p' ,
        //         'a','s','d','f','g','h','j','k','l',
        //         'z','x','c','v','b','n','m' ,
        //        'Q','W','E','R','T','Y','U','I','O','P',
        //        'A','S','D','F','G','H','J','K','L',
        //        'Z','X','C','V','B','N','M' ,
        //        '$' , '_' 
        //     ));
        context.subscriptions.push(
            vscode.languages.registerCompletionItemProvider(
                lucySelector, new GoCompletionItemProvider(), '.' , '\"'));
        
   
}


// this method is called when your extension is deactivated
export function deactivate() {
}

