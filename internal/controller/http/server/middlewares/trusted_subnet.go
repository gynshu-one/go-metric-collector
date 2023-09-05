package middlewares

import (
	"github.com/3th1nk/cidr"
	"github.com/gin-gonic/gin"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	"github.com/rs/zerolog/log"
)

var CIDR *cidr.CIDR

// Parsing trusted subnet from config once
func init() {
	var err error
	var subNet = config.GetConfig().TrustedSubNet
	if subNet == "" {
		return
	}
	CIDR, err = cidr.Parse(subNet)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing trusted subnet")
	}
}
func CheckSubnet() gin.HandlerFunc {
	return func(c *gin.Context) {
		if CIDR == nil {
			c.Next()
			return
		}
		xReal := c.Request.Header.Get("X-Real-IP")
		if xReal == "" {
			c.AbortWithStatus(403)
			return
		}
		if !CIDR.Contains(xReal) {
			c.AbortWithStatus(403)
			return
		}
		c.Next()
	}
}
