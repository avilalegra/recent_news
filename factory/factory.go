package factory

import (
	"avilego.me/recent_news/config"
	"avilego.me/recent_news/env"
	"avilego.me/recent_news/news"
	"avilego.me/recent_news/persistence"
	"avilego.me/recent_news/rss"
	"log"
	"os"
	"time"
)

func rssProvider() news.Provider {
	var sources []rss.Source
	for _, url := range config.Current.RNPConfig.Sources {
		sources = append(sources, rss.NewSource(url))
	}
	interval := time.Duration(config.Current.RNPConfig.MinutesPeriod) * time.Minute
	return rss.NewRssProvider(sources, time.Tick(interval))
}

func logger() *log.Logger {
	file, _ := os.OpenFile(env.LogFile(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	return log.New(file, "", log.LstdFlags)
}

func Collector() news.Collector {
	return news.Collector{
		Providers: []news.Provider{rssProvider()},
		Keeper:    Keeper(),
		Logger:    logger(),
	}
}

func Finder() news.Finder {
	return persistence.NewMongoKeeperFinder()
}

func Keeper() news.Keeper {
	return persistence.NewMongoKeeperFinder()
}

func Cleaner() news.Cleaner {
	ttl := config.Current.CleanerConfig.Ttl * 60
	interval := time.Duration(config.Current.CleanerConfig.MinutesPeriod) * time.Minute

	return news.Cleaner{
		KeeperFinder: persistence.NewMongoKeeperFinder(),
		Trigger:      time.Tick(interval),
		Ttl:          int64(ttl),
	}
}
