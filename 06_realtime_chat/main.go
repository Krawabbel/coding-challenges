package main

import "net"

func runServer(brok broker) error {

	listener, err := net.Listen("tcp", ":7007")
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go handleConn(brok, conn)
	}
}

func handleConn(brok broker, conn net.Conn) {

	logDbg("client %s connected", conn.RemoteAddr())

	go func() {
		logErr(handleSub(brok, conn))
	}()

	logErr(handlePub(brok, conn))

	logDbg("client %s disconnected", conn.RemoteAddr())

}

func main() {
	brok := newBroker()

	runnables := []func(broker) error{
		// runHeartbeat,
		runMonitor,
		runServer,
	}

	for _, run := range runnables {
		go func() {
			err := run(brok)
			logErr(err)
			brok.close()
		}()
	}

	<-brok.done()
}
