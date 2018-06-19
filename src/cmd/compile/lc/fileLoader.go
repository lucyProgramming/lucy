package lc

import (
	"encoding/binary"

	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/*
	load from file implementation
*/
type FileLoader struct {
	caches map[string]interface{}
}

/*
	lucy and java have no difference
*/
func (loader *FileLoader) loadInterfaces(astClass *ast.Class, c *cg.Class) error {
	astClass.InterfaceNames = make([]*ast.NameWithPos, len(c.Interfaces))
	for k, v := range c.Interfaces {
		astClass.InterfaceNames[k] = &ast.NameWithPos{
			Name: string(c.ConstPool[v].Info),
		}
	}
	astClass.Interfaces = make([]*ast.Class, len(astClass.InterfaceNames))
	for k, c := range astClass.InterfaceNames {
		i := &ast.Class{}
		i.Name = c.Name
		i.NotImportedYet = true
		astClass.Interfaces[k] = i
	}
	return nil
}

func (loader *FileLoader) loadAsJava(c *cg.Class) (*ast.Class, error) {
	//name
	if t := c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_SIGNATURE); t != nil && len(t) > 0 {
		//TODO:: support signature???
	}
	astClass := &ast.Class{}
	{
		nameIndex := binary.BigEndian.Uint16(c.ConstPool[c.ThisClass].Info)
		astClass.Name = string(c.ConstPool[nameIndex].Info)
		if astClass.Name != ast.JAVA_ROOT_CLASS {
			nameIndex = binary.BigEndian.Uint16(c.ConstPool[c.SuperClass].Info)
			astClass.SuperClassName = string(c.ConstPool[nameIndex].Info)
		}
	}
	loader.loadInterfaces(astClass, c)
	astClass.AccessFlags = c.AccessFlag
	astClass.IsJava = true // class compiled from java
	var err error
	astClass.Fields = make(map[string]*ast.ClassField)
	astClass.LoadFromOutSide = true
	for _, v := range c.Fields {
		f := &ast.ClassField{}
		f.LoadFromOutSide = true
		f.AccessFlags = v.AccessFlags
		f.JvmDescriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		f.Name = string(c.ConstPool[v.NameIndex].Info)
		_, f.Type, err = jvm.Descriptor.ParseType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		astClass.Fields[f.Name] = f
	}
	astClass.Methods = make(map[string][]*ast.ClassMethod)
	for _, v := range c.Methods {
		m := &ast.ClassMethod{}
		m.Func = &ast.Function{}
		m.LoadFromOutSide = true
		m.Func.Name = string(c.ConstPool[v.NameIndex].Info)
		m.Func.Descriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		m.Func.AccessFlags = v.AccessFlags
		m.Func.Type, err = jvm.Descriptor.ParseFunctionType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_METHOD_PARAMETERS); t != nil && len(t) > 0 {
			parseMethodParameter(c, t[0].Info, m.Func)
		}
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_RETURNLIST_NAMES); t != nil && len(t) > 0 {
			parseReturnListNames(c, t[0].Info, m.Func)
		}
		if astClass.Methods[m.Func.Name] == nil {
			astClass.Methods[m.Func.Name] = []*ast.ClassMethod{m}
		} else {
			astClass.Methods[m.Func.Name] = append(astClass.Methods[m.Func.Name], m)
		}
	}
	return astClass, nil
}

func (loader *FileLoader) loadAsLucy(c *cg.Class) (*ast.Class, error) {
	if t := c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_COMPILTER_AUTO); t != nil && len(t) > 0 {
		return nil, nil
	}
	//name
	astClass := &ast.Class{}
	{
		nameIndex := binary.BigEndian.Uint16(c.ConstPool[c.ThisClass].Info)
		astClass.Name = string(c.ConstPool[nameIndex].Info)
		if astClass.Name != ast.JAVA_ROOT_CLASS {
			nameIndex = binary.BigEndian.Uint16(c.ConstPool[c.SuperClass].Info)
			astClass.SuperClassName = string(c.ConstPool[nameIndex].Info)
		}
	}
	loader.loadInterfaces(astClass, c)
	astClass.AccessFlags = c.AccessFlag
	astClass.LoadFromOutSide = true
	var err error
	astClass.Fields = make(map[string]*ast.ClassField)
	for _, v := range c.Fields {
		f := &ast.ClassField{}
		f.Name = string(c.ConstPool[v.NameIndex].Info)
		f.JvmDescriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		f.LoadFromOutSide = true
		_, f.Type, err = jvm.Descriptor.ParseType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_FIELD_DESCRIPTOR); t != nil && len(t) > 0 {
			index := binary.BigEndian.Uint16(t[0].Info)
			_, f.Type, err = jvm.LucyFieldSignatureParser.Decode(c.ConstPool[index].Info)
			if err != nil {
				return nil, err
			}
		}
		if f.Type.Type == ast.VARIABLE_TYPE_ENUM {
			loadEnumForVariableType(f.Type)
		}
		f.AccessFlags = v.AccessFlags
		astClass.Fields[f.Name] = f
	}
	astClass.Methods = make(map[string][]*ast.ClassMethod)
	for _, v := range c.Methods {
		m := &ast.ClassMethod{}
		m.Func = &ast.Function{}
		m.Func.Name = string(c.ConstPool[v.NameIndex].Info)
		m.Func.Type, err = jvm.Descriptor.ParseFunctionType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		m.Func.AccessFlags = v.AccessFlags
		m.LoadFromOutSide = true
		m.Func.Descriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_METHOD_DESCRIPTOR); t != nil && len(t) > 0 {
			index := binary.BigEndian.Uint16(t[0].Info)
			err = jvm.LucyMethodSignatureParser.Decode(m.Func, c.ConstPool[index].Info)
			if err != nil {
				return nil, err
			}
		}
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_DEFAULT_PARAMETERS); t != nil && len(t) > 0 {
			dp := &cg.AttributeDefaultParameters{}
			dp.FromBytes(t[0].Info)
			jvm.FunctionDefaultValueParser.Decode(c, m.Func, dp)
		}
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_METHOD_PARAMETERS); t != nil && len(t) > 0 {
			parseMethodParameter(c, t[0].Info, m.Func)
		}
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_RETURNLIST_NAMES); t != nil && len(t) > 0 {
			parseReturnListNames(c, t[0].Info, m.Func)
		}
		err = loadEnumForFunction(m.Func)
		if err != nil {
			return nil, err
		}
		if astClass.Methods[m.Func.Name] == nil {
			astClass.Methods[m.Func.Name] = []*ast.ClassMethod{m}
		} else {
			astClass.Methods[m.Func.Name] = append(astClass.Methods[m.Func.Name], m)
		}
	}
	return astClass, nil
}

func (loader *FileLoader) loadLucyEnum(c *cg.Class) (*ast.Enum, error) {
	e := &ast.Enum{}
	{
		nameIndex := binary.BigEndian.Uint16(c.ConstPool[c.ThisClass].Info)
		e.Name = string(c.ConstPool[nameIndex].Info)
	}
	e.AccessFlags = c.AccessFlag
	for _, v := range c.Fields {
		en := &ast.EnumName{}
		name := string(c.ConstPool[v.NameIndex].Info)
		en.Name = name
		en.Enum = e
		constValue := v.AttributeGroupedByName[cg.ATTRIBUTE_NAME_CONST_VALUE][0] // must have this attribute
		en.Value = int32(binary.BigEndian.Uint32(c.ConstPool[binary.BigEndian.Uint16(constValue.Info)].Info))
		e.Enums = append(e.Enums, en)
	}
	return e, nil
}

func (loader *FileLoader) loadLucyMainClass(pack *ast.Package, c *cg.Class) error {
	var err error
	mainClassName := &cg.ClassHighLevel{}
	mainClassName.Name = pack.Name + "/main"
	pack.Block.Variables = make(map[string]*ast.VariableDefinition)
	pack.Block.Constants = make(map[string]*ast.Constant)
	pack.Block.Functions = make(map[string]*ast.Function)
	for _, f := range c.Fields {
		name := string(c.ConstPool[f.NameIndex].Info)
		constValue := f.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_CONST_VALUE)
		if len(constValue) > 1 {
			return fmt.Errorf("constant value length greater than  1 at class 'main'  field '%s'", name)
		}
		_, typ, err := jvm.Descriptor.ParseType(c.ConstPool[f.DescriptorIndex].Info)
		if err != nil {
			return err
		}
		if len(f.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_CONST)) > 0 {
			//const
			cos := &ast.Constant{}
			cos.Name = name
			cos.AccessFlags = f.AccessFlags
			cos.Type = typ
			_, cos.Type, err = jvm.Descriptor.ParseType(c.ConstPool[f.DescriptorIndex].Info)
			if err != nil {
				return err
			}
			valueIndex := binary.BigEndian.Uint16(constValue[0].Info)
			switch cos.Type.Type {
			case ast.VARIABLE_TYPE_BOOL:
				cos.Value = binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info) != 0
			case ast.VARIABLE_TYPE_BYTE:
				cos.Value = byte(binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_SHORT:
				fallthrough
			case ast.VARIABLE_TYPE_INT:
				cos.Value = int32(binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_LONG:
				cos.Value = int64(binary.BigEndian.Uint64(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_FLOAT:
				cos.Value = float32(binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_DOUBLE:
				cos.Value = float64(binary.BigEndian.Uint64(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_STRING:
				valueIndex = binary.BigEndian.Uint16(c.ConstPool[valueIndex].Info) // const_string_info
				cos.Value = string(c.ConstPool[valueIndex].Info)                   // utf 8
			}
			pack.Block.Constants[name] = cos
		} else {
			//global vars
			vd := &ast.VariableDefinition{}
			vd.Name = name
			vd.AccessFlags = f.AccessFlags
			vd.JvmDescriptor = string(c.ConstPool[f.DescriptorIndex].Info)
			vd.Type = typ
			vd.IsGlobal = true
			pack.Block.Variables[name] = vd
			if t := f.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_FIELD_DESCRIPTOR); t != nil && len(t) > 0 {
				index := binary.BigEndian.Uint16(t[0].Info)
				_, vd.Type, err = jvm.LucyFieldSignatureParser.Decode(c.ConstPool[index].Info)
				if err != nil {
					return err
				}
			}
			if typ.Type == ast.VARIABLE_TYPE_ENUM {
				loadEnumForVariableType(typ)
			}
		}
	}
	for _, m := range c.Methods {
		if t := m.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_TRIGGER_PACKAGE_INIT); t != nil && len(t) > 0 {
			pack.TriggerPackageInitMethodName = string(c.ConstPool[m.NameIndex].Info)
			continue
		}
		if m.AccessFlags&cg.ACC_METHOD_BRIDGE != 0 {
			continue
		}
		name := string(c.ConstPool[m.NameIndex].Info)
		if name == ast.MAIN_FUNCTION_NAME {
			// this is main function
			continue
		}
		function := &ast.Function{}
		function.Name = name
		function.AccessFlags = m.AccessFlags
		function.Descriptor = string(c.ConstPool[m.DescriptorIndex].Info)
		function.Type, err = jvm.Descriptor.ParseFunctionType(c.ConstPool[m.DescriptorIndex].Info)
		if err != nil {
			return err
		}
		if t := m.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_METHOD_DESCRIPTOR); t != nil && len(t) > 0 {
			index := binary.BigEndian.Uint16(t[0].Info)
			err = jvm.LucyMethodSignatureParser.Decode(function, c.ConstPool[index].Info)
			if err != nil {
				return err
			}
		}
		err = loadEnumForFunction(function)
		if err != nil {
			return err
		}
		if t := m.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_METHOD_PARAMETERS); t != nil && len(t) > 0 {
			parseMethodParameter(c, t[0].Info, function)
		}
		if t := m.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_RETURNLIST_NAMES); t != nil && len(t) > 0 {
			parseReturnListNames(c, t[0].Info, function)
		}
		function.ClassMethod = &cg.MethodHighLevel{}
		function.ClassMethod.Name = function.Name
		function.ClassMethod.Class = mainClassName
		function.ClassMethod.Descriptor = function.Descriptor
		function.IsGlobal = true
		pack.Block.Functions[name] = function
	}
	if pack.Block.TypeAlias == nil {
		pack.Block.TypeAlias = make(map[string]*ast.VariableType)
	}
	for _, v := range c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_TYPE_ALIAS) {
		index := binary.BigEndian.Uint16(v.Info)
		name, typ, err := jvm.LucyTypeAliasParser.Decode(c.ConstPool[index].Info)
		if err != nil {
			return err
		}
		typ.Alias = name
		pack.Block.TypeAlias[name] = typ
		if typ.Type == ast.VARIABLE_TYPE_ENUM {
			err = loadEnumForVariableType(typ)
			if err != nil {
				return err
			}
		}
	}
	for _, v := range c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_TEMPLATE_FUNCTION) {
		attr := &cg.AttributeTemplateFunction{}
		attr.FromBytes(c, v.Info)
		f, es := ParseFunctionHandler([]byte(attr.Code), &ast.Pos{
			Filename:    attr.Filename,
			StartLine:   int(attr.StartLine),
			StartColumn: int(attr.StartColumn),
		})
		if len(es) > 0 { // looks impossible
			return es[0]
		}
		f.TemplateFunction = &ast.TemplateFunction{}
		pack.Block.Functions[attr.Name] = f
	}
	return nil
}

func (loader *FileLoader) loadLucyPackage(r *Resource) (*ast.Package, error) {
	fis, err := ioutil.ReadDir(r.realPath)
	if err != nil {
		return nil, err
	}
	fisM := make(map[string]os.FileInfo)
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".class") {
			fisM[v.Name()] = v
		}
	}
	_, ok := fisM[mainClassName]
	if ok == false {
		return nil, fmt.Errorf("main class not found")
	}
	bs, err := ioutil.ReadFile(filepath.Join(r.realPath, mainClassName))
	if err != nil {
		return nil, fmt.Errorf("read main.class failed,err:%v", err)
	}
	c, err := (&ClassDecoder{}).decode(bs)
	if err != nil {
		return nil, fmt.Errorf("decode main class failed,err:%v", err)
	}
	p := &ast.Package{}
	p.Name = r.name
	err = loader.loadLucyMainClass(p, c)
	if err != nil {
		return nil, fmt.Errorf("parse main class failed,err:%v", err)
	}
	delete(fisM, mainClassName)
	mkEnums := func(e *ast.Enum) {
		if p.Block.Enums == nil {
			p.Block.Enums = make(map[string]*ast.Enum)
		}
		if p.Block.EnumNames == nil {
			p.Block.EnumNames = make(map[string]*ast.EnumName)
		}
		p.Block.Enums[filepath.Base(e.Name)] = e
		for _, v := range e.Enums {
			p.Block.EnumNames[v.Name] = v
		}
	}
	for _, v := range fisM {
		bs, err := ioutil.ReadFile(filepath.Join(r.realPath, v.Name()))
		if err != nil {
			return p, fmt.Errorf("read class failed,err:%v", err)
		}
		c, err := (&ClassDecoder{}).decode(bs)
		if err != nil {
			return nil, fmt.Errorf("decode class failed,err:%v", err)
		}
		if len(c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_ENUM)) > 0 {
			e, err := loader.loadLucyEnum(c)
			if err != nil {
				return nil, err
			}
			mkEnums(e)
			continue
		}
		class, err := loader.loadAsLucy(c)
		if err != nil {
			return nil, fmt.Errorf("decode class failed,err:%v", err)
		}
		if p.Block.Classes == nil {
			p.Block.Classes = make(map[string]*ast.Class)
		}
		p.Block.Classes[filepath.Base(class.Name)] = class
	}
	return p, nil
}

func (loader *FileLoader) loadJavaPackage(r *Resource) (*ast.Package, error) {
	fis, err := ioutil.ReadDir(r.realPath)
	if err != nil {
		return nil, err
	}
	ret := &ast.Package{}
	ret.Block.Classes = make(map[string]*ast.Class)
	for _, v := range fis {
		var rr Resource
		rr.kind = RESOURCE_KIND_JAVA_CLASS
		if strings.HasSuffix(v.Name(), ".class") == false || strings.Contains(v.Name(), "$") {
			continue
		}
		rr.realPath = filepath.Join(r.realPath, v.Name())
		class, err := loader.loadClass(&rr)
		if err != nil {
			return nil, err
		}
		if c, ok := class.(*ast.Class); ok && class != nil {
			ret.Block.Classes[filepath.Base(c.Name)] = c
		}
	}
	return ret, nil
}

func (loader *FileLoader) loadClass(r *Resource) (interface{}, error) {
	bs, err := ioutil.ReadFile(r.realPath)
	if err != nil {
		return nil, err
	}
	c, err := (&ClassDecoder{}).decode(bs)
	if r.kind == RESOURCE_KIND_LUCY_CLASS {
		if t := c.AttributeGroupedByName[cg.ATTRIBUTE_NAME_LUCY_ENUM]; len(t) > 0 {
			return loader.loadLucyEnum(c)
		} else {
			return loader.loadAsLucy(c)
		}
	}
	return loader.loadAsJava(c)
}

func (loader *FileLoader) LoadImport(importName string) (interface{}, error) {
	if loader.caches != nil && loader.caches[importName] != nil {
		return loader.caches[importName], nil
	}
	var realPaths []*Resource
	for _, v := range compiler.lucyPaths {
		p := filepath.Join(v, "class", importName)
		f, err := os.Stat(p)
		if err == nil && f.IsDir() { // directory is package
			realPaths = append(realPaths, &Resource{
				kind:     RESOURCE_KIND_LUCY_PACKAGE,
				realPath: p,
				name:     importName,
			})
		}
		p = filepath.Join(v, "class", importName+".class")
		f, err = os.Stat(p)
		if err == nil && f.IsDir() == false { // class file
			realPaths = append(realPaths, &Resource{
				kind:     RESOURCE_KIND_LUCY_CLASS,
				realPath: p,
				name:     importName,
			})
		}
	}
	for _, v := range compiler.ClassPaths {
		p := filepath.Join(v, importName)
		f, err := os.Stat(p)
		if err == nil && f.IsDir() { // directory is package
			realPaths = append(realPaths, &Resource{
				kind:     RESOURCE_KIND_JAVA_PACKAGE,
				realPath: p,
				name:     importName,
			})
		}
		p = filepath.Join(v, importName+".class")
		f, err = os.Stat(p)
		if err == nil && f.IsDir() == false { // directory is package
			realPaths = append(realPaths, &Resource{
				kind:     RESOURCE_KIND_JAVA_CLASS,
				realPath: p,
				name:     importName,
			})
		}
	}
	if len(realPaths) == 0 {
		return nil, fmt.Errorf("resource '%v' not found", importName)
	}
	realPathMap := make(map[string][]*Resource)
	for _, v := range realPaths {
		_, ok := realPathMap[v.realPath]
		if ok {
			realPathMap[v.realPath] = append(realPathMap[v.realPath], v)
		} else {
			realPathMap[v.realPath] = []*Resource{v}
		}
	}
	if len(realPathMap) > 1 {
		errMsg := "not 1 resource named '" + importName + "' present:\n"
		for _, v := range realPathMap {
			switch v[0].kind {
			case RESOURCE_KIND_JAVA_CLASS:
				errMsg += fmt.Sprintf("\t in '%s' is a java class\n", v[0].realPath)
			case RESOURCE_KIND_JAVA_PACKAGE:
				errMsg += fmt.Sprintf("\t in '%s' is a java package\n", v[0].realPath)
			case RESOURCE_KIND_LUCY_CLASS:
				errMsg += fmt.Sprintf("\t in '%s' is a lucy class\n", v[0].realPath)
			case RESOURCE_KIND_LUCY_PACKAGE:
				errMsg += fmt.Sprintf("\t in '%s' is a lucy package\n", v[0].realPath)
			}
		}
		return nil, fmt.Errorf(errMsg)
	}
	if realPaths[0].kind == RESOURCE_KIND_LUCY_CLASS {
		if filepath.Base(realPaths[0].realPath) == mainClassName {
			return nil, fmt.Errorf("%s is special class for global variable and other things", mainClassName)
		}
	}
	if realPaths[0].kind == RESOURCE_KIND_JAVA_CLASS {
		class, err := loader.loadClass(realPaths[0])
		if class != nil {
			loader.caches[importName] = class
		}
		return class, err
	} else if realPaths[0].kind == RESOURCE_KIND_LUCY_CLASS {
		t, err := loader.loadClass(realPaths[0])
		if t != nil {
			loader.caches[importName] = t
		}
		return t, err
	} else if realPaths[0].kind == RESOURCE_KIND_JAVA_PACKAGE {
		p, err := loader.loadJavaPackage(realPaths[0])
		if p != nil {
			loader.caches[importName] = p
		}
		return p, err
	} else { // lucy package
		p, err := loader.loadLucyPackage(realPaths[0])
		if p != nil {
			loader.caches[importName] = p
		}
		return p, err
	}
}
