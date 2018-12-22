

// 'use strict';
// // The module 'vscode' contains the VS Code extensibility API
// // Import the module and reference it with the alias vscode in your code below
// import * as vscode from 'vscode';


// module.exports = function(definitions : any ): vscode.DocumentSymbol[] {
//     var infos = new Array();
//     for(var i = 0 ;  i <  definitions.length ; i++ ) {
//         var v = definitions[i];
//         var kind : vscode.SymbolKind ;
//         switch(v.Type) {
//             case "variable":
//                 kind = vscode.SymbolKind.Variable;
//                 break;
//             case "function":
//                 kind = vscode.SymbolKind.Function;
//                 break;
//             case "constant":
//                 kind = vscode.SymbolKind.Constant;
//                 break;
//             case "class":
//                 kind = vscode.SymbolKind.Class;
//                 break;
//             case "enum":
//                 kind = vscode.SymbolKind.Enum;
//                 break;
//             default:  
//                 console.log(v.name , v.Type , "have not match use default");
//                 kind = vscode.SymbolKind.Property;
//         }
        
//         var info = new vscode.SymbolInformation(
//             v.name , 
//             kind , 
//             new vscode.Range(
//                 new vscode.Position(v.pos.startLine , v.pos.startColumnOffset) ,
//                 new vscode.Position(v.pos.endLine , v.pos.endColumnOffset)
//                 ),
//             vscode.Uri.file(v.pos.filename),
//             );

//         infos.push(info);
//         if (v.inners) {
//             for(var j = 0 ;  j <  v.inners.length ; j++ ) {
//                 var vv =  v.inners[j];
//                 switch(vv.Type) {
//                     case "variable":
//                         kind = vscode.SymbolKind.Variable;
//                         break;
//                     case "function":
//                         kind = vscode.SymbolKind.Function;
//                         break;
//                     case "constant":
//                         kind = vscode.SymbolKind.Constant;
//                         break;
//                     case "class":
//                         kind = vscode.SymbolKind.Class;
//                         break;
//                     case "enum":
//                         kind = vscode.SymbolKind.Enum;
//                         break;
//                     case "method":
//                         kind = vscode.SymbolKind.Method;
//                         break;
//                     case "field":
//                         kind = vscode.SymbolKind.Field;
//                         break;
//                     case "enumItem":
//                         kind = vscode.SymbolKind.EnumMember;
//                         break;
//                     default:  
//                         console.log(vv.name , vv.Type , "have not match use default");
//                         kind = vscode.SymbolKind.Property;
//                 }
//                 var location22 = new vscode.Location(vscode.Uri.file(vv.pos.filename) , new vscode.Position(vv.pos.endLine , vv.pos.endColumnOffset)) ;
//                 var info2 = new vscode.SymbolInformation(vv.name , kind , v.name , location22);
//                 console.log(vv.name , v.name , info.name , info2.containerName);
//                 infos.push(info2);
//             }
//         }
//     }
//     return infos;
// };
