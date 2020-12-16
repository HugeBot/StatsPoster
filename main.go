package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-akka/configuration"
	_ "github.com/lib/pq"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	confPath := filepath.Join(fmt.Sprintf("%s/credentials.conf", wd))
	log.Printf("Config file directory %s", confPath)

	conf := configuration.LoadConfig(confPath)
	client := &http.Client{}

	totalCount, err := getTotalGuildCount(conf)
	if err != nil {
		panic(err)
	}

	err = postStats(client, conf, totalCount)
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
	rows, err := db.Query("SELECT id, count FROM discord.shard_guilds")
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

func postStats(client *http.Client, conf *configuration.Config, totalCount int) error {

	log.Println("[StatsPoster]", "Posting stats...")

	botID := conf.GetString("bot.id")
	topgg := conf.GetString("tokens.topgg")
	bfd := conf.GetString("tokens.bfd")
	dboats := conf.GetString("tokens.dboats")
	dbotsgg := conf.GetString("tokens.dbotsgg")

	if len(topgg) > 0 {
		json := fmt.Sprintf(`{"server_count": %d}`, totalCount)
		service := Service{
			name:  "[TopGG]",
			url:   fmt.Sprintf(`https://top.gg/api/bots/%s/stats`, botID),
			token: topgg,
			body:  bytes.NewBufferString(json),
		}
		postServiceStats(client, &service)
	}

	if len(bfd) > 0 {
		json := fmt.Sprintf(`{"server_count": %d}`, totalCount)
		service := Service{
			name:  "[BotsForDiscord]",
			url:   fmt.Sprintf(`https://botsfordiscord.com/api/bot/%s`, botID),
			token: bfd,
			body:  bytes.NewBufferString(json),
		}
		postServiceStats(client, &service)
	}

	if len(dboats) > 0 {
		json := fmt.Sprintf(`{"server_count": %d}`, totalCount)
		service := Service{
			name:  "[DiscordBoats]",
			url:   fmt.Sprintf(`https://discord.boats/api/bot/%s`, botID),
			token: dboats,
			body:  bytes.NewBufferString(json),
		}
		postServiceStats(client, &service)
	}

	if len(dbotsgg) > 0 {
		json := fmt.Sprintf(`{"guildCount": %d}`, totalCount)
		service := Service{
			name:  "[DiscordBotsGG]",
			url:   fmt.Sprintf(`https://discord.bots.gg/api/v1/bots/%s/stats`, botID),
			token: dbotsgg,
			body:  bytes.NewBufferString(json),
		}
		postServiceStats(client, &service)
	}

	return nil
}

func postServiceStats(client *http.Client, service *Service) {
	req, err := http.NewRequest("POST", service.url, service.body)

	req.Header.Add("Authorization", service.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(service.name, "Error on request.\n[ERRO] -", err)
	}

	defer resp.Body.Close()
	log.Println(service.name, resp.Status)
}

//Service ...
type Service struct {
	name  string
	url   string
	token string
	body  *bytes.Buffer
}
