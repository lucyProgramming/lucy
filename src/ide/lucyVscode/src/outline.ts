

'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';

const querystring = require('querystring');
const request = require('request');


//FIXME vscode don't goto the right place
module.exports = class GoDocumentSymbolProvider implements vscode.DocumentSymbolProvider {
    public provideDocumentSymbols(
        document: vscode.TextDocument, token: vscode.CancellationToken):
        Thenable<vscode.SymbolInformation[]> {
        return new Promise(function(resolve ,reject) {
            var u = "http://localhost:2018/ide/outline?filename=" + querystring.escape(document.fileName);
            console.log(u);
            request(u , function(error : any, response : any, body:any) {
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
                for(var i = 0 ;  i <  definitions.length ; i++ ) {
                    var v = definitions[i];
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
                        default:  
                            console.log(v.name , v.Type , "have not match use default");
                            kind = vscode.SymbolKind.Property;
                    }
                    var info = new vscode.DocumentSymbol(
                        v.name , 
                        "",
                        kind , 
                        new vscode.Range(
                            new vscode.Position(v.pos.startLine , v.pos.startColumnOffset) ,
                            new vscode.Position(v.pos.endLine , v.pos.endColumnOffset)
                            ),
                        new vscode.Range(
                            new vscode.Position(v.pos.startLine , v.pos.startColumnOffset) ,
                            new vscode.Position(v.pos.endLine , v.pos.endColumnOffset)
                            )
                        );
                    infos.push(info);
                    
                    if (v.inners) {
                        for(var j = 0 ;  j <  v.inners.length ; j++ ) {
                            var vv =  v.inners[j];
                            switch(vv.Type) {
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
                                case "method":
                                    kind = vscode.SymbolKind.Method;
                                    break;
                                case "constructor":
                                    kind = vscode.SymbolKind.Constructor;
                                    break;
                                case "field":
                                    kind = vscode.SymbolKind.Field;
                                    break;
                                case "enumItem":
                                    kind = vscode.SymbolKind.EnumMember;
                                    break;
                                default:  
                                    console.log(vv.name , vv.Type , "have not match use default");
                                    kind = vscode.SymbolKind.Property;
                            }
                            var info2 = new vscode.DocumentSymbol(
                                vv.name , 
                                "",
                                kind , 
                                new vscode.Range(
                                    new vscode.Position(vv.pos.startLine , vv.pos.startColumnOffset) ,
                                    new vscode.Position(vv.pos.endLine , vv.pos.endColumnOffset)
                                    ),
                                new vscode.Range(
                                    new vscode.Position(vv.pos.startLine , vv.pos.startColumnOffset) ,
                                    new vscode.Position(vv.pos.endLine , vv.pos.endColumnOffset)
                                    ));
                            
                            info.children.push(info2);
                        }
                    }
                }
                resolve(infos);
            });
        });
    }
};


