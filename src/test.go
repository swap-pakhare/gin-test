package main

import (
	"fmt"
	"log"
	"os"
)

func checkError(err error)  {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	tmpFile := "text.log"
	f, err := os.OpenFile(tmpFile, os.O_RDWR, os.ModePerm)
	defer f.Close()
	checkError(err)

	off := int64(1)
	buf := make([]byte, off)
	stat, err := os.Stat(tmpFile)
	fmt.Println("Size : ", stat.Size())
	start := stat.Size() - off

	str := "This\n]"
	tmp := []byte(str)
	_, err = f.WriteAt(tmp, start)
	checkError(err)

	_, err = f.ReadAt(buf, start+int64(len(str)-1))
	checkError(err)
	fmt.Println(string(buf))
}
