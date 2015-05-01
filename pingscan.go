package main

//
// Fast pingscan with JSON output
//
// ICMP echo part from https://github.com/golang/net/blob/master/icmp/ping_test.go
//

import (
	"errors"
	"flag"
	"fmt"
	"net"
	//"time"
	"github.com/analogic/pingscan/echo"
	"sync"
	"syscall"
	"time"

	// output
	"bytes"
	"encoding/json"
	"os"
)

type Host struct {
	Domain   string
	IP       *net.IP
	Sent     int64
	Received int64
	Err      error
	V6       bool
}

func (r *Host) RTT() int64 {
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
	// args
	timeout := flag.Int("timeout", 5, "DNS resolve + ping timeout (in s)")
	flag.Parse()
	domains := flag.Args()

	if len(domains) == 0 {
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("# pingscan -timeout=5 google.com yahoo.com ...")
		fmt.Println("")
		fmt.Println("You must set pingscan setuid bit, or set sysctl net.ipv4.ping_group_range=\"0   1000\", or run as root")
		fmt.Println("")
		os.Exit(1)
	}

	// we try if we have setuid bit
	syscall.Setuid(0)

	// real work
	results := ping(timeout, &domains)

	// json print
	js, _ := json.Marshal(results)
	var out bytes.Buffer
	json.Indent(&out, js, "", "\t")
	out.WriteTo(os.Stdout)
	fmt.Println("")
}

func ping(timeout *int, domains *[]string) *[]Host {

	hosts := make([]Host, len(*domains))
	retrieves := make([]chan int64, len(*domains))
	for i := range retrieves {
		retrieves[i] = make(chan int64)
	}

	for i, target := range *domains {
		hosts[i] = Host{Domain: target}
	}

	var wg sync.WaitGroup
	wg.Add(len(hosts))

	// work
	s := echo.StartSocket(false)

	// timed out receive
	for i := range hosts {
		go func(host *Host, retr *chan int64) {
			select {
			case <-time.After(time.Second * time.Duration(*timeout)):
				host.Err = errors.New("Echo timed out")
			case host.Received = <-*retr:
			}
			wg.Done()
		}(&hosts[i], &retrieves[i])
	}

	// receive echos from socket
	go func() {
		for {
			ip := <-s.In
			received := time.Now().UnixNano()

			go func() {
				for i, host := range hosts {
					if host.IP != nil && host.IP.String() == ip {
						retrieves[i] <- received
						return
					}

				}
			}()
		}
	}()

	// send echos
	for i := range hosts {

		go func(h *Host, retr *chan int64) {
			err := h.Resolve()
			if err != nil {
				h.Err = err
				*retr <- 0
			} else {
				h.Sent = time.Now().UnixNano()
				s.Echo(h.IP)
			}
		}(&hosts[i], &retrieves[i])

	}

	// output
	wg.Wait()

	return &hosts
}
