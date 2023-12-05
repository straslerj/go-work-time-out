package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gen2brain/beeep"
)

const (
	notificationTimeout     = 30 * time.Second
	oneMinuteSeconds        = 60
	seventeenMinutesSeconds = 17 * oneMinuteSeconds
)

var (
	startTime           = time.Now()
	screenTimeTitle     = "Screen Time Alert"
	screenTimeMessage   = "You have been looking at your screen for 20 minutes. Consider looking at something 20 feet away for 20 minutes."
	screenTimeIcon      = "resources/screen_time.ico"
	screenTimeBreaks    = 0
	mentalBreakTitle    = "Mental Break Alert"
	mentalBreakMessage  = "You have been working for 52 minutes. Research shows you should take a 17-minute break to combat fatigue."
	mentalBreakIcon     = "resources/mental_break.ico"
	mentalBreakBreaks   = 0
	notificationCounter = 0
)

func main() {
	fmt.Println("Work Time Out has began... Work hard!")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				handleTimeElapsed()
			case <-signalCh:
				handleInterrupt()
				return
			}
		}
	}()

	select {}
}

func handleTimeElapsed() {
	elapsedMinutes := int(time.Since(startTime).Minutes())
	fmt.Printf("[%s] Elapsed time: %d minutes.\n", currentTime(), elapsedMinutes)

	if elapsedMinutes != 0 {
		if elapsedMinutes%20 == 0 {
			fmt.Printf("[%s] Screen time alert issued.\n", currentTime())
			notify(screenTimeTitle, screenTimeMessage, screenTimeIcon)
			screenTimeBreaks++
		}
		if elapsedMinutes%52 == 0 {
			fmt.Printf("[%s] Mental break alert issued.\n", currentTime())
			notify(mentalBreakTitle, mentalBreakMessage, mentalBreakIcon)
			sleepUntil := time.Now().Add(17 * time.Minute).Format("15:04:05")
			fmt.Printf("[%s] Program will sleep for 17 minutes. It will resume at %s\n", currentTime(), sleepUntil)
		}
	}
}

func currentTime() string {
	return time.Now().Format("15:04:05")
}

func notify(title, message, icon string) {
	err := beeep.Notify(title, message, icon)
	if err != nil {
		fmt.Println("Error sending notification:", err)
	}
	notificationCounter++
}

func handleInterrupt() {
	fmt.Printf("\n[%s] Program terminated.\n", currentTime())
	fmt.Printf("\t    Screen time breaks taken : %d\n", screenTimeBreaks)
	fmt.Printf("\t    Mental break breaks taken: %d\n", mentalBreakBreaks)
	fmt.Printf("\t    Total notifications sent : %d\n", notificationCounter)
	os.Exit(0)
}
