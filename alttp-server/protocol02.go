package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

type P02Kind byte

const (
	RequestIndex      = P02Kind(0x00)
	Broadcast         = P02Kind(0x01)
	BroadcastToSector = P02Kind(0x02)
	RequestSnapshot   = P02Kind(0x03)
	SnapshotData      = P02Kind(0x04)
)

func (k P02Kind) String() string {
	switch k {
	case RequestIndex:
		return "request_index"
	case Broadcast:
		return "broadcast"
	case BroadcastToSector:
		return "broadcast_to_sector"
	case RequestSnapshot:
		return "request_snapshot"
	case SnapshotData:
		return "snapshot_data"
	}
	return "unknown"
}

func make02Packet(groupBuf []byte, kind P02Kind) (buf *bytes.Buffer) {
	// construct message:
	buf = &bytes.Buffer{}
	header := uint16(25887)
	binary.Write(buf, binary.LittleEndian, &header)
	protocol := byte(0x02)
	buf.WriteByte(protocol)

	// protocol packet:
	buf.Write(groupBuf)
	responseKind := kind | 0x80
	buf.WriteByte(byte(responseKind))

	return
}

func processProtocol02(message UDPMessage, buf *bytes.Buffer) (fatalErr error) {
	groupBuf := make([]byte, 20)
	_, err := buf.Read(groupBuf)
	if err != nil {
		log.Print(err)
		return
	}
	group := string(groupBuf[:])
	if len(group) != 20 {
		log.Fatal("bug! group name must be exactly 20 bytes")
		return
	}

	var kind P02Kind
	if err := binary.Read(buf, binary.LittleEndian, &kind); err != nil {
		log.Print(err)
		return
	}

	// what the client thinks its index is:
	var index uint16
	if err := binary.Read(buf, binary.LittleEndian, &index); err != nil {
		log.Print(err)
		return
	}

	// trim whitespace and convert to lowercase for key lookup:
	groupKey := calcGroupKey(group)
	clientGroup := findGroupOrCreate(groupKey)

	// create a key that represents the client from the received address:
	addr := message.ReceivedFrom
	clientKey := ClientKey{
		Port: addr.Port,
		Zone: addr.Zone,
	}
	copy(clientKey.IP[:], addr.IP)

	client, ci := findClientOrCreate(clientGroup, clientKey, addr, group, groupKey)

	// record number of bytes received:
	networkMetrics.ReceivedBytes(len(message.Envelope), kind.String(), clientGroup, client)

	switch kind {
	case RequestIndex:
		// client requests its own client index, no need to broadcast to other clients:

		// construct message:
		rsp := make02Packet(groupBuf, kind)

		// emit client index:
		index := uint16(ci)
		binary.Write(rsp, binary.LittleEndian, &index)

		// send message back to client:
		rspBytes := rsp.Bytes()
		_, fatalErr = conn.WriteToUDP(rspBytes, &client.UDPAddr)
		if fatalErr != nil {
			return
		}
		networkMetrics.SentBytes(len(rspBytes), kind.String(), clientGroup, client)
		rsp = nil

		break
	case RequestSnapshot:
		// client requests a full-state snapshot resync
		// find the "host" client (lowest active index) in the group:
		hostIdx := -1
		hostAddr := &client.UDPAddr
		for i := range clientGroup.Clients {
			c := &clientGroup.Clients[i]
			if !c.IsAlive {
				continue
			}
			if c == client {
				continue
			}
			if hostIdx == -1 || int(c.Index) < hostIdx {
				hostIdx = int(c.Index)
				hostAddr = &c.UDPAddr
			}
		}
		if hostIdx == -1 {
			// no other players to request from
			break
		}

		// forward the request to the host, include requesting client's index:
		rsp := make02Packet(groupBuf, kind)
		requestingIdx := uint16(ci)
		binary.Write(rsp, binary.LittleEndian, &requestingIdx)
		payload := buf.Bytes()
		rsp.Write(payload)

		rspBytes := rsp.Bytes()
		_, fatalErr = conn.WriteToUDP(rspBytes, hostAddr)
		if fatalErr != nil {
			return
		}
		networkMetrics.SentBytes(len(rspBytes), kind.String(), clientGroup, client)
		rsp = nil
		break
	case SnapshotData:
		// host is responding with snapshot data to a specific client
		// target client index is the first uint16 in the payload:
		var targetIdx uint16
		if err := binary.Read(buf, binary.LittleEndian, &targetIdx); err != nil {
			log.Print(err)
			break
		}
		payload := buf.Bytes()

		// find target client and unicast the snapshot:
		for i := range clientGroup.Clients {
			c := &clientGroup.Clients[i]
			if !c.IsAlive || int(c.Index) != int(targetIdx) {
				continue
			}
			rsp := make02Packet(groupBuf, kind)
			index := uint16(ci)
			binary.Write(rsp, binary.LittleEndian, &index)
			rsp.Write(payload)

			rspBytes := rsp.Bytes()
			_, fatalErr = conn.WriteToUDP(rspBytes, &c.UDPAddr)
			if fatalErr != nil {
				return
			}
			networkMetrics.SentBytes(len(rspBytes), kind.String(), clientGroup, client)
			rsp = nil
			break
		}
		break
	case Broadcast:
		// broadcast message received to all other clients:
		payload := buf.Bytes()
		for i := range clientGroup.Clients {
			c := &clientGroup.Clients[i]
			if !c.IsAlive {
				continue
			}
			if c == client {
				continue
			}

			responseKind := kind

			// construct message:
			rsp := make02Packet(groupBuf, responseKind)
			index := uint16(ci)
			binary.Write(rsp, binary.LittleEndian, &index)
			// write the payload:
			rsp.Write(payload)

			// send message to this client:
			rspBytes := rsp.Bytes()
			_, fatalErr = conn.WriteToUDP(rspBytes, &c.UDPAddr)
			if fatalErr != nil {
				return
			}
			networkMetrics.SentBytes(len(rspBytes), kind.String(), clientGroup, client)
			rsp = nil
			//log.Printf("[group %s] (%v) sent message to (%v)\n", groupKey, client, other)
		}
		break
	case BroadcastToSector:
		// broadcast message received to all other clients in the same sector:
		var sector uint32
		if err := binary.Read(buf, binary.LittleEndian, &sector); err != nil {
			log.Print(err)
			return
		}

		// join this client to the sector they're broadcasting to:
		client.Sector = sector

		// broadcast message:
		payload := buf.Bytes()
		for i := range clientGroup.Clients {
			c := &clientGroup.Clients[i]
			if !c.IsAlive {
				continue
			}
			if c.Sector != sector {
				continue
			}
			if c == client {
				continue
			}

			// construct message:
			rsp := make02Packet(groupBuf, Broadcast)
			index := uint16(ci)
			binary.Write(rsp, binary.LittleEndian, &index)

			// write the payload:
			rsp.Write(payload)

			// send message to this client:
			rspBytes := rsp.Bytes()
			_, fatalErr = conn.WriteToUDP(rspBytes, &c.UDPAddr)
			if fatalErr != nil {
				return
			}
			networkMetrics.SentBytes(len(rspBytes), kind.String(), clientGroup, client)
			rsp = nil
			//log.Printf("[group %s] (%v) sent message to (%v)\n", groupKey, client, other)
		}
		break
	}

	return
}
