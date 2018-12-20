
'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';

const querystring = require('querystring');
const syncHttpRequest = require('sync-request');

module.exports = class GoReferenceProvider implements vscode.ReferenceProvider {
    public provideReferences(
        document: vscode.TextDocument, position: vscode.Position,
        options: { includeDeclaration: boolean }, token: vscode.CancellationToken):
        Thenable<vscode.Location[]> {
        return new Promise(function(resolve ,reject) {
            var u = "http://localhost:2018/ide/findUsage?file=" + querystring.escape(document.fileName) + "&line=" + 
            position.line + "&column=" + position.character; 
            console.log(u);
            var res  = syncHttpRequest("GET" , u);
            var usages = JSON.parse(res.getBody());
            if (!usages) {
                reject("not found usages");
                return;
            }
            console.log(usages);
            var items = new Array();
            for(var i = 0 ; i < usages.length ; i ++ ) {
                let v = usages[i];
                var uri2 = vscode.Uri.file(v.pos.filename);
                var position2 = new vscode.Position(v.pos.endLine,v.pos.endColumnOffset);
                var location =  new vscode.Location(uri2, position2);
                items[i] = location; 
            } 
            resolve(items);
        });
    }
};
