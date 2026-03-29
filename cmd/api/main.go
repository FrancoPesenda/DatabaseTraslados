package main

func main() {
	_ = LoadDotEnv(".env")
	NewFxApp().Run()
}
