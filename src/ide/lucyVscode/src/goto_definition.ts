'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';

const child_process = require('child_process');
const fs = require('fs');
const path = require('path');
const process = require('process');

module.exports = class GoDefinitionProvider implements vscode.DefinitionProvider {
    public provideDefinition(
        document: vscode.TextDocument, position: vscode.Position, token: vscode.CancellationToken):
        Thenable<vscode.Location> {
        // let holeText = document.getText();
        // let bufferFileName =  document.fileName + ".buffer";
        //fs.writeFileSync(bufferFileName);
        console.log(document.uri.toString());
        let args = [
            "lucy.cmd.langtools.ide.gotodefinition.main",
            "-file",
            document.fileName,
            "-line",
            position.line,
            "-column",
            position.character
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
        let result = child_process.execFileSync("java", args );
        let definition = JSON.parse(result);
        if (!definition) {
            return ;
        }
        let uri2 =  vscode.Uri.file(path.normalize(definition.filename)) ; 
        let position2 = new vscode.Position(definition.endLine, definition.endColumnOffset);
        //fs.unlink(bufferFileName);
        return new vscode.Location(uri2, position2);
        }
    };
};






