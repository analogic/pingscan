package echo

import (
    "os"
    "errors"
    "net"
    "golang.org/x/net/icmp"
    "golang.org/x/net/internal/iana"
    "golang.org/x/net/ipv4"
    "golang.org/x/net/ipv6"
    "fmt"
)

func StartSocket(v6 bool) (s *Socket) {
    s = &Socket{V6: v6}
    s.Start()
    go s.Listen()

    return s
}

type Socket struct {
    C *icmp.PacketConn
    Out chan string
    In chan string
    V6 bool
}

func (es *Socket) Start() {
    var err error
    es.Out = make(chan string)
    es.In = make(chan string)

    if es.V6 == true {
        es.C, err = icmp.ListenPacket("udp6", "::")
        if err != nil {
            es.C, err = icmp.ListenPacket("ip6:ipv6-icmp", "::")
            if err != nil {
                panic(err)
            }
        }
    } else {
        es.C, err = icmp.ListenPacket("udp4", "0.0.0.0")
        if err != nil {
            es.C, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
            if err != nil {
                panic(err)
            }
        }
    }
}

func (es *Socket) Listen() {
    for {
        rb := make([]byte, 1500)
        n, peer, err := es.C.ReadFrom(rb)

        if err != nil {
            panic(err)
        }

        go es.handlePacket(rb[:n],  peer);
    }
}

func (es *Socket) handlePacket(packet []byte, peer net.Addr) {
    var protocol int

    if es.V6 {
        protocol = iana.ProtocolIPv6ICMP
    } else {
        protocol = iana.ProtocolICMP
    }

    rm, err := icmp.ParseMessage(protocol, packet); if err != nil {
        return
    }

    switch rm.Type {
        case ipv4.ICMPTypeEchoReply, ipv6.ICMPTypeEchoReply:
        es.In <- peer.String()
    }

}

func (es *Socket) Echo(ip *net.IP) (err error) {
    var messageType icmp.Type

    if !es.V6 {
        messageType = ipv4.ICMPTypeEcho
    } else {
        messageType = ipv6.ICMPTypeEchoRequest
    }

    // create message for ping
    wm := icmp.Message{
        Type: messageType, Code: 0,
        Body: &icmp.Echo{
            ID: os.Getpid() & 0xffff, Seq: 1 << 1,
            Data: []byte("HELLO-R-U-THERE"),
        },
    }

    // encode
    wb, err := wm.Marshal(nil)
    if err != nil {
        return err
    }

    // send
    if n, err := es.C.WriteTo(wb, &net.IPAddr{IP: *ip}); err != nil {
        return err
    } else if n != len(wb) {
        return errors.New(fmt.Sprintf("got %v; want %v", n, len(wb)))
    } else {
        return nil
    }
}