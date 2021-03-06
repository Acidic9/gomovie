package gomovie

import (
	"fmt"
	"encoding/hex"
)

var arrChrs []string
var reversegetFChars map[string]int
var getFStr string
var getFCount int

func DecryptPutlocker(str string) string {
	arrChrs = []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","0","1","2","3","4","5","6","7","8","9","+","/"}
	reversegetFChars = make(map[string]int, len(arrChrs))
	for i, char := range arrChrs {
		reversegetFChars[char] = i
	}

	return doit(str)
}

func ntos(e int) string {
	h := fmt.Sprintf("%x", e)
	if len(h) == 1 {
		h = "0" + h
	}
	return unescape(h)
}

func readReversegetF() int {
	if len(getFStr) == 0 {
		return -1
	}
	for true {
		if getFCount >= len(getFStr) {
			return -1
		}
		e := getFStr[getFCount:getFCount+1]
		getFCount++
		if reversegetFChars[e] > 0 {
			return reversegetFChars[e]
		}
		if e == "A" {
			return 0
		}
	}
	return -1
}

func setgetFStr(e string) {
	getFStr = e
	getFCount = 0
}

func getF(e string) string {
	setgetFStr(e)
	var t string
	n := make([]int, 4)
	var r bool
	for !r {
		n[0] = readReversegetF()
		n[1] = readReversegetF()
		if n[0] == -1 && n[1] == -1 {
			break
		}
		n[2] = readReversegetF()
		n[3] = readReversegetF()

		t += ntos(n[0] << 2 & 255 | n[1] >> 4)
		if n[2] != -1 {
			t += ntos(n[1] << 4 & 255 | n[2] >> 2)
			if n[3] != -1 {
				t += ntos(n[2] << 6 & 255 | n[3])
			} else {
				r = true
			}
		} else {
			r = true
		}
	}
	return t
}

func doit(e string) string {
	return getF(getF(e))
}

func unescape(hexi string) string {
	byte, _ := hex.DecodeString(hexi)
	return string(byte)
}