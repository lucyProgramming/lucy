

'use strict';
// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';

const path = require('path');
const querystring = require('querystring');
const syncHttpRequest = require('sync-request');

//FIXME vscode don't goto the right place
module.exports = class GoDocumentSymbolProvider implements vscode.DocumentSymbolProvider {
    public provideDocumentSymbols(
        document: vscode.TextDocument, token: vscode.CancellationToken):
        Thenable<vscode.SymbolInformation[]> {
        return new Promise(function(resolve ,reject) {
            let dir = path.dirname(document.fileName);
            var u = "http://localhost:2018/ide/allDefinition?dir=" + querystring.escape(dir);
            console.log(u);
            var res  = syncHttpRequest("GET" , u);
            var definitions = JSON.parse(res.getBody());
            if (!definitions) {
                reject("not found");
                return;
            }
            //TODO:: 
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
                        kind = vscode.SymbolKind.Property;
                }
                var location2 = new vscode.Location(vscode.Uri.file(v.pos.filename) , new vscode.Range(new vscode.Position(v.pos.endLine , v.pos.endColumnOffset) , new vscode.Position(v.pos.endLine , v.pos.endColumnOffset))   );
                console.log( v.name , location2.uri.toString());
                var info = new vscode.SymbolInformation(v.name , kind  , "" , location2);
                infos.push(info);
            }
            resolve(infos);
        });
    }
};













// 'use strict';
// // The module 'vscode' contains the VS Code extensibility API
// // Import the module and reference it with the alias vscode in your code below
// import * as vscode from 'vscode';

// const path = require('path');
// const querystring = require('querystring');
// const syncHttpRequest = require('sync-request');

// //FIXME vscode don't goto the right place
// module.exports = class GoDocumentSymbolProvider implements vscode.DocumentSymbolProvider {
//     public provideDocumentSymbols(
//         document: vscode.TextDocument, token: vscode.CancellationToken):
//         Thenable<vscode.SymbolInformation[]> {
//         let dir = path.dirname(document.fileName);
//         var u = "http://localhost:2018/ide/allDefinition?dir=" + querystring.escape(dir);
//         console.log(u);
//         var res  = syncHttpRequest("GET" , u);
//         var definitions = JSON.parse(res.getBody());
//         if (!definitions) {
//             return null;
//         }
//         //TODO:: 
//         var infos = new Array();
//         for(var i = 0 ;  i <  definitions.length ; i++ ) {
//             var v = definitions[i];
//             var kind : vscode.SymbolKind ;
//             switch(v.Type) {
//                 case "variable":
//                     kind = vscode.SymbolKind.Variable;
//                     break;
//                 case "function":
//                     kind = vscode.SymbolKind.Function;
//                     break;
//                 case "constant":
//                     kind = vscode.SymbolKind.Constant;
//                     break;
//                 case "class":
//                     kind = vscode.SymbolKind.Class;
//                     break;
//                 case "enum":
//                     kind = vscode.SymbolKind.Enum;
//                     break;
//             }
//             var uri2 = vscode.Uri.file(v.pos.filename);
//             var position2 = new vscode.Position(v.pos.endLine,v.pos.endColumnOffset);
//             var location2 = new vscode.Location(uri2 , position2);
//             var info = new vscode.SymbolInformation(v.name , kind  , "" , location2);
//             infos[i] = info;
//         }
//         return infos;
//     }
// };











