
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
        for(var name in definitions.variables) {
            let v = definitions.variables[name];
            let location = new vscode.Location(vscode.Uri.file(v.pos.filename) , new vscode.Position(v.pos.endLine , v.pos.columnOffset));
            let info = new vscode.SymbolInformation(name ,  vscode.SymbolKind.Variable  , "" , location );
            infos[i] = info;
            i++;
        }
        for(var name in definitions.constants) {
            let v = definitions.constants[name];
            let location = new vscode.Location(vscode.Uri.file(v.pos.filename) , new vscode.Position(v.pos.endLine , v.pos.columnOffset));
            let info = new vscode.SymbolInformation(name ,  vscode.SymbolKind.Constant  , "" , location );
            infos[i] = info;
            i++;
        }
        for(var name in definitions.functions) {
            let v = definitions.functions[name];
            let location = new vscode.Location(vscode.Uri.file(v.pos.filename) , new vscode.Position(v.pos.endLine , v.pos.columnOffset));
            let info = new vscode.SymbolInformation(name ,  vscode.SymbolKind.Function  , "" , location );
            infos[i] = info;
            i++;
        }
        for(var name in definitions.classes) {
            let v = definitions.classes[name];
            let location = new vscode.Location(vscode.Uri.file(v.pos.filename) , new vscode.Position(v.pos.endLine , v.pos.columnOffset));
            let info = new vscode.SymbolInformation(name ,  vscode.SymbolKind.Function  , "" , location );
            infos[i] = info;
            i++;
        }
        for(var name in definitions.enums) {
            let v = definitions.enums[name];
            let location = new vscode.Location(vscode.Uri.file(v.pos.filename) , new vscode.Position(v.pos.endLine , v.pos.columnOffset));
            let info = new vscode.SymbolInformation(name ,  vscode.SymbolKind.Function  , "" , location );
            infos[i] = info;
            i++;
        }
        return infos;
    }
};











