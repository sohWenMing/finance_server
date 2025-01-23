package envvars

import "github.com/joho/godotenv"

func LoadEnv(envFilePath string) error {
	err := godotenv.Load(envFilePath)
	if err != nil {
		return err
	}
	return nil
}
