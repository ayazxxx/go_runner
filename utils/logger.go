package utils

import "fmt"

func Info(message string) {
    fmt.Println("\033[34m[INFO]\033[0m " + message)
}

func Success(message string) {
    fmt.Println("\033[32m[SUCCESS]\033[0m " + message)
}

func Error(message string) {
    fmt.Println("\033[31m[ERROR]\033[0m " + message)
}
