package configs

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	LNDHost              string `form:"LNDHost"`
	TaroHost             string `form:"TaroHost"`
	MacaroonLocation     string `form:"MacaroonLocation"`
	Macaroon             string `form:"Macaroon"`
	TaroMacaroonLocation string `form:"TaroMacaroonLocation"`
	TaroMacaroon         string `form:"TaroMacaroon"`
}

func (configs Config) GetConfigMap() (configMap map[string]string) {
	inrec, _ := json.Marshal(configs)
	json.Unmarshal(inrec, &configMap)
	return configMap
}

func GetConfig(ctx context.Context) (configs Config) {
	return ctx.Value("configs").(Config)
}

func LoadConfig(ctx context.Context) (context.Context, error) {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file, falling back to .env.sample: %v", err)
		if fatalErr := godotenv.Load("env/.env.sample"); fatalErr != nil {
			// load file from bindata.go
			// create dependencies
			// data, _ := Asset("env/.env.sample")
			// os.WriteFile(".env.sample", data, 0644)

			files := AssetNames()

			for _, file := range files {
				log.Println(file)
				data, _ := Asset(file)

				dir, _ := filepath.Split(file)

				if _, err := os.Stat(dir); os.IsNotExist(err) {
					// your file does not exist
					os.MkdirAll(dir, 0700)
				}

				err := os.WriteFile(file, data, 0644)
				log.Println(err)
			}

			if fatalErr := godotenv.Load("env/.env.sample"); fatalErr != nil {
				log.Fatalf(fatalErr.Error())
			}
		}
	}

	configs := Config{
		LNDHost:              os.Getenv("LNDHost"),
		TaroHost:             os.Getenv("TaroHost"),
		MacaroonLocation:     os.Getenv("MacaroonLocation"),
		Macaroon:             os.Getenv("Macaroon"),
		TaroMacaroonLocation: os.Getenv("TaroMacaroonLocation"),
		TaroMacaroon:         os.Getenv("TaroMacaroon"),
	}

	ctx = context.WithValue(ctx, "configs", configs)

	return ctx, err
}

// use godot package to load/read the .env file and
// return the value of the key.
func GoDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
