package main

import (
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"strings"
)

const (
	PORT = ":4444"
)

var (
	a = big.NewInt(-3)
	b = big.NewInt(0)
	p = big.NewInt(0)
	//FLAG = os.Getenv("KILLER_CURVE_PRIVATE_KEY")
	FLAG = "flag{1nv@l1d-cUrv3}"
	//FLAG = "flag{bad-ECC}"
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

	sharedSecretX := scalarMult(point, serverPrivKey).x

	conn.Write([]byte((sharedSecretX).String() + "\n"))
	fmt.Println("Shared secret derived and sent to client.")
}

func generatePrivateKey(curve elliptic.Curve) (*big.Int, *big.Int, *big.Int) {
	serverPrivKey := new(big.Int)
	serverPrivKey.SetString(hex.EncodeToString([]byte(FLAG)), 16)
	//serverPrivKey.SetInt64(0xf1a6ecc101)
	fmt.Println(serverPrivKey)

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

func isInfinity(p *Point) bool {
	return p.x == nil && p.y == nil
}

func pointAdd(P1, P2 *Point) *Point {
	if isInfinity(P1) {
		return &Point{new(big.Int).Set(P2.x), new(big.Int).Set(P2.y)}
	}
	if isInfinity(P2) {
		return &Point{new(big.Int).Set(P1.x), new(big.Int).Set(P1.y)}
	}

	if P1.x.Cmp(P2.x) == 0 && P1.y.Cmp(P2.y) == 0 {
		return pointDouble(P1)
	}

	lambda := new(big.Int).Sub(P2.y, P1.y)
	denom := new(big.Int).Sub(P2.x, P1.x)
	denom.ModInverse(denom, p)
	lambda.Mul(lambda, denom)
	lambda.Mod(lambda, p)

	x3 := new(big.Int).Mul(lambda, lambda)
	x3.Sub(x3, P1.x)
	x3.Sub(x3, P2.x)
	x3.Mod(x3, p)

	y3 := new(big.Int).Sub(P1.x, x3)
	y3.Mul(lambda, y3)
	y3.Sub(y3, P1.y)
	y3.Mod(y3, p)

	return &Point{x3, y3}
}

func pointDouble(P1 *Point) *Point {
	lambda := new(big.Int).Mul(big.NewInt(3), new(big.Int).Mul(P1.x, P1.x))
	lambda.Add(lambda, a)
	denom := new(big.Int).Mul(big.NewInt(2), P1.y)
	denom.ModInverse(denom, p)
	lambda.Mul(lambda, denom)
	lambda.Mod(lambda, p)

	x3 := new(big.Int).Mul(lambda, lambda)
	x3.Sub(x3, new(big.Int).Mul(big.NewInt(2), P1.x))
	x3.Mod(x3, p)

	y3 := new(big.Int).Sub(P1.x, x3)
	y3.Mul(lambda, y3)
	y3.Sub(y3, P1.y)
	y3.Mod(y3, p)

	return &Point{x3, y3}
}

func scalarMult(P *Point, k *big.Int) *Point {
	result := &Point{nil, nil}

	for i := k.BitLen() - 1; i >= 0; i-- {
		if !isInfinity(result) {
			result = pointDouble(result)
		}

		if k.Bit(i) == 1 {
			if isInfinity(result) {
				result = &Point{new(big.Int).Set(P.x), new(big.Int).Set(P.y)}
			} else {
				result = pointAdd(result, P)
			}
		}
	}

	return result
}
