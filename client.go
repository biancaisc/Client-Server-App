package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

type Config struct {
	ArrSize     int `json:"array_size"`
	MaxRoutines int `json:"maximum_routines"`
}

// Format nume fisier input:
// Ex: Problema nr 1 are: Client1_01.txt si Client1_02.txt
//	   Problema nr 2 are: Client2_01.txt si Client2_02.txt
//	   ...

func main() {

	configFile, err := os.Open("configurations.json")
	if err != nil {
		fmt.Println("Eroare deschidere configurations.json", err)
		return
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Eroare decode configurations.json", err)
		return
	}

	fmt.Print("	 Nume client: ")
	reader := bufio.NewReader(os.Stdin)
	clientName, _ := reader.ReadString('\n')
	clientName = strings.TrimSpace(clientName)

	adresaTCP, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		println("Eroare resolveTCPAddr ", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, adresaTCP)
	if err != nil {
		println("Eroare DialTCP:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	connectMessage := fmt.Sprintf("  Client %s conectat\n", clientName)
	fmt.Printf(connectMessage)
	_, err = conn.Write([]byte(connectMessage))
	if err != nil {
		fmt.Println("Eroare trimitere nume client:", err)
		os.Exit(1)
	}

	fmt.Printf("  Numar exercitiu: ")
	readEx := bufio.NewReader(os.Stdin)
	problemNr, err := readEx.ReadString('\n')
	if err != nil {
		fmt.Println("Eroare citire nr exercitiu: ", err)
		os.Exit(1)
	}
	problemNr = strings.TrimSpace(problemNr)
	_, err = conn.Write([]byte(problemNr + "\n"))
	if err != nil {
		fmt.Println("Eroare trimitere nr exercitiu:", err)
		os.Exit(1)
	}

	fmt.Printf("  Fisier input: ")
	readFile := bufio.NewReader(os.Stdin)
	fileName, err := readFile.ReadString('\n')
	fileName = strings.TrimSpace(fileName)
	if err != nil {
		fmt.Println("Eroare citire nume fisier:", err)
		os.Exit(1)
	}

	filePath := fmt.Sprintf("input/%s", fileName)

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Eroare citire continut fisier:", err)
		os.Exit(1)

	}
	if len(fileContent) > config.ArrSize {
		fmt.Printf("Eroare dimensiunea input file=%d bytes depaseste limita=%d\n", len(fileContent), config.ArrSize)
		os.Exit(1)
	}
	_, err = conn.Write(fileContent)
	if err != nil {
		println("Eroare scriere continut fisier:", err.Error())
		os.Exit(1)
	}

	fmt.Printf("  Clientul %s a facut request cu datele: %s\n", clientName, fileContent)
	_, err = conn.Write([]byte(fileContent))
	if err != nil {
		fmt.Printf("Eroare afisare date request:", err)
		os.Exit(1)
	}

	received := make([]byte, config.ArrSize)
	_, err = conn.Read(received)
	if err != nil {
		println("Eroare citire raspuns de la server:", err.Error())
		os.Exit(1)
	}

	fmt.Printf("  Client %s a primit raspunsul:%s ", clientName, string(received))
	conn.Close()
}
