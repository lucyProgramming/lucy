
'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';

const path = require('path');
const querystring = require('querystring');
const syncHttpRequest = require('sync-request');


module.exports = class GoDefinitionProvider implements vscode.DefinitionProvider {
    public provideDefinition(
        document: vscode.TextDocument, position: vscode.Position, token: vscode.CancellationToken):
        Thenable<vscode.Location> {
            return new Promise<vscode.Location>(function(resolve, reject) {
            var u = "http://localhost:2018/ide/gotoDefinition?file=" + querystring.escape(document.fileName) + "&line=" + 
            position.line + "&column=" + position.character ; 
            console.log(u);
            let buffer = document.getText();
            var res  = syncHttpRequest("POST" , u ,  {
                "body": buffer,
                "timeout" : 2000
            });
            let definition = JSON.parse(res.getBody());
            if (!definition) {
                reject("definition not found");
                return ;
            }
            let uri2 =  vscode.Uri.file(path.normalize(definition.filename)) ; 
            // var r = new vscode.Range(
            //         new vscode.Position(definition.satrtLine, definition.startColumnOffset),
            //         new vscode.Position(definition.endLine, definition.endColumnOffset)
            // );
            var r = new vscode.Position( definition.endLine, definition.endColumnOffset);
            resolve(new vscode.Location(uri2, r));
        });
    }
};



