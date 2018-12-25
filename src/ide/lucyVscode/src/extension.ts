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

const querystring = require('querystring');
const request = require('request');



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
    const collection = vscode.languages.createDiagnosticCollection('lucy');
    context.subscriptions.push(collection);
    context.subscriptions.push(vscode.workspace.onDidSaveTextDocument(e => updateDiagnostics(e, collection , true)));
    context.subscriptions.push(vscode.workspace.onDidOpenTextDocument(e => updateDiagnostics(e, collection , true)));
    context.subscriptions.push(vscode.workspace.onDidChangeTextDocument(e => updateDiagnostics(e.document, collection , false)));
}


function updateDiagnostics(document: vscode.TextDocument ,collection: vscode.DiagnosticCollection , isSaveOrOpen : boolean): void {
    if (document.isUntitled) {
        return;  
    }
    if (document.languageId !== "lucy") {
        return;
    }
    collection.clear();
    var u = "http://localhost:2018/ide/diagnose?file=" + querystring.escape(document.fileName);
    console.log("!!!!!!!!!!!!!" , u);
    request({
        url: u ,
        method : "POST",
        body : document.getText(),
    } , function(error : any, response : any, body:any){
        if(error) {
            console.log(error);
            return ; 
        }
        var errs = JSON.parse(body);
        if(!errs) {
            return ; 
        }
        console.log(errs);
        for(let filename in errs) {
            console.log(filename , errs[filename]);
            var d = new Array();
            for(var i = 0 ; i < errs[filename].length ; i++) {
                var v = errs[filename][i];
                var t = new vscode.Diagnostic(
                    new vscode.Range(
                        new vscode.Position(v.pos.startLine , v.pos.startColumnOffset),
                        new vscode.Position(v.pos.startLine , v.pos.startColumnOffset)
                    ),
                    v.err
                );
                d.push(t);
            }
            collection.set(  vscode.Uri.file(filename ), d);
        }
    });
}




// this method is called when your extension is deactivated
export function deactivate() {

}






