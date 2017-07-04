package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/goby-lang/goby/compiler"
	"github.com/goby-lang/goby/igb"
	"github.com/goby-lang/goby/vm"
	"github.com/pkg/profile"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Version stores current Goby version
const Version string = "0.0.8"

func main() {
	compileOptionPtr := flag.Bool("c", false, "Compile to bytecode")
	profileOptionPtr := flag.Bool("p", false, "Profile program execution")
	versionOptionPtr := flag.Bool("v", false, "Show current Goby version")
	interactiveOptionPtr := flag.Bool("i", false, "Run interactive goby")

	flag.Parse()

	if *interactiveOptionPtr {
		scanner := bufio.NewScanner(os.Stdin)

		ch := make(chan string)

		go func() {
			for {
				scanned := scanner.Scan()

				if !scanned {
					continue
				}

				ch <- scanner.Text()
			}
		}()

		igb.Start(ch, os.Stdout)
	}

	if *profileOptionPtr {
		defer profile.Start().Stop()
	}

	if *versionOptionPtr {
		fmt.Println(Version)
		os.Exit(0)
	}

	filepath := flag.Arg(0)

	if filepath == "" {
		flag.Usage()
		os.Exit(0)
	}

	args := flag.Args()[1:]

	dir, filename, fileExt := extractFileInfo(filepath)
	file := readFile(filepath)

	switch fileExt {
	case "gb", "rb":
		bytecodes, err := compiler.CompileToBytecode(string(file))

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if !*compileOptionPtr {
			v := vm.New(dir, args)
			v.ExecBytecodes(bytecodes, filepath)
			return
		}

		writeByteCode(bytecodes, dir, filename)
	case "gbbc":
		bytecodes := string(file)
		v := vm.New(dir, args)
		v.ExecBytecodes(bytecodes, filepath)
	default:
		fmt.Printf("Unknown file extension: %s", fileExt)
	}
}

func extractFileInfo(fp string) (dir, filename, fileExt string) {
	dir, filename = filepath.Split(fp)
	dir, _ = filepath.Abs(dir)
	fileExt = filepath.Ext(fp)
	splited := strings.Split(filename, ".")
	filename, fileExt = splited[0], splited[1]
	return
}

func writeByteCode(bytecodes, dir, filename string) {
	f, err := os.Create(dir + filename + ".gbbc")

	if err != nil {
		panic(err)
	}

	f.WriteString(bytecodes)
}

func readFile(filepath string) []byte {
	file, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	return file
}
