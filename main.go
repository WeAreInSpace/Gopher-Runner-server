package main

import (
	"log"
	"net"
	"sync"

	"github.com/WeAreInSpace/Gopher-Runner-server/config"
	"github.com/WeAreInSpace/Gopher-Runner-server/network"
)

func main() {
	wg := new(sync.WaitGroup)
	mx := new(sync.Mutex)

	config := config.GetConfig()

	server, listenE := net.Listen("tcp", config.Address)
	if listenE != nil {
		log.Fatal(listenE)
	}

	wg.Add(1)

	go func() {
		for {
			conn, acceptE := server.Accept()
			if acceptE != nil {
				log.Printf("ERROR: %s\n", acceptE)
				if acceptE == net.ErrClosed {
					conn.Close()
					continue
				}
			}
			ib, og := network.HandleConn(conn)
			pm := network.PacketManager{
				Mx: mx,

				Conn: conn,
				Ib:   &ib,
				Og:   &og,
			}
			go pm.HandleClientConn()
		}
	}()
	wg.Wait()
	defer wg.Done()
}
