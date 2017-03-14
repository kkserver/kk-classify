package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kkserver/kk-classify/classify"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"log"
	"os"
)

func main() {

	log.SetFlags(log.Llongfile | log.LstdFlags)

	env := "./config/env.ini"

	if len(os.Args) > 1 {
		env = os.Args[1]
	}

	a := classify.ClassifyApp{}

	err := app.Load(&a, "./app.ini")

	if err != nil {
		log.Panicln(err)
	}

	err = app.Load(&a, env)

	if err != nil {
		log.Panicln(err)
	}

	app.Obtain(&a)

	app.Handle(&a, &app.InitTask{})

	kk.DispatchMain()

	app.Recycle(&a)

}
