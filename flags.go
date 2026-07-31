package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Flags representa los argumentos de línea de comandos
type Flags struct {
	CreatePR    bool
	MergePR     bool
	MergePRNum  int // 0 significa no especificado
	Help        bool
}

// parseFlags parsea los argumentos de línea de comandos
func parseFlags() Flags {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "tuipr - PR Lifecycle TUI Manager\n\n")
		fmt.Fprintf(os.Stderr, "Uso:\n")
		fmt.Fprintf(os.Stderr, "  tuipr              Abrir dashboard principal\n")
		fmt.Fprintf(os.Stderr, "  tuipr -c           Abrir directamente en Create PR\n")
		fmt.Fprintf(os.Stderr, "  tuipr -m           Abrir pantalla de Merge (seleccionar PR)\n")
		fmt.Fprintf(os.Stderr, "  tuipr -m <num>     Merge directo del PR #<num>\n")
		fmt.Fprintf(os.Stderr, "  tuipr --help       Mostrar esta ayuda\n")
		fmt.Fprintf(os.Stderr, "\nNavegación:\n")
		fmt.Fprintf(os.Stderr, "  1-3    Cambiar entre paneles\n")
		fmt.Fprintf(os.Stderr, "  j/k    Navegar arriba/abajo\n")
		fmt.Fprintf(os.Stderr, "  i      Modo Insert\n")
		fmt.Fprintf(os.Stderr, "  Esc    Volver a Normal / Volver\n")
		fmt.Fprintf(os.Stderr, "\nAcciones:\n")
		fmt.Fprintf(os.Stderr, "  c      Crear PR\n")
		fmt.Fprintf(os.Stderr, "  m      Merge PR\n")
		fmt.Fprintf(os.Stderr, "  e      Abrir nvim para resolver conflictos\n")
		fmt.Fprintf(os.Stderr, "  r      Refrescar lista\n")
		fmt.Fprintf(os.Stderr, "  q      Salir\n")
	}

	createPtr := flag.Bool("c", false, "Abrir directamente en pantalla de Create PR")
	mergePtr := flag.Bool("m", false, "Abrir directamente en pantalla de Merge PR")
	helpPtr := flag.Bool("help", false, "Mostrar ayuda")

	flag.Parse()

	// Verificar si hay un número después de -m
	mergeNum := 0
	if flag.NArg() > 0 {
		numStr := flag.Arg(0)
		// Permitir "tuipr -m 134" o "tuipr 134"
		if num, err := strconv.Atoi(numStr); err == nil {
			mergeNum = num
		}
	}

	return Flags{
		CreatePR:   *createPtr,
		MergePR:    *mergePtr || mergeNum > 0,
		MergePRNum: mergeNum,
		Help:       *helpPtr,
	}
}
