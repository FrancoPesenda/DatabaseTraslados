package main

import "github.com/FrancoPesenda/eventra/internal/utils/config"

func main() {
	_ = config.LoadDotEnv(".env")
	NewFxApp().Run()
}
