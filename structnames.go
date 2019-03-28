package main

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"github.com/urfave/cli"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"fmt"
)
var flags struct {
	FileIn     string
	FileOut     string

}

func main() {
	app := cli.NewApp()
	app.Name = "structnames"
	app.Usage = "parse a file and cat all go struct names to an output file"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "input",
			Usage:       "inputfile",
			Destination: &flags.FileIn,
		},
		cli.StringFlag{
			Name:        "output",
			Usage:       "output file",
			Value: "./structNames.txt",
			Destination: &flags.FileOut,
		},
	}
	app.Action = launch
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

}

func launch(_ *cli.Context) error{
	file, _ := filepath.Abs(flags.FileIn)
	fmt.Println("Parsing file:", file)

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "Failed to read file.")
	}

	r, err := regexp.Compile("type\\s[A-Z]\\w+\\sstruct")
	if err != nil {
		return errors.Wrap(err, "Failed to execute regex.")
	}
	s := r.FindAllString(string(b), -1)


	fOut, err := os.Create(flags.FileOut)
	defer fOut.Close()
	var structNames []string
	for _, item := range s {
		item = strings.Replace(item, "type ", "", -1)
		item = strings.Replace(item, " struct", "\n", -1)
		structNames = append(structNames, item)
		fOut.WriteString(item)
	}
	return nil
}