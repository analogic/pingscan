package main

//
// Fast pingscan with JSON output
//
// ICMP echo part from https://github.com/golang/net/blob/master/icmp/ping_test.go
//

import (
    "errors"
    "net"
    "flag"
    "fmt"
    //"time"
    "sync"
    "github.com/analogic/pingscan/echo"
    "time"

    // output
    "encoding/json"
    "bytes"
    "os"
)

type Host struct {
    Domain string
    IP *net.IP
    Sent int64
    Received int64
    Err error
    V6 bool
}

func (r *Host) RTT() (int64) {
    return r.Received - r.Sent
}

func (h *Host) Resolve() (err error) {
    ips, err := net.LookupIP(h.Domain)
    if err != nil {
        h.Err = err
        return err
    }

    for _, ip := range ips {
        if (h.V6 && (ip.To4() == nil && ip.To16() != nil)) || (!h.V6 && ip.To4() != nil) {
            h.IP = &ip
            return nil
        }
    }

    h.Err = errors.New("No IP address found")
    return h.Err
}

func main() {
	//timeout := flag.Int("timeout", 5, "Single dns resolve + ping timeout (in s)")

    // prepare
	flag.Parse()
	domains := flag.Args()

    hosts := make([]Host, len(domains))
    for i, target := range domains {
        hosts[i] = Host{Domain: target}
    }

    var wg sync.WaitGroup
    wg.Add(len(hosts))

    // work

    s := echo.StartSocket(false)

    // receive
    go func() {
        for {
            ip := <- s.In
            received := time.Now().UnixNano()

            go func() {
                for i, host := range hosts {
                    if host.IP != nil && host.IP.String() == ip {

                        hosts[i].Received = received
                        //fmt.Println("Packet od", ip, hosts[i].RTT())
                        wg.Done()

                        return
                    }

                }
            } ()
        }
    } ()

    // send
    for i := range hosts {

        go func(h *Host) {
            go func() {
                time.Sleep(5 * time.Second)
            } ()

            if h.Resolve() != nil {
                wg.Done()
            } else {
                h.Sent = time.Now().UnixNano()
                s.Echo(h.IP)
            }
        } (&hosts[i])

    }



    // output
    wg.Wait()

    // json print
    js, _ := json.Marshal(hosts)
    //fmt.Println(string(js))
    var out bytes.Buffer
    json.Indent(&out, js, "", "\t")
    out.WriteTo(os.Stdout)
    fmt.Println("")
}