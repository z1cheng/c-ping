package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

// ICMPEchoPacket is the structure of ICMP echo packet without data sections
//
// ICMP Echo or Echo reply message
//  0               1               2               3
//  0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 8 9 A B C D E F
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Type      |     Code      |          Checksum             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |           Identifier          |        Sequence Number        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |  	                        Data...                            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type ICMPEchoPacket struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func main() {
	// The type of echo request is 8
	packet := ICMPEchoPacket{Type: 8}
	laddr := net.IPAddr{IP: net.ParseIP("0.0.0.0")}
	host := os.Args[1]
	raddr, _ := net.ResolveIPAddr("ip", host)
	var buffer bytes.Buffer

	// Generate checksum
	binary.Write(&buffer, binary.BigEndian, packet)
	packet.Checksum = ComputeChecksum(buffer.Bytes())
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, packet)

	conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
	if err != nil {
		fmt.Printf("Ping request could not find host %s. Please check the name and try again.\n", host)
		return
	}
	defer conn.Close()

	// Send the packet
	fmt.Printf("Pinging %s [%s]\n", host, raddr.String())
	if _, err := conn.Write(buffer.Bytes()); err != nil {
		fmt.Println(err.Error())
		return
	}
	conn.SetReadDeadline((time.Now().Add(3 * time.Second)))

	// Receive the reply
	startTime := time.Now()
	recv := make([]byte, 1024)
	_, err = conn.Read(recv)
	duration := time.Since(startTime).Milliseconds()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Check if the reply is correct
	if ComputeChecksum(recv) != 0 {
		fmt.Println("The checksum of the reply is incorrect")
		return
	}

	fmt.Printf("Reply from %s, time=%dms\n", raddr.String(), duration)
}

func ComputeChecksum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		cur    int
	)
	// Add each 16 bits integer to sum
	for length > 1 {
		sum += uint32(data[cur])<<8 + uint32(data[cur+1])
		cur += 2
		length -= 2
	}
	// If it's odd, add remaining byte
	if length > 0 {
		sum += uint32(data[cur])
	}
	// Add the high 16 bits to the low 16 bits
	for sum>>16 != 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}

	return uint16(^sum)
}
