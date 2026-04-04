package main

import "github.com/FrancoPesenda/eventra/internal/config"

func main() {
	_ = config.LoadDotEnv(".env")
	NewFxApp().Run()
}
