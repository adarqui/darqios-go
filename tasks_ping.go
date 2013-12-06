package main

import (
	"fmt"
	"bytes"
	"os"
	"errors"
	"time"
	"net"
)


func Task_Ping(M *Main, S *State, TD *Task_Data) (bool) {
	Debug("Task_Ping:Entered\n")

	if os.Getuid() != 0 || len(TD.Policy.Params) == 0 {
		return false
	}

	found := false

	for _,param := range TD.Policy.Params {
		Debug("PINGING: %v\n", param)
		truth := PING_Ipv4(param, 1)

		if truth == false {
			found = true
			mon := S.MON_Gen_Task("high", param, TD.Policy, fmt.Sprintf("Unable to ping %s", param), "None.")

			M.M<-mon
		}
	}

	return found
}

func PING_Ipv4(IP string, Timeout time.Duration) (bool) {
	c, err := net.Dial("ip4:icmp", IP)
	if err != nil {
		return false
	}

	c.SetDeadline(time.Now().Add(Timeout * time.Second))
	defer c.Close()

	typ := icmpv4EchoRequest

	i := 0

	xid, xseq := os.Getpid()&0xffff, i+1
	wb, err := (&icmpMessage{
		Type: typ, Code: 0,
		Body: &icmpEcho{
			ID: xid, Seq: xseq,
			Data: bytes.Repeat([]byte("Go Go Gadget Ping!!!"), 3),
		},
	}).Marshal()
	if err != nil {
		return false
	}
	if _, err := c.Write(wb); err != nil {
		return false
	}
	var m *icmpMessage
	rb := make([]byte, 20+len(wb))
	for {
		if _, err := c.Read(rb); err != nil {
			return false
		}
		rb = ipv4Payload(rb)
		if m, err = parseICMPMessage(rb); err != nil {
			return false
		}
		switch m.Type {
		case icmpv4EchoRequest, icmpv6EchoRequest:
			continue
		}
		break
	}
	switch p := m.Body.(type) {
	case *icmpEcho:
		if p.ID != xid || p.Seq != xseq {
			return false
		}
	default:
		return false
	}

	return true
}


func ipv4Payload(b []byte) []byte {
	if len(b) < 20 {
		return b
	}
	hdrlen := int(b[0]&0x0f) << 2
	return b[hdrlen:]
}


const (
	icmpv4EchoRequest = 8
	icmpv4EchoReply   = 0
	icmpv6EchoRequest = 128
	icmpv6EchoReply   = 129
)

// icmpMessage represents an ICMP message.
type icmpMessage struct {
	Type     int             // type
	Code     int             // code
	Checksum int             // checksum
	Body     icmpMessageBody // body
}

// icmpMessageBody represents an ICMP message body.
type icmpMessageBody interface {
	Len() int
	Marshal() ([]byte, error)
}

// Marshal returns the binary enconding of the ICMP echo request or
// reply message m.
func (m *icmpMessage) Marshal() ([]byte, error) {
	b := []byte{byte(m.Type), byte(m.Code), 0, 0}
	if m.Body != nil && m.Body.Len() != 0 {
		mb, err := m.Body.Marshal()
		if err != nil {
			return nil, err
		}
		b = append(b, mb...)
	}
	switch m.Type {
	case icmpv6EchoRequest, icmpv6EchoReply:
		return b, nil
	}
	csumcv := len(b) - 1 // checksum coverage
	s := uint32(0)
	for i := 0; i < csumcv; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if csumcv&1 == 0 {
		s += uint32(b[csumcv])
	}
	s = s>>16 + s&0xffff
	s = s + s>>16
	// Place checksum back in header; using ^= avoids the
	// assumption the checksum bytes are zero.
	b[2] ^= byte(^s)
	b[3] ^= byte(^s >> 8)
	return b, nil
}

// parseICMPMessage parses b as an ICMP message.
func parseICMPMessage(b []byte) (*icmpMessage, error) {
	msglen := len(b)
	if msglen < 4 {
		return nil, errors.New("message too short")
	}
	m := &icmpMessage{Type: int(b[0]), Code: int(b[1]), Checksum: int(b[2])<<8 | int(b[3])}
	if msglen > 4 {
		var err error
		switch m.Type {
		case icmpv4EchoRequest, icmpv4EchoReply, icmpv6EchoRequest, icmpv6EchoReply:
			m.Body, err = parseICMPEcho(b[4:])
			if err != nil {
				return nil, err
			}
		}
	}
	return m, nil
}

// imcpEcho represenets an ICMP echo request or reply message body.
type icmpEcho struct {
	ID   int    // identifier
	Seq  int    // sequence number
	Data []byte // data
}

func (p *icmpEcho) Len() int {
	if p == nil {
		return 0
	}
	return 4 + len(p.Data)
}

// Marshal returns the binary enconding of the ICMP echo request or
// reply message body p.
func (p *icmpEcho) Marshal() ([]byte, error) {
	b := make([]byte, 4+len(p.Data))
	b[0], b[1] = byte(p.ID>>8), byte(p.ID)
	b[2], b[3] = byte(p.Seq>>8), byte(p.Seq)
	copy(b[4:], p.Data)
	return b, nil
}

// parseICMPEcho parses b as an ICMP echo request or reply message
// body.
func parseICMPEcho(b []byte) (*icmpEcho, error) {
	bodylen := len(b)
	p := &icmpEcho{ID: int(b[0])<<8 | int(b[1]), Seq: int(b[2])<<8 | int(b[3])}
	if bodylen > 4 {
		p.Data = make([]byte, bodylen-4)
		copy(p.Data, b[4:])
	}
	return p, nil
}
