package cg

//
//import (
//	"encoding/binary"
//)
//
//type AttributeBootstrapMethods struct {
//	methods []*BootStrapMethod
//}
//
//func (a *AttributeBootstrapMethods) ToAttributeInfo() *AttributeInfo {
//	ret := &AttributeInfo{}
//	ret.Info = make([]byte, 2)
//	binary.BigEndian.PutUint16(ret.Info, uint16(len(a.methods)))
//	for _, v := range a.methods {
//		bs := make([]byte, 4+len(v.bootStrapMethodArguments)*2)
//		binary.BigEndian.PutUint16(bs, uint16(v.bootStrapMethodRef))
//		binary.BigEndian.PutUint16(bs[2:], uint16(len(v.bootStrapMethodArguments)))
//		for kk, vv := range v.bootStrapMethodArguments {
//			binary.BigEndian.PutUint16(bs[4+kk*2:], uint16(vv))
//		}
//		ret.Info = append(ret.Info, bs...)
//	}
//	ret.attributeLength = uint32(len(ret.Info))
//	return ret
//}
//
//type BootStrapMethod struct {
//	bootStrapMethodRef       uint16
//	bootStrapMethodArguments []uint16
//}
