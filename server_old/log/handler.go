package log

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/felixge/httpsnoop"
	"github.com/pbivrell/gatekeeper/log"
)

func Middlewear(logger log.Logger, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := "n/a"
		c, err := r.Cookie("token")
		if err == nil {
			claims := Claims{}
			_, err = jwt.ParseWithClaims(c.Value, &claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err == nil {
				user = claims.Username
			}
		}

		t := time.Now()

		// Preform the request and capture metrics about the result
		m := httpsnoop.CaptureMetrics(handler, w, r)

		// Parse useragent
		ua := ua.Parse(r.Header.Get("User-Agent"))

		deviceType := "n/a"
		if ua.Bot {
			deviceType = "bot"
		} else if ua.Mobile {
			deviceType = "mobile"
		} else if ua.Tablet {
			deviceType = "tablet"
		} else if ua.Desktop {
			deviceType = "desktop"
		}

		logger.Write(log.Data{
			Method:           r.Method,
			Endpoint:         r.URL.String(),
			Referer:          r.Header.Get("Referer"),
			Code:             m.Code,
			Duration:         m.Duration.Milliseconds(),
			IP:               r.RemoteAddr[0:strings.LastIndex(r.RemoteAddr, ":")],
			UserAgent:        ua.Name,
			UserAgentVersion: ua.Version,
			OS:               ua.OS,
			OSVersion:        ua.OSVersion,
			Device:           ua.Device,
			DeviceType:       deviceType,
			Time:             t,
			User:             user,
		})
	}
}
