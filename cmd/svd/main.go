//-----------------------------------------------------------------------------
/*

Test code for SVD reading.

*/
//-----------------------------------------------------------------------------

package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
)

//-----------------------------------------------------------------------------

func svdParse(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	defer file.Close()
	defer gz.Close()

	scanner := bufio.NewScanner(gz)

	_ = scanner

	return nil
}

//-----------------------------------------------------------------------------

const filePath = "./vendor/gigadevice/svd/gd32vf103.svd.gz"

func main() {

	err := svdParse(filePath)

	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

//-----------------------------------------------------------------------------
