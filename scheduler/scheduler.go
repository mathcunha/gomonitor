package scheduler

import (
	"github.com/mathcunha/gomonitor/db"
	"log"
	"regexp"
	"strconv"
	"time"
)

func LoadMonitors() error {
	var monitor db.Monitor
	monitors, err := monitor.FindAll()

	if err != nil {
		log.Printf("error loading the montitors %v", err)
		return err
	}

	length := len(monitors)
	tickers := make([]*time.Ticker, length, length)
	controls := make([]chan bool, length, length)

	for i, m := range monitors {
		log.Printf("monitor %v", m.Id)
		duration := getInterval(m.Interval)

		if duration > 0 {
			tickers[i] = time.NewTicker(duration)
			controls[i] = make(chan bool)
			go schedule(tickers[i], controls[i], m)
		}
	}
	return nil
}

func schedule(t *time.Ticker, q chan bool, m db.Monitor) {
	Evaluate(m)
	for {
		select {
		case <-t.C:
			log.Printf("event - %v", m.Id)
			Evaluate(m)
		case <-q:
			t.Stop()
			return
		}
	}
}

func getInterval(interval string) time.Duration {
	log.Printf("interval %v", interval)
	nPattern := "^[0-9]*"
	dPattern := "[hms]$"

	if matched, _ := regexp.MatchString(nPattern+dPattern, interval); matched {
		re := regexp.MustCompile(nPattern)
		num, _ := strconv.Atoi(re.FindString(interval))

		re = regexp.MustCompile(dPattern)
		duration := re.FindString(interval)

		log.Printf("Num = %v - Duration = %v", num, duration)

		switch {
		case "h" == duration:
			return time.Duration(num) * time.Hour
		case "m" == duration:
			return time.Duration(num) * time.Minute
		case "s" == duration:
			return time.Duration(num) * time.Second
		}
	}

	return -1
}
