package server

import (
	"bytes"
	"context"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/quhive/qumine-ingress/proto"
	"github.com/sirupsen/logrus"
)

const (
	handshakeTimeout = 5 * time.Second
)

var (
	noDeadline time.Time
)

// Server represents the server
type Server struct {
	// Status is the current status of the server.
	Status string
	// Addr is the current address of the server.
	Addr  string
	state proto.State
}

// NewServer creates a new server instance with the given host and port
func NewServer(host string, port int) *Server {
	return &Server{
		Addr: net.JoinHostPort(host, strconv.Itoa(port)),
	}
}

// Start the server
func (server *Server) Start(context context.Context) {
	logrus.WithFields(logrus.Fields{
		"addr": server.Addr,
	}).Info("starting server...")

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		logrus.WithError(err).Fatal("server failed to start")
	}

	server.acceptConnections(context, listener)
}

func (server *Server) acceptConnections(context context.Context, listener net.Listener) {
	defer listener.Close()
	server.Status = "up"

	for {
		select {
		case <-context.Done():
			server.Status = "down"
			return
		default:
			connection, err := listener.Accept()
			if err != nil {
				logrus.WithError(err).Error("connection accept failed")
			} else {
				go server.handleConnection(context, connection)
			}
		}
	}
}

func (server *Server) handleConnection(context context.Context, client net.Conn) {
	defer client.Close()
	defer logrus.WithFields(logrus.Fields{
		"client": client.RemoteAddr(),
	}).Info("closed client connection")

	logrus.WithFields(logrus.Fields{
		"client": client.RemoteAddr(),
	}).Info("inbound client connection")

	buffer := new(bytes.Buffer)
	reader := io.TeeReader(client, buffer)

	if err := client.SetReadDeadline(time.Now().Add(handshakeTimeout)); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"client": client.RemoteAddr(),
		}).Error("setting deadline failed")
		return
	}
	packet, err := proto.ReadPacket(reader, client.RemoteAddr(), server.state)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"client": client.RemoteAddr(),
		}).Error("reading packet failed")
		return
	}
	logrus.WithFields(logrus.Fields{
		"client":       client.RemoteAddr(),
		"packetLength": packet.Length,
		"packetID":     packet.PacketID,
	}).Debug("received packet")

	if packet.PacketID == proto.HandshakeID {
		handshake, err := proto.ReadHandshake(packet.Data)
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"client": client.RemoteAddr(),
			}).Error("decoding handshake packet failed")
			return
		}
		logrus.WithFields(logrus.Fields{
			"client":    client.RemoteAddr(),
			"handshake": handshake,
		}).Debug("decoded handshake")

		hostname := handshake.ServerAddress
		server.findAndConnectBackend(context, client, buffer, hostname, "handshake")
	} else if packet.PacketID == proto.LegacyServerListPingID {
		handshake, ok := packet.Data.(*proto.LegacyServerListPing)
		if !ok {
			logrus.WithError(err).WithFields(logrus.Fields{
				"client": client.RemoteAddr(),
			}).Error("decoding legacyServerListPing packet failed")
			return
		}
		logrus.WithFields(logrus.Fields{
			"client":    client.RemoteAddr(),
			"handshake": handshake.ServerAddress,
		}).Debug("decoded legacyServerListPing")

		hostname := handshake.ServerAddress
		server.findAndConnectBackend(context, client, buffer, hostname, "legacyServerListPing")
	} else {
		logrus.WithFields(logrus.Fields{
			"client":   client.RemoteAddr(),
			"packetID": packet.PacketID,
		}).Error("received unexpected packet, expected handshake or legacyServerListPing")
		return
	}
}

func (server *Server) findAndConnectBackend(context context.Context, client net.Conn, preReadContent io.Reader, hostname string, packet string) {
	route, err := ReadRoute(hostname)
	if err != nil {
		logrus.WithError(err).Warn("no matching route found")
		metricsErrorsTotal.With(prometheus.Labels{"error": "no-route"}).Inc()
		return
	}
	logrus.WithFields(logrus.Fields{
		"client": client.RemoteAddr(),
		"route":  route,
	}).Debug("found matching route")

	upstream, err := net.Dial("tcp", route)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"client": client.RemoteAddr(),
			"route":  route,
		}).Error("connecting to upstream failed")
		metricsErrorsTotal.With(prometheus.Labels{"error": "upstream-unavailable"}).Inc()
		return
	}
	metricsConnectionsTotal.With(prometheus.Labels{"route": route}).Inc()
	logrus.WithFields(logrus.Fields{
		"client":   client.RemoteAddr(),
		"upstream": upstream.RemoteAddr(),
	}).Info("connected to upstream")

	amount, err := io.Copy(upstream, preReadContent)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"client":   client.RemoteAddr(),
			"upstream": upstream.RemoteAddr(),
		}).Error("failed to relay packet to upstream")
		return
	}
	logrus.WithFields(logrus.Fields{
		"client":   client.RemoteAddr(),
		"upstream": upstream.RemoteAddr(),
		"amount":   amount,
	}).Debugf("relayed %s to upstream", packet)

	if err = client.SetReadDeadline(noDeadline); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"client":   client.RemoteAddr(),
			"upstream": upstream.RemoteAddr(),
		}).Error("clearing deadline failed")
		return
	}
	server.relayConnections(context, route, client, upstream)
	return
}

func (server *Server) relayConnections(context context.Context, route string, client net.Conn, upstream net.Conn) {
	defer upstream.Close()
	defer logrus.WithFields(logrus.Fields{
		"client":   client.RemoteAddr(),
		"upstream": upstream.RemoteAddr(),
	}).Info("closed upstream connection")

	errors := make(chan error, 2)
	go server.relay(client, upstream, errors, "upstream", route, client, upstream)
	go server.relay(upstream, client, errors, "downstream", route, client, upstream)
	logrus.WithFields(logrus.Fields{
		"client":   client.RemoteAddr(),
		"upstream": upstream.RemoteAddr(),
	}).Debug("relayed connection to upstream")

	select {
	case err := <-errors:
		if err != io.EOF {
			logrus.WithError(err).WithFields(logrus.Fields{
				"client":   client.RemoteAddr(),
				"upstream": upstream.RemoteAddr(),
			}).Error("clearing deadline failed")
		}

	case <-context.Done():
		return
	}
}

func (server *Server) relay(incoming io.Reader, outgoing io.Writer, errors chan<- error, direction string, route string, client net.Conn, upstream net.Conn) {
	amount, err := io.Copy(outgoing, incoming)
	metricsBytesTotal.With(prometheus.Labels{"direction": direction, "route": route}).Add(float64(amount))
	logrus.WithFields(logrus.Fields{
		"client":    client.RemoteAddr(),
		"upstream":  upstream.RemoteAddr(),
		"direction": direction,
		"amount":    amount,
	}).Debug("stopped connection relay")

	if err != nil {
		errors <- err
	} else {
		errors <- io.EOF
	}
}
