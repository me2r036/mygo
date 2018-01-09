package main

import "fmt"
import "os"
import "crypto/md5"
import "io"
import "encoding/hex"

func getKey(email, expireDate string) string {
	k := "EX1," + email + "," + expireDate
	h := md5.New()
	io.WriteString(h, k+","+"356B4B5C")
	return k + "," + hex.EncodeToString(h.Sum(nil))
}

func main() {
	fmt.Println(getKey(os.Args[1], os.Args[2]))
}
