package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-akka/configuration"
	_ "github.com/lib/pq"
)

func main() {
	conf := configuration.LoadConfig("credentials.conf")

	totalCount, err := getTotalGuildCount(conf)
	if err != nil {
		panic(err)
	}

	err = postStats(totalCount)
	if err != nil {
		panic(err)
	}

}

func getTotalGuildCount(conf *configuration.Config) (totalCount int, err error) {
	var (
		id    int
		count int
	)

	log.Println("[Database]", "Getting database connection info...")
	dbHost := conf.GetString("database.host")
	dbPort := conf.GetInt32("database.port")
	dbUser := conf.GetString("database.user")
	dbPass := conf.GetString("database.pass")
	dbName := conf.GetString("database.name")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	log.Println("[Database]", "Connecting with the database...")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return 0, err
	}
	log.Println("[Database]", "Successfully connected with the database!")

	err = db.Ping()
	if err != nil {
		return 0, err
	}

	log.Println("[Database]", "Fetching shard stats...")
	defer db.Close()
	rows, err := db.Query("SELECT * FROM discord.shard_guilds")
	if err != nil {
		return 0, err
	}

	log.Println("[Database]", "Counting results...")
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &count)
		if err != nil {
			return 0, err
		}
		totalCount = totalCount + count
	}
	err = rows.Err()
	if err != nil {
		return 0, err
	}

	log.Println("[Database]", "Success!")
	return totalCount, nil
}

func postStats(totalCount int) error {

	log.Println("[StatsPoster]", "Posting stats...")
	return fmt.Errorf("Not implemented")
}
