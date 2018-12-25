

'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';

const path = require('path');
const querystring = require('querystring');
const request = require('request');


module.exports = class GoDefinitionProvider implements vscode.DefinitionProvider {
    public provideDefinition(
        document: vscode.TextDocument, position: vscode.Position, token: vscode.CancellationToken):
        Thenable<vscode.Location> {
            return new Promise<vscode.Location>(function(resolve, reject) {
            request({
                method : "POST",
                url: "http://localhost:2018" + "/ide/gotoDefinition?file=" + querystring.escape(document.fileName) + "&line=" + 
                    position.line + "&column=" + position.character,
                body : document.getText(),
            } , function(error : any, response : any, body:any) {
                if(error) {
                    console.log(error);
                    return ; 
                }
                let definition = JSON.parse(body);
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
        });
    }
};




