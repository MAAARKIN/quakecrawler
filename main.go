package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/maaarkin/quakecrawler/internal/handler"
)

func main() {

	fmt.Println("init read file")
	file, err := os.Open("quake.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Println("start scan")
	scanner := bufio.NewScanner(file)
	handler := handler.NewQuakeHandler()
	report, killByMeansReport := handler.Run(scanner)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	fmt.Println("------------ Report ------------")
	enc.Encode(report)
	fmt.Println("------------ Report KillByMeans ------------")
	enc.Encode(killByMeansReport)
}
