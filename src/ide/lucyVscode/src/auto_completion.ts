'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';

const querystring = require('querystring');
const syncHttpRequest = require('sync-request');



module.exports = class GoCompletionItemProvider implements vscode.CompletionItemProvider {
    public provideCompletionItems(
        document: vscode.TextDocument, position: vscode.Position, token: vscode.CancellationToken):
        Thenable<vscode.CompletionItem[]> {
        var u = "http://localhost:2018/ide/autoCompletion?file=" + querystring.escape(document.fileName) + "&line=" + 
            position.line + "&column=" + position.character; 
        console.log(u);
        let buffer = document.getText();
        var res  = syncHttpRequest("POST" , u , {
            "body": buffer,
        });
        let lucyItems = JSON.parse(res.getBody());
        console.log(lucyItems);
        if(lucyItems.length === 0 ){
            console.log("auto completion length is 0");
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
                    v.name = v.suggest;
                    break;
                case "enumItem":
                    kind = vscode.CompletionItemKind.EnumMember;
                    break;
                default:
                    kind = vscode.CompletionItemKind.Text ;  
            }
            let item = new vscode.CompletionItem(v.name , kind);
            item.sortText = "" + i ; 
            items[i] = item;
        }
        return items;
    }
};