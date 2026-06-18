package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func MySQLDSN() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		rawURL := dsn
		if !strings.HasPrefix(rawURL, "mysql://") && !strings.Contains(rawURL, "://") {
			rawURL = "mysql://" + rawURL
		}
		if parsed, err := url.Parse(rawURL); err == nil && parsed.Scheme == "mysql" {
			user := parsed.User.Username()
			pass, _ := parsed.User.Password()
			host := parsed.Host
			name := parsed.Path
			if len(name) > 0 && name[0] == '/' {
				name = name[1:]
			}
			if strings.Contains(host, "tcp(") {
				return dsn
			}
			
			var userPass string
			if pass != "" {
				userPass = fmt.Sprintf("%s:%s", user, pass)
			} else if user != "" {
				userPass = user
			} else {
				userPass = "root"
			}
			
			query := parsed.RawQuery
			newDSN := fmt.Sprintf("%s@tcp(%s)/%s", userPass, host, name)
			if query != "" {
				newDSN += "?" + query
				if !strings.Contains(query, "parseTime=") {
					newDSN += "&parseTime=True"
				}
				if !strings.Contains(query, "loc=") {
					newDSN += "&loc=Local"
				}
			} else {
				newDSN += "?charset=utf8mb4&parseTime=True&loc=Local"
			}
			return newDSN
		}
		return dsn
	}
	user := GetEnv("DB_USER", "root")
	pass := GetEnv("DB_PASSWORD", "")
	host := GetEnv("DB_HOST", "localhost")
	port := GetEnv("DB_PORT", "3306")
	name := GetEnv("DB_NAME", "rpa")
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)
}
