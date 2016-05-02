//通讯协议处理，主要处理封包和解包的过程
package protocol

import (
	"bytes"
	"encoding/binary"
)

type ProtoCommand interface  {

}

const (
	HeaderLength = 4 //包头长度4个字节，表示的是后面的包体长度
	ProtoCodeLength=4 //协议号4个字节
)
//封包 协议结构： 协议长度|协议号|协议内容
func Packet(protoCode int ,message []byte) []byte {
	message=append(IntToBytes(protoCode),message...) //先把协议号加到协议内容前面
	message=append(IntToBytes(len(message)),message...)
	return message
}
//解包
func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)
	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ HeaderLength {
			break
		}
		messageLength := BytesToInt(buffer[i : i+ HeaderLength]) //读取包头的长度信息
		if length < i+ HeaderLength +messageLength { //长度不够一个整包长度，直接跳出
			break
		}else { //足够整包，取出，并且把读取位置后移到取出这个整包之后的位置
			data := buffer[i+ HeaderLength : i+ HeaderLength +messageLength]
			readerChannel <- data //往通道里面放整个包的数据， 在取channel值的地方应该做一个协议的routing
			i += HeaderLength + messageLength - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}
//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
