package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	outputDir = "output"
)

func main() {
	log.SetFlags(log.Lshortfile)

	buildOnly := flag.Bool("buildOnly", false, "only build all the metrics")
	flag.Parse()

	os.RemoveAll(outputDir)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	if *buildOnly {
		for _, metric := range metrics() {
			fmt.Println("-->", metric)

			cmd := exec.Command("go", "build", "main.go", metric)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				fmt.Println(err, cmd.Args)
			}
		}

		err = os.Remove("main")
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	for _, metric := range metrics() {
		log.Println(metric)

		cmd := exec.Command("go", "run", "main.go", metric, "-batch")

		output, err := cmd.CombinedOutput()
		if err != nil {
			os.Stdout.Write(output)
			log.Fatalln(metric, err, cmd.Args)
		}

		ioutil.WriteFile(filepath.Join(outputDir, metric+".txt"), output, 0644)
	}
}

func metrics() []string {
	d, err := os.Open(".")
	if err != nil {
		log.Fatalln(err)
	}

	fis, err := d.Readdir(-1)
	if err != nil {
		log.Fatalln(err)
	}

	metrics := []string{}
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}

		if filepath.Ext(fi.Name()) != ".go" {
			continue
		}

		switch fi.Name() {
		case "main.go", "all.go":
			// skip
		default:
			metrics = append(metrics, fi.Name())
		}
	}

	err = d.Close()
	if err != nil {
		log.Fatalln(err)
	}

	return metrics
}
