package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func resolveHost(host string) (net.IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return ip, nil
		}
	}

	return nil, fmt.Errorf("no IPv4 found for %s", host)
}

func resolveIP(ip net.Addr) (string, error) {

	names, err := net.LookupAddr(ip.String())
	if err != nil {
		return "", err
	}

	return strings.Join(names, ", "), nil
}

func main() {
	if len(os.Args) != 2 {
		panic("usage: gotraceroute host")
	}

	hostName := os.Args[1]

	hostIP, err := resolveHost(hostName)
	must(err)

	maxHops := 64

	payload := []byte("codingchallenges.fyi trace route")

	fmt.Printf("traceroute to %s (%v), %d hops max, %d byte packets\n", hostName, hostIP, maxHops, len(payload))

	hostPort := "33434"
	host, err := net.ResolveUDPAddr("udp4", hostIP.String()+":"+hostPort)
	must(err)

	conn, err := net.DialUDP("udp", nil, host)
	must(err)
	defer conn.Close()

	ipv4Conn := ipv4.NewPacketConn(conn)

loop:
	for i := 1; i <= maxHops; i++ {
		must(ipv4Conn.SetTTL(i))

		routes := make(chan Route)
		ready := make(chan Signal)
		go listenForICMP(routes, ready)
		<-ready

		start := time.Now()
		_, err := conn.Write(payload)
		must(err)

		route := <-routes

		switch route.typ {
		case ROUTE_TIMEOUT:
			fmt.Printf("%d * * *\n", i)
		case ROUTE_INTERMEDIATE, ROUTE_FINAL:
			routeName, err := resolveIP(route.addr)
			if err != nil {
				routeName = route.addr.String()
			}
			rtt := route.timestamp.Sub(start)
			fmt.Printf("%d %s (%s) %v\n", i, routeName, route.addr, rtt)
			if route.typ == ROUTE_FINAL {
				break loop
			}
		default:
			panic("unexpected route type")
		}
	}
}

const (
	ROUTE_INTERMEDIATE = iota
	ROUTE_TIMEOUT
	ROUTE_FINAL
)

type Route struct {
	addr      net.Addr
	typ       int
	timestamp time.Time
}

type Signal struct{}

var signal Signal

func listenForICMP(routes chan Route, ready chan Signal) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	must(err)
	defer conn.Close()

	ready <- signal

	buffer := make([]byte, 1500)
	conn.SetReadDeadline(time.Now().Add(time.Second))

	n, peer, err := conn.ReadFrom(buffer)
	timestamp := time.Now()

	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			routes <- Route{typ: ROUTE_TIMEOUT, timestamp: timestamp}
			return
		}
		panic(err)
	}

	msg, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), buffer[:n])
	must(err)

	switch msg.Type {
	case ipv4.ICMPTypeTimeExceeded:
		routes <- Route{addr: peer, typ: ROUTE_INTERMEDIATE, timestamp: timestamp}
		return
	case ipv4.ICMPTypeDestinationUnreachable:
		routes <- Route{addr: peer, typ: ROUTE_FINAL, timestamp: timestamp}
		return
	default:
	}
	fmt.Printf("[DEBUG] ignoring message %v of type %v (proto: %d) from %v\n", msg.Body, msg.Type, msg.Type.Protocol(), peer)
}
