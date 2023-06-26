package main

import (
	"context"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/spf13/viper"
)

// Notion Sync -> Todoist
// Get notion db
// init todoist columns
// store sync in notion db

// Query Objects with create/updated dates < sync
// for objects with no todoist id, create, for ones with id update

// Todoist Sync
// Get since last sync
// Query Notion db to see if they exist.
// Loop through results

type Config struct {
	NotionToken        string `mapstructure:"NOTION_TOKEN"`
	NotionDBID         string `mapstructure:"NOTION_DB_ID"`
	CreateDB           bool
	TodoistAccessToken string `mapstructure:"TODOIST_ACCESS_TOKEN"`
}

// type RichTextPropertyConfig struct {
// 	ID       notionapi.PropertyID         `json:"id,omitempty"`
// 	Type     notionapi.PropertyConfigType `json:"type"`
// 	RichText struct{}                     `json:"rich_text"`
// 	Name     string                       `json:"name,omitempty"`
// }

// func (r *RichTextPropertyConfig) GetType() notionapi.PropertyConfigType {
// 	return notionapi.PropertyConfigType(notionapi.PropertyTypeRichText)
// }

func main() {
	// lastSync := time.Now()
	cfg, err := NewConfig("")
	if err != nil {
		panic(err.Error())
	}
	cfg.Print()

	notionClient := notionapi.NewClient(notionapi.Token(cfg.NotionToken))

	db, err := notionClient.Database.Get(context.TODO(), notionapi.DatabaseID(cfg.NotionDBID))
	if err != nil {
		panic(err.Error())
		// if db not found && CreateDB then create
	}

	_, found := db.Properties["todoist_id"]
	if !found {
		db, err = notionClient.Database.Update(context.TODO(), notionapi.DatabaseID(db.ID), &notionapi.DatabaseUpdateRequest{
			Properties: map[string]notionapi.PropertyConfig{
				"todoist_id": notionapi.RichTextPropertyConfig{
					Type: "rich_text",
				},
			},
		})
		if err != nil {
			panic(err.Error())
		}
	}

	// pages, err := GetNotionPagesFromDb(notionClient, string(db.ID))
	// if err != nil {
	// 	panic(err.Error())
	// }

	// // log.Printf("Pages: %v\n", pages)

	// for i, v := range pages {
	// 	log.Printf("PAGE #%v\n", i)
	// 	log.Printf("ID: %v \n", v.ID)
	// 	log.Printf("Obj: %v \n", v.Object)
	// 	log.Printf("Properties: %v \n", v.Properties)
	// 	for j, k := range v.Properties {
	// 		log.Printf("K: %v | V: %v | type: %v \n", j, k, k.GetType())
	// 	}

	// 	// convert to todoist tasks
	// 	// sync

	// }

	// tCfg := todoist.Config{
	// 	AccessToken: cfg.TodoistAccessToken,
	// }
	// todoistClient := todoist.NewClient(&tCfg)
	// err = todoistClient.Sync(context.TODO())
	// if err != nil {
	// 	panic(err.Error())
	// }
	// for k, v := range todoistClient.Store.Items {
	// 	log.Printf("#%v: %v \n\n", k, v)
	// }
}

// Gets all items from db.
func GetNotionPagesFromDb(client *notionapi.Client, dbID string) ([]notionapi.Page, error) {
	startCursor := notionapi.Cursor("")
	hasMore := true
	pages := []notionapi.Page{}

	for hasMore {
		res, err := client.Database.Query(context.TODO(), notionapi.DatabaseID(dbID), &notionapi.DatabaseQueryRequest{
			StartCursor: startCursor,
		})
		if err != nil {
			return nil, err
		}
		pages = append(pages, res.Results...)
		hasMore = res.HasMore
		startCursor = res.NextCursor
	}
	return pages, nil
}

func NewConfig(dir string) (*Config, error) {
	if dir == "" {
		dir = "."
	}

	viper.SetDefault("NOTION_TOKEN", "")
	viper.SetDefault("NOTION_DB_ID", "")
	viper.SetDefault("TODOIST_ACCESS_TOKEN", "")
	// viper.SetEnvPrefix("")
	viper.SetConfigName("config")
	viper.AddConfigPath(dir)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.ReadInConfig()

	// Environment variables beat config file
	var C Config
	err := viper.Unmarshal(&C)
	if err != nil {
		return nil, err
	}
	return &C, nil
}

func (c *Config) Print() {
	println("NOTION_TOKEN", c.NotionToken)
}
