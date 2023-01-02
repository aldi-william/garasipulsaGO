package server

import (
	"os"

	"github.com/gin-contrib/cors"
)

// allowCors is used to register url to user our resources
func AllowCors(environment string) {
	var (
		subDomainURL string = os.Getenv("DOMAIN_URL")
		whitelists   []string
	)
	config := cors.DefaultConfig()
	switch environment {
	case "DEVELOPMENT":
		whitelists = []string{"*"}
	default:
		whitelists = []string{subDomainURL}
	}

	config.AllowWildcard = true
	config.AllowOrigins = whitelists
	config.AllowHeaders = append(config.AllowHeaders, "Authorization", "Accept-Encoding", "X-Digiflazz-Event", "X-Digiflazz-Delivery", "X-Hub-Signature", "Content-Type", "X-Forwarded-For", "X-MOOTA-WEBHOOK", "X-MOOTA-USER")

	router.Use(cors.New(config))
}
