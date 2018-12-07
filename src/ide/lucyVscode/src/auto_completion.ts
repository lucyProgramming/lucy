'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';
const fs = require('fs');
const child_process = require('child_process');

const bufferFile = "./ffsdfw3er2233242342wewe4233423.buffer";


module.exports = class GoCompletionItemProvider implements vscode.CompletionItemProvider {
    public provideCompletionItems(
        document: vscode.TextDocument, position: vscode.Position, token: vscode.CancellationToken):
        Thenable<vscode.CompletionItem[]> {
        console.log(process.cwd());
        fs.writeFileSync(bufferFile , document.getText());
        let args = [
            "lucy.cmd.langtools.ide.autocompletion.main",
            "-file",
            document.fileName,
            "-line",
            position.line,
            "-column",
            position.character,
            "-bufferFile",
            bufferFile
        ];
        let result = child_process.execFileSync("java", args);
        console.log(result.toString());
        let lucyItems = JSON.parse(result);
        if (!lucyItems) {
            return null;
        }
        var items = new Array();
        for (var i = 0 ; i < lucyItems.length ; i++) {
            var kind : vscode.CompletionItemKind = vscode.CompletionItemKind.Text;
            var v = lucyItems[i] ; 
            switch (v.Type) {
                case "constant":
                    kind = vscode.CompletionItemKind.Constant;
                    break;
                case "variable":
                    kind = vscode.CompletionItemKind.Variable;
                    break;
                case "function":
                    v.name = v.suggest;
                    kind = vscode.CompletionItemKind.Function;
                    break;
                case "class":
                    kind = vscode.CompletionItemKind.Class;
                    break;
                case "field":
                    kind = vscode.CompletionItemKind.Field;
                    break;
                case "method":
                    kind = vscode.CompletionItemKind.Method;
                    break;
                case "enumItem":
                    kind = vscode.CompletionItemKind.EnumMember;
                    break;
                default:
                    kind = vscode.CompletionItemKind.Text ;  
            }
            let item = new vscode.CompletionItem(v.name , kind);
            items[i] = item;
        }
        return items;
    }
};