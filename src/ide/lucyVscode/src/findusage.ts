'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';
const child_process = require('child_process');
const fs = require('fs');


module.exports = class GoReferenceProvider implements vscode.ReferenceProvider {
    public provideReferences(
        document: vscode.TextDocument, position: vscode.Position,
        options: { includeDeclaration: boolean }, token: vscode.CancellationToken):
        Thenable<vscode.Location[]> {
        var args = [
            "lucy.cmd.langtools.ide.gotodefinition.main",
            "-file",
            document.fileName,
            "-line",
            position.line,
            "-column",
            position.character
        ];
        var result = child_process.execFileSync("java", args );
        var pos = JSON.parse(result);
        if (!pos) {
            console.log("definition not found");
            return null ;
        }
        args = [
            "lucy.cmd.langtools.ide.findusage.main",
            "-file",
            pos.filename,
            "-line",
            pos.endLine ,
            "-column",
            pos.columnOffset -1 
        ];
        {
            let s =  "java ";
            for (var i = 0 ; i < args.length ; i++) {
                s += args[i] ; 
                if(i !== args.length - 1 ) {
                    s += " ";
                }
            }
            console.log(s);
        }
        result = child_process.execFileSync("java", args);
        pos = JSON.parse(result);
        if (!pos) {
            return null;
        }
        console.log(pos);
        var items = new Array();
        for(var i = 0 ; i < pos.length ; i ++ ) {
            let v = pos[i];
            var uri2 = vscode.Uri.file(v.pos.filename);
            var position2 = new vscode.Position(v.pos.endLine,v.pos.columnOffset);
            var location =  new vscode.Location(uri2, position2);
            items[i] = location; 
        } 
        return items; 
    }

};
