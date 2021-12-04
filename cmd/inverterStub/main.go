package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
)

var count = 0

func handleConnection(c net.Conn) {
	// fmt.Print(".")
	defer c.Close()
	for {
		// netData, err := bufio.NewReader(c).ReadString(0x15)
		netData, err := bufio.NewReader(c).ReadBytes(0x15)

		// netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		// temp := strings.TrimSpace(string(netData))
		// if temp == "STOP" {
		// 	break
		// }
		// fmt.Println(temp)
		fmt.Printf("Data reci: %X\n", netData)

		Message2 := "A5 33 00 10 15 00 02 FF D9 F9 1E 02 01 6F 3E D1 01 07 00 00 00 00 00 00 00 01 03 20 00 00 00 00 00 00 00 00 00 00 00 00 1B 1D 00 00 1F F3 00 62 01 55 00 16 00 22 05 25 05 33 05 26 54 D5 E6 15"
		Message1 := "A5 63 00 10 15 00 A5 FF D9 F9 1E 02 01 6E 1A D0 01 44 2C 00 00 E2 48 CF 5F 01 03 50 00 02 00 00 00 00 00 00 00 00 00 00 07 95 00 68 07 5F 00 39 00 14 00 0A 00 1D 00 00 13 88 09 32 00 8B 00 00 00 00 00 00 00 00 00 00 1B 18 00 00 1F E2 00 1A 00 92 00 18 00 25 0E 90 07 95 07 5F 00 3C 00 00 00 01 00 00 06 AC 2E E0 2E E0 00 07 DE C3 45 15"
		// dummy := []byte(netData)
		var listje []string
		if netData[len(netData)-2] == 0xA2 {
			listje = strings.Split(Message1, " ")
		} else {
			listje = strings.Split(Message2, " ")
		}

		data := []byte{}
		for _, v := range listje {
			d, err := hex.DecodeString(v)
			if err != nil {
				panic(err)
			}
			data = append(data, d...)
		}
		fmt.Printf("Data send: % X\n", data)

		c.Write([]byte(data))

	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
		count++
	}
}
