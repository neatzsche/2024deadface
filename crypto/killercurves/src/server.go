package main

import (
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
)

const (
	PORT = ":4444"
	//FLAG = "FLAG{1nv@l1d-cUrv3-atT@ck-on-ECDH-W}"
)

var (
	a    = big.NewInt(-3)
	b    = big.NewInt(0)
	p    = big.NewInt(0)
	FLAG = os.Getenv("KILLER_CURVE_PRIVATE_KEY")
)

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Error creating server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	curve := elliptic.P256()
	b = curve.Params().B
	p = curve.Params().P
	serverPrivKey, serverPubX, serverPubY := generatePrivateKey(curve)

	serverPubKey := fmt.Sprintf("%s:%s", serverPubX.Text(16), serverPubY.Text(16))
	conn.Write([]byte(serverPubKey + "\n"))
	fmt.Println("Server public key sent:", serverPubKey)

	buffer := make([]byte, 256)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading client data:", err)
		return
	}
	clientPubKeyHex := strings.TrimSpace(string(buffer[:n]))
	clientPubX, clientPubY, err := parsePublicKey(clientPubKeyHex)
	if err != nil {
		fmt.Println("Invalid client public key:", err)
		return
	}

	point := &Point{
		x: clientPubX,
		y: clientPubY,
	}

	sharedSecretX := ScalarMult(point, serverPrivKey).x

	conn.Write([]byte((sharedSecretX).String() + "\n"))
	fmt.Println("Shared secret derived and sent to client.")
}

func generatePrivateKey(curve elliptic.Curve) (*big.Int, *big.Int, *big.Int) {
	serverPrivKey := new(big.Int)
	serverPrivKey.SetString(hex.EncodeToString([]byte(FLAG)), 16)

	serverPubX, serverPubY := curve.ScalarBaseMult(serverPrivKey.Bytes())
	return serverPrivKey, serverPubX, serverPubY
}

func parsePublicKey(pubKeyHex string) (*big.Int, *big.Int, error) {
	parts := strings.Split(pubKeyHex, ":")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid public key format")
	}

	clientPubX := new(big.Int)
	clientPubY := new(big.Int)

	clientPubX.SetString(parts[0], 16)
	clientPubY.SetString(parts[1], 16)

	return clientPubX, clientPubY, nil
}

type Point struct {
	x, y *big.Int
}

func modInverse(k, p *big.Int) *big.Int {
	return new(big.Int).ModInverse(k, p)
}

func modSubtract(a, b, p *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Sub(a, b), p)
}

func modAdd(a, b, p *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Add(a, b), p)
}

func modMultiply(a, b, p *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Mul(a, b), p)
}

func modSquare(a, p *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Mul(a, a), p)
}

func PointAddition(P, Q *Point) *Point {
	if P.x == nil && P.y == nil {
		return Q
	}
	if Q.x == nil && Q.y == nil {
		return P
	}

	if P.x.Cmp(Q.x) == 0 && P.y.Cmp(Q.y) == 0 {
		return PointDoubling(P)
	}

	lambda := modMultiply(modSubtract(Q.y, P.y, p), modInverse(modSubtract(Q.x, P.x, p), p), p)

	x3 := modSubtract(modSubtract(modSquare(lambda, p), P.x, p), Q.x, p)

	y3 := modSubtract(modMultiply(lambda, modSubtract(P.x, x3, p), p), P.y, p)

	return &Point{x: x3, y: y3}
}

func PointDoubling(P *Point) *Point {
	if P.x == nil && P.y == nil {
		return P
	}

	lambda := modMultiply(modAdd(modMultiply(big.NewInt(3), modSquare(P.x, p), p), a, p),
		modInverse(modMultiply(big.NewInt(2), P.y, p), p), p)

	x3 := modSubtract(modSquare(lambda, p), modMultiply(big.NewInt(2), P.x, p), p)

	y3 := modSubtract(modMultiply(lambda, modSubtract(P.x, x3, p), p), P.y, p)

	return &Point{x: x3, y: y3}
}

func ScalarMult(P *Point, k *big.Int) *Point {
	result := &Point{nil, nil}
	addend := P

	for i := k.BitLen() - 1; i >= 0; i-- {
		result = PointDoubling(result)

		if k.Bit(i) == 1 {
			result = PointAddition(result, addend)
		}
	}

	return result
}
