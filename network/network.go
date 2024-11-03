package network

import (
	"encoding/json"
	"log"
	"math"
	"net"
	"sync"

	"github.com/WeAreInSpace/dotio"
)

func HandleConn(conn net.Conn) (ib dotio.Inbound, og dotio.Outgoing) {
	ib = dotio.Inbound{
		Conn: conn,
	}
	og = dotio.Outgoing{
		Conn: conn,
	}
	return ib, og
}

type PacketManager struct {
	Mx *sync.Mutex

	Conn net.Conn
	Ib   *dotio.Inbound
	Og   *dotio.Outgoing
}

type Player struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

func (pm *PacketManager) HandleClientConn() {
	defer pm.Conn.Close()

	for {
		log.Printf("INBOUND: waiting for receive packet from client\n")

		fristPacketId, _, readErr := pm.Ib.Read()
		if readErr != nil {
			log.Printf("CLOSE: %s", pm.Conn.RemoteAddr().String())
			break
		}

		log.Println(fristPacketId)

		log.Printf("INBOUND: received packet from client\n")

		if fristPacketId == 0 {
			log.Printf("INBOUND: from %s\n", pm.Conn.RemoteAddr().String())
			motd := pm.Og.Write()
			motd.WriteString("Welcome to gopher runner server")
			motd.WriteString("By We Are In Space")
			motd.Sent(dotio.WriteInt32(2))
			continue
		}

		if fristPacketId == 1 {
			log.Println(" - Handshake - ")
			log.Printf("From %s\n", pm.Conn.RemoteAddr().String())

			/*reqLoginId*/
			_, reqLogin, readReqLoginE := pm.Ib.Read()
			if readReqLoginE != nil {
				log.Printf("CLOSE: %s\n", pm.Conn.RemoteAddr().String())
				break
			}

			data := reqLogin.ReadString()

			player := Player{}
			JSONUnmarshalE := json.Unmarshal([]byte(data), &player)
			if JSONUnmarshalE != nil {
				log.Printf("CLOSE: %s\n", pm.Conn.RemoteAddr().String())
				break
			}

			loginRes := pm.Og.Write()
			loginRes.Sent(dotio.WriteInt32(0))

			playE := pm.Play(&player)
			if playE != nil {
				log.Printf("CLOSE: %s\n", pm.Conn.RemoteAddr().String())
				break
			}
		}

		if fristPacketId == math.MaxInt32 {
			log.Printf("CLOSE: %s\n", pm.Conn.RemoteAddr().String())
			break
		}
	}
}

func (pm *PacketManager) Play(player *Player) error {
	for {
		playId, _, playE := pm.Ib.Read()
		if playE != nil {
			return playE
		}

		if playId == 1 {
			_, playerPos, playerPosE := pm.Ib.Read()
			if playerPosE != nil {
				return playerPosE
			}
			getX := playerPos.ReadFloat64()
			getY := playerPos.ReadFloat64()
			log.Println(getX, getY)
		} else if playId == 2 {
			_, playerShift, playerShiftE := pm.Ib.Read()
			if playerShiftE != nil {
				return playerShiftE
			}
			playerShift.ReadInt64()
		} else {
			break
		}
	}
	return nil
}
