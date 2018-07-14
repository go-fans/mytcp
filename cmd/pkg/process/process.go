package process

import (
	"bytes"
	"encoding/binary"
)

//bytesToInt get msg length
func bytesToInt(b []byte)int{
	byteBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(byteBuffer, binary.BigEndian, &x)
	return int(x)
}

// UnPacket xxx
func UnPacket(buffer []byte, readerCh chan []byte)[]byte{
	length := len(buffer)
	var i int
	for i = 0;i< length;i++{
		if length < i + 4{		//TODO replace magic number with const value
			break
		}
		msgLength := bytesToInt(buffer[i:i+4])

		if length < i + 4 + msgLength{
			break
		}
		data := buffer[i+4:i+4+msgLength]
		readerCh <- data
		i += msgLength + 4 - 1
	}
	if i == length{
		return make([]byte,0)
	}
	return buffer[i:]
}


func intToByte(n int)[]byte{
	x := int32(n)
	byteBuffer := bytes.NewBuffer([]byte{})
	binary.Write(byteBuffer, binary.BigEndian,x)
	return byteBuffer.Bytes()
}

func Packet(msg []byte)[]byte{
	s := make([]byte,0)
	return append(append(s, intToByte(len(msg))...), msg...)
}
