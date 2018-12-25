

'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';

const querystring = require('querystring');
const request = require('request');

//FIXME vscode don't goto the right place
module.exports = class GoHoverProvider implements vscode.HoverProvider {
    public provideHover(
        document: vscode.TextDocument, position: vscode.Position, token: vscode.CancellationToken):
        Thenable<vscode.Hover> {
        return new Promise(function(resolve , reject) {
            var u = "http://localhost:2018/ide/getHover?file=" + querystring.escape(document.fileName) + "&line=" + 
            position.line + "&column=" + position.character ; 
            console.log(u);
            request({
                method : "POST",
                url : u ,
                body : document.getText()
            } , function(error : any, response : any, body:any){
                if(error) {
                    console.log(error);
                    return ; 
                }
                var value = body;
                if(!value) {
                    reject("not found");
                    return ;
                }
                resolve(new vscode.Hover({language:"lucy" ,value : value }));
            });
        });
    }
};
