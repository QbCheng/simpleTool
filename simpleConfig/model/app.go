package model

type App struct {
	Mysql struct {
		Host     string
		User     string
		Password string
	}
	Redis struct {
		Host string
		Auth string
	}
}
