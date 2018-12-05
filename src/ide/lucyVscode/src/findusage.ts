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
        let holeText = document.getText();
        let bufferFileName =  document.fileName + ".buffer";
        fs.writeFileSync(bufferFileName);
        let args = [
            "lucy.cmd.langtools.ide.findusage.main",
            "-file",
            document.fileName,
            "-line",
            position.line + 1 ,
            "-column",
            position.character
        ]
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
        let result = child_process.execFileSync("java", args );
        let pos = JSON.parse(result);
        if (!pos) {
            return null;
        }
        let uri2 = new vscode.Uri("file" , "" , pos.filename , "" , "" , "");
        let position2 = new vscode.Position(pos.endLine - 1 , pos.columnOffset);
        fs.unlink(bufferFileName);
        return new vscode.Location(uri2, position2);
    }

};
