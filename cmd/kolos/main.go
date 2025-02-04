package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProtonEvgeny/kolos/internal/aig"
	"github.com/ProtonEvgeny/kolos/internal/model"
	"github.com/ProtonEvgeny/kolos/internal/stats"
	"github.com/chzyer/readline"
)

var (
	loadedAIG   *model.AIG
	currentFile string
)

func main() {
	fmt.Println("KoLoS (Kosmos Logic Synthesis) v0.1.0")
	fmt.Println("Type 'help' for available commands")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       ">> ",
		HistoryFile:  "/tmp/kolos_history.tmp",
		AutoComplete: completer(),
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // EOF
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		args := strings.Fields(line)
		cmd := args[0]
		args = args[1:]

		switch cmd {
		case "read_aiger":
			if len(args) != 1 {
				fmt.Println("Usage: read_aiger <file>")
				continue
			}
			handleReadAiger(args[0])

		case "print_stats":
			handlePrintStats()

		case "clear":
			handleClear()

		case "quit", "exit":
			fmt.Println("Bye! UwU")
			return

		case "help":
			printHelp()

		default:
			fmt.Printf("Unknown command: %s\n", cmd)
		}
	}
}

func completer() readline.AutoCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("read_aiger"),
		readline.PcItem("print_stats"),
		readline.PcItem("clear"),
		readline.PcItem("quit"),
		readline.PcItem("help"),
	)
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("	read_aiger <file>         Load AIGER file")
	fmt.Println("	print_stats               Show statistics")
	fmt.Println("	clear                     Clear current network")
	fmt.Println("	quit                      Quit the program")
	fmt.Println("	help                      Show this help")
}

func handleClear() {
	loadedAIG = nil
	currentFile = ""
	fmt.Println("Cleared")
}

func handleReadAiger(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error loading file: %v\n", err)
		return
	}
	defer f.Close()
	aigGraph, err := aig.ParseAIG(f)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}
	aig.LinkNodes(aigGraph)
	loadedAIG = aigGraph
	currentFile = file
	fmt.Println("Successfully loaded")
}

func handlePrintStats() {
	if loadedAIG == nil {
		fmt.Println("Error: No network loaded. Use 'read_aiger <file>' first.")
		return
	}

	s := stats.Calculate(loadedAIG)
	name := strings.TrimSuffix(
		filepath.Base(currentFile),
		filepath.Ext(currentFile),
	)

	fmt.Printf("%s :\n", name)
	fmt.Printf("		I / O      = %d / %d\n", s.Inputs, s.Outputs)
	fmt.Printf("		Latches    = %d\n", s.Latches)
	fmt.Printf("		AND        = %d\n", s.AndGates)
	fmt.Printf("		Level      = %d\n", s.MaxLevel)
}
