package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"strings"
)

const (
	defaultPort = "4444"
)

func main() {
	publicKeyArg := flag.String("pubkey", "", "Public key in hex format (x:y)")
	serverAddr := flag.String("server", "", "Server IP or DNS name")
	port := flag.String("port", defaultPort, "Port number (default is 4444)")
	filePath := flag.String("file", "", "Path to the file to be encrypted")

	flag.Parse()

	if *publicKeyArg == "" || *serverAddr == "" || *filePath == "" {
		fmt.Println("Usage: -pubkey <public key> -server <server IP/DNS> -file <file path> [-port <port>]")
		return
	}

	fmt.Println("Checking if the public key is on the secp256k1 curve...")

	curve := elliptic.P256()
	clientPubX, clientPubY, err := parsePublicKey(*publicKeyArg)
	if err != nil {
		fmt.Println("Invalid public key format:", err)
		return
	}

	if !curve.IsOnCurve(clientPubX, clientPubY) {
		fmt.Println("Public key is not on the curve.")
		return
	}

	fmt.Println("Public key is valid and on the curve.")

	fmt.Println("Generating random scalar and performing scalar multiplication...")

	clientPrivKey, err := rand.Int(rand.Reader, curve.Params().N)
	if err != nil {
		fmt.Println("Error generating random scalar:", err)
		return
	}

	clientSharedX, clientSharedY := curve.ScalarMult(clientPubX, clientPubY, clientPrivKey.Bytes())
	fmt.Printf("Randomized Public Key Computed")

	fmt.Println("Connecting to the server to exchange public keys...")

	conn, err := net.Dial("tcp", *serverAddr+":"+*port)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 256)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving server public key:", err)
		return
	}
	serverPubKeyHex := strings.TrimSpace(string(buffer[:n]))
	fmt.Println("Received server public key:", serverPubKeyHex)

	clientPubKey := fmt.Sprintf("%s:%s", clientSharedX.Text(16), clientSharedY.Text(16))
	fmt.Println("Sending client's randomized public key:", clientPubKey)
	conn.Write([]byte(clientPubKey + "\n"))

	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving shared secret X Value from server:", err)
		return
	}
	sharedSecretXHex := strings.TrimSpace(string(buffer[:n]))
	fmt.Println("Received shared secret X value from server:", sharedSecretXHex)

	fmt.Println("Encrypting the file with the shared secret...")

	sharedSecretX, _ := new(big.Int).SetString(sharedSecretXHex, 16)
	if sharedSecretX == nil {
		fmt.Println("Invalid shared secret X from server.")
		return
	}

	fileContent, err := ioutil.ReadFile(*filePath)
	if err != nil {
		fmt.Println("Error reading the file:", err)
		return
	}

	encryptedContent := xorEncrypt(fileContent, sharedSecretX.Bytes())

	encryptedFilePath := *filePath + ".enc"
	err = ioutil.WriteFile(encryptedFilePath, encryptedContent, 0644)
	if err != nil {
		fmt.Println("Error saving encrypted file:", err)
		return
	}

	fmt.Printf("File successfully encrypted. Saved as %s\n", encryptedFilePath)
}

func parsePublicKey(pubKeyHex string) (*big.Int, *big.Int, error) {
	parts := strings.Split(pubKeyHex, ":")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid public key format")
	}

	pubX := new(big.Int)
	pubY := new(big.Int)

	_, success := pubX.SetString(parts[0], 16)
	if !success {
		return nil, nil, fmt.Errorf("invalid X coordinate")
	}
	_, success = pubY.SetString(parts[1], 16)
	if !success {
		return nil, nil, fmt.Errorf("invalid Y coordinate")
	}

	return pubX, pubY, nil
}

func xorEncrypt(content []byte, key []byte) []byte {
	encrypted := make([]byte, len(content))
	keyLen := len(key)

	for i := range content {
		encrypted[i] = content[i] ^ key[i%keyLen]
	}

	return encrypted
}
