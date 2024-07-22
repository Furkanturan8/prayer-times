package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"namaz-vakitleri/database"
	"namaz-vakitleri/handlers"
	"namaz-vakitleri/routes"
	"namaz-vakitleri/services"
	"time"
)

func main() {
	fmt.Println("\n--------------BİSMİLLAH--------------\n")

	app := fiber.New()

	// Middleware to log requests
	app.Use(requestLogger)

	// DB Init
	db, err := database.DBInstance()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("DB Address: ", db)

	// API Base URL
	apiBaseURL := "https://api.aladhan.com/v1/calendarByCity/2024/7?country=turkey&city="

	prayerTimesService := services.NewPrayerTimeService(apiBaseURL)
	cityService := services.NewCityService()
	phraseService := services.NewPhraseService(db)

	prayerTimesHandler := handlers.NewPrayerTimeHandler(prayerTimesService)
	cityHandler := handlers.NewCityHandler(cityService)
	phraseHandler := handlers.NewPhraseHandler(phraseService)

	routes.PhraseRoutes(app, phraseHandler)
	routes.CityRoutes(app, cityHandler)
	routes.PrayerTimeRoutes(app, prayerTimesHandler)

	app.Listen(":3000")
}

func requestLogger(c *fiber.Ctx) error {
	start := time.Now()

	// Proceed to the next middleware or handler
	err := c.Next()

	stop := time.Now()
	latency := stop.Sub(start)

	// Get the status code and method
	status := c.Response().StatusCode()
	method := c.Method()
	url := c.OriginalURL()

	// Get the client IP
	clientIP := c.IP()

	// Log format: time | status | latency | clientIP | method | url
	fmt.Printf("%s | %3d | %9v | %15s | %-7s | %s\n",
		start.Format("15:04:05"),
		status,
		latency,
		clientIP,
		method,
		url,
	)

	return err
}