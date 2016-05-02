package main

import (
	"../protocol"
	"fmt"
)
func ParsePacket(u User ,data []byte) {
	length := len(data)
	if length < protocol.ProtoCodeLength {
		fmt.Println("数据长度少于头协议头长度")
		return
	}

	protoCode := protocol.BytesToInt(data[:protocol.ProtoCodeLength])
	if protoCode == 1000 {

	}
}

