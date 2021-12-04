package inverter

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/sigurn/crc16"
)

const inverterSerial = 519690751

var (
	verbose            *bool
	inverterStringData *bool
)

// readJSONfile convert json file data into ConversionData struct
func readJSONfile(filename string) ConversionData {
	data := ConversionData{}

	file, err := ioutil.ReadFile(filename)
	CheckErr(err)

	err = json.Unmarshal([]byte(file), &data)
	CheckErr(err)
	return data
}

func SendAndProcess(register_start, register_end int, con net.Conn) []OutputData {
	msg := createMsg(register_start, register_end)
	if *verbose {
		fmt.Printf("Data send: %X\n", msg)
	}

	data := readJSONfile("data.json")

	_, err := con.Write([]byte(msg))
	CheckErr(err)

	reply := make([]byte, 128)

	_, err = con.Read(reply)
	CheckErr(err)

	if *verbose {
		fmt.Printf("Date recv: %X\n", reply)
	}
	return printAllData(reply, register_start, register_end, data)
}

// msg := "A5 17 00 10 45 00 00 ff D9 F9 1E 02 00 00 00 00 00 00 00 00 00 00 00 00 00 00 01 04 20 00 00 0E 7A 0E 00 15 <- HWdata inverter
//         A5 17 00 10 45 00 00 FF D9 F9 1E 02 00 00 00 00 00 00 00 00 00 00 00 00 00 00 01 03 00 00 00 28 45 D4 A2 15 <- inverter data
//         A5 17 00 10 45 00 00 FF D9 F9 1E 02 00 00 00 00 00 00 00 00 00 00 00 00 00 00 01 03 00 10 00 10 45 C3 89 15 <- string
//         -\ ----\ ----\ ----\ ----------\ -------------------------------------------\ ----\ ----------\ ----\ -\ -\
//        start                 serial      datafield                                          pos S to E| crc |cs|endcode
func createMsg(register_start, register_end int) []byte {
	startCode := "A5"                    // # start A5
	length := "1700"                     // # datalength
	controlcode := "1045"                // # controlCode
	serial := "0000"                     // # serial
	endCode, _ := hex.DecodeString("15") // # end 15

	hexInverterSerial := strconv.FormatInt(int64(inverterSerial), 16)
	inverter_sn2, err := hex.DecodeString(hexInverterSerial[6:8] + hexInverterSerial[4:6] + hexInverterSerial[2:4] + hexInverterSerial[0:2])
	CheckErr(err)

	datafield, _ := hex.DecodeString("020000000000000000000000000000")

	// register_start2 int = 0x0105
	// register_end2   int = 0x0114
	// 00 10 00 10
	// registerhw_start int = 0x2000
	// registerhw_end   int = 0x200D
	// 20 00 00 0E

	// 	# Modbus request begin
	// pos_ini=str(hex_zfill(pini)[2:])
	// pos_fin=str(hex_zfill(pfin-pini+1)[2:])
	// businessfield= binascii.unhexlify('0104' + pos_ini + pos_fin) # Modbus data to count crc

	pos_ini := Zfill(strconv.FormatInt(int64(register_start), 16), 4)[1:3]
	pos_fin := Zfill(strconv.FormatInt(int64(register_end-register_start+1), 16), 4)[2:4]

	businessfield, _ := hex.DecodeString("0103" + Zfill(pos_ini, 4) + Zfill(pos_fin, 4))

	a := Zfill(fmt.Sprintf("%X", crc16.Checksum(businessfield, crc16.MakeTable(crc16.CRC16_MODBUS))), 4)
	crc, _ := hex.DecodeString(a[2:4] + a[0:2]) // Swap the fields

	msg, _ := hex.DecodeString(startCode + length + controlcode + serial)
	msg = append(msg, inverter_sn2...)
	msg = append(msg, datafield...)
	msg = append(msg, businessfield...)
	msg = append(msg, crc...)
	msg = append(msg, checksumMsg(msg))
	msg = append(msg, endCode...)

	return msg
}

// checksumMsg calculates the sum of all bytes in array, last byte of result is checksum
func checksumMsg(m []byte) byte {
	var checksum uint64
	for i := 1; i < len(m); i++ {
		checksum += uint64(m[i])
	}
	return byte(checksum & 0xff)
}

// A56300101500D4FFD9F91E02013779D101653C0000A7EED15F01035000020000...
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// 0                   1                   2                   3

func printAllData(reply []byte, register_start, register_end int, data ConversionData) []OutputData {
	returnValue := []OutputData{}
	item := OutputData{}
	for i := 0; i < (register_end-register_start)*2; i++ {
		for x := range data {
			for y := range data[x].Items {
				for z := range data[x].Items[y].Registers {
					if data[x].Items[y].Registers[z] == fmt.Sprintf("0x%04X", register_start+i) {
						// OptionRanges can be contain multiple adresses
						if len(data[x].Items[y].OptionRanges) == 0 {
							if *verbose {
								fmt.Printf("0x%04X - %v: %.2f%v\n",
									register_start+i,
									data[x].Items[y].TitleEN,
									float64(getValue(i*2, reply))*data[x].Items[y].Ratio,
									data[x].Items[y].Unit)
							}

							item.Key = data[x].Items[y].TitleEN
							item.Value = fmt.Sprintf("%f", float64(getValue(i*2, reply))*data[x].Items[y].Ratio)

						} else {
							if *verbose {
								fmt.Printf("0x%04X - %v: %v\n",
									register_start+i,
									data[x].Items[y].TitleEN,
									data[x].Items[y].OptionRanges[getValue(i*2, reply)].ValueEN)
							}

							item.Key = data[x].Items[y].TitleEN
							item.Value = data[x].Items[y].OptionRanges[getValue(i*2, reply)].ValueEN
						}
						returnValue = append(returnValue, item)
					}
				}
			}
		}
	}
	return returnValue
}

// Zfill fills the provide string with the max of (overall) zeros in front of string
//
// example: Zfill("9",4) result in "0009"
func Zfill(s string, overall int) string {
	l := overall - len(s)
	return strings.Repeat("0", l) + s
}

// getValue get part of []byte data. Data always start from 28th byte
//
// s: is start position
// As values alway exist of two bytes this is set fixed
func getValue(s int, data []byte) int64 {
	r := ""
	s += 28
	if len(data) >= s+2 {
		for i := s; i < s+2; i++ {
			r = r + fmt.Sprintf("%v", strconv.FormatInt(int64(data[i]), 16))
		}
	}
	output, _ := strconv.ParseInt(r, 16, 64)
	return output
}

// CheckErr checks error result. In case of an error print error and stop
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
