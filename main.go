package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"project/tech-blog-go/db"
	"project/tech-blog-go/handlers"
	"project/tech-blog-go/utils"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	var cfg handlers.Config

	// .envファイル作成
	err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		log.Fatalln(err)
	}

	ep, _ := strconv.Atoi(os.Getenv("PORT")) // .envのPORT
	ee := os.Getenv("ENV")                   // .envのENV
	ed := os.Getenv("DATABASE_URL")          // .envのDATABASE_URL

	flag.IntVar(&cfg.Port, "port", ep, "Server port to Listen on")
	flag.StringVar(&cfg.Env, "env", ee, "Application environment (development|staging|production)")
	flag.StringVar(&cfg.Db.Dsn, "dsn", ed, "Postgres connetion string")
	flag.Parse()

	cfg.Logfile = "tb.log"
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	utils.LoggingSetting(cfg.Logfile)

	db, err := db.OpenDB(ed)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	srv := handlers.ServerSetting(cfg, logger, db)

	logger.Println("Starting server on port", cfg.Port)

	err = srv.ListenAndServe()

	if err != nil {
		log.Println(err)
	}
}
