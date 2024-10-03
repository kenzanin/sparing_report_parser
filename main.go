package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "report_parser",
		Action: report_parser,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func find_col(header *[]string, find string) int {
	for i, e := range *header {
		if strings.Contains(strings.ToUpper(e), find) {
			return i
		}
	}
	return -1
}

type max_min struct {
	col int
	max float32
	min float32
}

func report_parser(cli *cli.Context) error {
	log.Print("Open file: ", cli.Args().Get(0))

	f, err := os.Open(cli.Args().Get(0))
	if err != nil {
		log.Fatal("error opening file: ", err)
	}
	defer f.Close()

	csv_reader := csv.NewReader(f)
	csv_reader.Comma = ';'
	header, err := csv_reader.Read()
	if err != nil {
		log.Fatal("error read header: ", err)
	}

	log.Print(header)

	col_waktu := find_col(&header, "WAKTU")
	if col_waktu == -1 {
		log.Fatal("error colloum waktu tidak ditemukan")
	}

	PH := max_min{col: find_col(&header, "PH"), min: 99999.0}
	COD := max_min{col: find_col(&header, "COD"), min: 99999.0}
	TSS := max_min{col: find_col(&header, "TSS"), min: 99999.0}
	NH3N := max_min{col: find_col(&header, "NH3N"), min: 99999.0}

	datas, err := csv_reader.ReadAll()
	if err != nil {
		log.Fatal("error decode csv: ", err)
	}

	var data_bolong []string
	var data_double []string

	first := int64(0)
	for index, el := range datas {

		p := &PH
		tmp, _ := strconv.ParseFloat(el[p.col], 32)
		if tmp > float64(p.max) {
			p.max = float32(tmp)
		}
		if tmp < float64(p.min) {
			p.min = float32(tmp)
		}

		p = &COD
		tmp, _ = strconv.ParseFloat(el[p.col], 32)
		if tmp > float64(p.max) {
			p.max = float32(tmp)
		}
		if tmp < float64(p.min) {
			p.min = float32(tmp)
		}

		p = &TSS
		tmp, _ = strconv.ParseFloat(el[p.col], 32)
		if tmp > float64(p.max) {
			p.max = float32(tmp)
		}
		if tmp < float64(p.min) {
			p.min = float32(tmp)
		}

		p = &NH3N
		tmp, _ = strconv.ParseFloat(el[p.col], 32)
		if tmp > float64(p.max) {
			p.max = float32(tmp)
		}
		if tmp < float64(p.min) {
			p.min = float32(tmp)
		}

		t, err := time.Parse("02/01/2006, 15:04", el[col_waktu])
		if err != nil {
			log.Print("error parse time: ", err)
		}
		if index == 0 {
			first = t.Unix()
			continue
		}

		current_time := t.Unix() - first
		first = t.Unix()
		if current_time == 0 {
			data_double = append(data_double, el[col_waktu])
		} else if current_time > 120 {
			data_bolong = append(data_bolong, fmt.Sprintf("%s len: %d", el[col_waktu], (current_time/120)-1))
		}

	}

	fmt.Println("data bolong total:", len(data_bolong))
	for _, e := range data_bolong {
		fmt.Println(e)
	}

	fmt.Println("data double total:", len(data_double))
	for _, e := range data_double {
		fmt.Println(e)
	}

	p := &PH
	fmt.Printf("data PH max:\t %.3f\t| min: %.3f\n", p.max, p.min)
	p = &COD
	fmt.Printf("data COD max:\t %.3f\t| min: %.3f\n", p.max, p.min)
	p = &TSS
	fmt.Printf("data TSS max:\t %.3f\t| min: %.3f\n", p.max, p.min)
	p = &NH3N
	fmt.Printf("data NH3N max:\t %.3f\t| min: %.3f\n", p.max, p.min)

	return nil
}
