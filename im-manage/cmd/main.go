package main

func main() {
	// 初始化
	app, err := initApp()
	if err != nil {
		panic(err)
	}

	app.run()
}
