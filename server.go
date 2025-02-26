package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unicode"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

type Configure struct {
	ArrSize     int `json:"array_size"`
	MaxRoutines int `json:"maximum_routines"`
}

var (
	rutineGoActive int32
	config         Configure
)

func main() {

	configFile, err := os.Open("configurations.json")
	if err != nil {
		fmt.Println("Eroare deschidere configurations.json", err)
		return
	}

	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Eroare decode configurations.json", err)
		return
	}

	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		for atomic.LoadInt32(&rutineGoActive) >= int32(config.MaxRoutines) {
			fmt.Println("Numar maxim de rutine go atins..In astepare...")
			time.Sleep(2 * time.Second)
		}

		atomic.AddInt32(&rutineGoActive, 1)

		go request(conn)
	}

}

func request(conn net.Conn) {
	defer atomic.AddInt32(&rutineGoActive, -1)
	defer conn.Close()

	read := bufio.NewReader(conn)
	connectMessage, err := read.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(connectMessage)

	problemNr := make([]byte, config.ArrSize)
	n, _ := conn.Read(problemNr)
	problem := strings.TrimSpace(string(problemNr[:n]))

	content := make([]byte, config.ArrSize)
	p, err := conn.Read(content)
	if err != nil {
		println("Eroare citire content fisier", err.Error())
		os.Exit(1)
	}
	fileContent := string(content[:p])
	fmt.Printf("  Serverul a primit requestul\n  ")
	fmt.Printf("  Serverul proceseaza datele\n  ")

	switch problem {
	case "1":
		response := cerinta1(fileContent)
		fmt.Printf("Server trimite raspunsul: %s catre client\n", response)
		_, err = conn.Write([]byte(response))

	case "2":
		response := cerinta2(fileContent)
		fmt.Printf("Server trimite raspunsul: %s catre client\n", response)
		_, err = conn.Write([]byte(response))
	case "3":
		response := cerinta3(fileContent)
		fmt.Printf("Server trimite raspunsul: %s catre client\n", response)
		_, err = conn.Write([]byte(response))
	case "4":
		response := cerinta4(fileContent)
		fmt.Printf("Server trimite raspunsul: %s catre client\n", response)
		_, err = conn.Write([]byte(response))
	case "5":
		response := cerinta5(fileContent)
		fmt.Printf("Server trimite raspunsul: %s catre client\n", response)
		_, err = conn.Write([]byte(response))
	case "7":
		response := cerinta7(fileContent)
		fmt.Printf("Server trimite raspunsul: %s catre client\n", response)
		_, err = conn.Write([]byte(response))
	case "8":
		response := cerinta8(fileContent)
		fmt.Printf("Server trimite raspunsul: %s catre client\n", response)
		_, err = conn.Write([]byte(response))
	case "12":
		response := cerinta12(fileContent)
		fmt.Printf("Server trimite raspunsul: %s catre client\n", response)
		_, err = conn.Write([]byte(response))
	default:
		fmt.Println("Nu exista problema rezolvata pentru acest numar.")
	}

	conn.Close()
}

func cerinta1(file string) string {
	array := strings.Split(file, ", ")
	minLength := len(array[0])
	output := make([]string, minLength)
	for i := 0; i < minLength; i++ {
		var result string
		for _, word := range array {
			result += string(word[i])
		}
		output[i] = result
	}
	outputStr := strings.Join(output, ", ")
	return outputStr

}

func cerinta2(file string) string {
	array := strings.Split(file, ", ")
	count := 0

	patratPerfect := func(nr int) bool {
		sqrt := int(math.Sqrt(float64(nr)))
		return sqrt*sqrt == nr
	}

	for _, cuvant := range array {
		var number string

		for _, char := range cuvant {
			if unicode.IsDigit(char) {
				number += string(char)
			}
		}

		if number != "" {
			nr, _ := strconv.Atoi(number)
			if patratPerfect(nr) {
				count++
			}
		}
	}
	return fmt.Sprintf("<%d patrate perfecte>", count)
}

func cerinta3(file string) string {
	numbers := strings.Split(file, ", ")
	sum := 0

	for _, nr := range numbers {
		nr = strings.TrimSpace(nr)
		reversed := ""
		for i := len(nr) - 1; i >= 0; i-- {
			reversed += string(nr[i])
		}
		num, _ := strconv.Atoi(reversed)
		sum += num
	}

	return fmt.Sprintf("<Suma %d>", sum)
}

func cerinta4(file string) string {
	sir := strings.Split(file, ",")
	a, _ := strconv.Atoi(strings.TrimSpace(sir[0]))
	b, _ := strconv.Atoi(strings.TrimSpace(sir[1]))
	// n, _ := strconv.Atoi(strings.TrimSpace(sir[2]))
	numere := sir[3:]

	sumaCifre := func(nr int) int {
		suma := 0
		for nr != 0 {
			suma += nr % 10
			nr /= 10
		}
		return suma
	}

	var numereInInterval []int
	for _, nrStr := range numere {
		nr, err := strconv.Atoi(strings.TrimSpace(nrStr))
		if err != nil {
			continue
		}
		if sumaCifre(nr) >= a && sumaCifre(nr) <= b {
			numereInInterval = append(numereInInterval, nr)
		}
	}

	if len(numereInInterval) == 0 {
		return fmt.Sprintf("Nu exista numere care se afla in intervalul [%d, %d]", a, b)
	}

	suma := 0
	for _, nr := range numereInInterval {
		suma += nr
	}
	medie := float64(suma) / float64(len(numereInInterval))

	return fmt.Sprintf("<Media %.2f>", medie)
}

func cerinta5(file string) string {
	sir := strings.Split(file, ",")
	var numereConvertite []int

	binar := func(numar string) bool {
		for _, cifra := range numar {
			if cifra != '0' && cifra != '1' {
				return false
			}
		}
		return true
	}

	conversie := func(binar string) int {
		suma := 0
		power := len(binar) - 1
		for _, cifra := range binar {
			if cifra == '1' {
				suma += int(math.Pow(2, float64(power)))
			}
			power--
		}
		return suma
	}

	for _, element := range sir {
		element = strings.TrimSpace(element)
		if binar(element) {
			numarBaza10 := conversie(element)
			numereConvertite = append(numereConvertite, numarBaza10)
		}
	}

	if len(numereConvertite) == 0 {
		return "Nu exista numere binare."
	}

	return fmt.Sprintf("<Numerele convertite sunt %v>", numereConvertite)
}

func cerinta7(file string) string {
	decodedText := strings.Builder{}
	n := len(file)
	i := 0
	for i < n {
		numStart := i
		for i < n && file[i] >= '0' && file[i] <= '9' {
			i++
		}
		if numStart == i {
			return "Textul trebuie sa inceapa cu o cifra\n"
		}

		count, _ := strconv.Atoi(file[numStart:i])
		if i >= n {
			return "Lipsa caracter dupa numar\n"
		}
		character := file[i]
		i++
		decodedText.WriteString(strings.Repeat(string(character), count))
	}

	return decodedText.String()
}

func cerinta8(file string) string {
	sir := strings.Split(file, ", ")
	count := 0

	estePrim := func(nr string) bool {
		numar, _ := strconv.Atoi(strings.TrimSpace(nr))

		for i := 2; i*i <= numar; i++ {
			if numar%i == 0 {
				return false
			}
		}
		return true
	}

	for _, nr := range sir {
		if estePrim(nr) {
			count += len(nr)
		}
	}
	return fmt.Sprintf("<Total cifre %d>", count)
}

func cerinta12(file string) string {
	sir := strings.Split(file, ", ")
	suma := 0

	for _, nr := range sir {
		nr = strings.TrimSpace(nr)
		primaC := nr[0]
		numar, _ := strconv.Atoi(string(primaC) + nr)
		suma += numar
	}
	return fmt.Sprintf("<Suma este %d>", suma)
}
