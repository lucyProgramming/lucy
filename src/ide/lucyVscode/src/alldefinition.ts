
'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';
const path = require('path');

const child_process = require('child_process');

module.exports = class GoDocumentSymbolProvider implements vscode.DocumentSymbolProvider {
    public provideDocumentSymbols(
        document: vscode.TextDocument, token: vscode.CancellationToken):
        Thenable<vscode.SymbolInformation[]> {
        let dir = path.dirname( document.fileName);
        let args = [
            "lucy.cmd.langtools.ide.alldefinition.main",
            "-dir",
            dir
        ];
        let result = child_process.execFileSync("java", args);
        console.log(result);
        let definitions = JSON.parse(result);
        if (!definitions) {
            return ;
        }
        //TODO:: 
        console.log(definitions);
        var i = 0 ;
        var infos = new Array();
        for(var pro in definitions) {
            let v = definitions[pro];
            let location = new vscode.Location(vscode.Uri.file(v.pos.filename) , new vscode.Position(v.pos.endLine , v.pos.endColumnOffset));
            console.log(v.name , location);
            var kind : vscode.SymbolKind ;
            switch(v.Type) {
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
            }
            let info = new vscode.SymbolInformation(v.name , kind  , "" , location);
            infos[i] = info;
            i++;
        }
        return infos;
    }
};











