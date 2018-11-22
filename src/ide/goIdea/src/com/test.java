package com;

import com.intellij.ide.ui.EditorOptionsTopHitProvider;
import com.intellij.openapi.actionSystem.AnAction;
import com.intellij.openapi.actionSystem.AnActionEvent;


public class test extends AnAction {
    public  test(){
        super("Text _Boxes");
    }
    @Override
    public void actionPerformed(AnActionEvent e)   {
        System.out.println("Hello World");
    }
}
