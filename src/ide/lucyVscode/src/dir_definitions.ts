

'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';

const querystring = require('querystring');
const request = require('request');
const path = require('path');


//FIXME vscode don't goto the right place
module.exports = class GoWorkSpaceSymbolProvider implements vscode.WorkspaceSymbolProvider {
    public provideWorkspaceSymbols(
        query: string, token: vscode.CancellationToken):
        Thenable<vscode.SymbolInformation[]> {
        return new Promise(function (resolve, reject) {
            if (!vscode.window.activeTextEditor) {
                reject("no active text editor");
                return;
            }
            var file = vscode.window.activeTextEditor.document.fileName;
            file = path.dirname(file);
            request("http://localhost:2018/ide/allDefinition?dir=" + querystring.escape(file), function(error : any, response : any, body:any) {
                if(error) {
                    console.log(error);
                    return ; 
                }
                var definitions = JSON.parse(body);
                if (!definitions) {
                    reject("not found");
                    return;
                }
                var infos = new Array();
                for (var i = 0; i < definitions.length; i++) {
                    var v = definitions[i];
                    var kind: vscode.SymbolKind;
                    switch (v.Type) {
                        case "variable":
                            kind = vscode.SymbolKind.Variable;
                            break;
                        case "function":
                            kind = vscode.SymbolKind.Function;
                            break;
                        case "constant":
                            kind = vscode.SymbolKind.Constant;
                            break;
                        case "class":
                            kind = vscode.SymbolKind.Class;
                            break;
                        case "enum":
                            kind = vscode.SymbolKind.Enum;
                            break;
                        case "typealias":
                            kind = vscode.SymbolKind.Operator;
                            break;
                        default:
                            console.log(v.name, v.Type, "have not match use default");
                            kind = vscode.SymbolKind.Property;
                    }
                    var info = new vscode.SymbolInformation(
                        v.name,
                        kind,
                        new vscode.Range(
                            new vscode.Position(v.pos.startLine, v.pos.startColumnOffset),
                            new vscode.Position(v.pos.endLine, v.pos.endColumnOffset)
                        ),
                        vscode.Uri.file(v.pos.filename),
                    );
                    infos.push(info);
                }
                resolve(infos);
            });
            
        });
}
};