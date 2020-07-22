package cmd

import (
	"context"
	"goteach/logger"
	"goteach/routes"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start goteach http service",
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()
		e.Pre(middleware.RemoveTrailingSlash())
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: logger.MiddlewareLog,
		}))
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}))
		// Handler for hooking any request in routers registered and log it
		e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
			Handler: logger.APILogHandler,
			Skipper: logger.APILogSkipper,
		}))
		// Handler for putting goteach request and response timestamp. This used for get elapsed time
		e.Use(ServiceRequestTime)

		routes.Api(e)

		e.GET("/", func(c echo.Context) error {
			message := "Wani mati Wedi luwe, by Garasiman"
			return c.String(http.StatusOK, message)
		})

		// Start server
		// e.Logger.Fatal(e.Start(":8000"))
		go func() {
			if err := e.Start(":"+viper.GetString("port")); err != nil {
				e.Logger.Info("Shutting down the server")
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}

	},
}

// ServiceRequestTime middleware adds a `Server` header to the response.
func ServiceRequestTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("X-goteach-RequestTime", time.Now().Format(time.RFC3339))
		return next(c)
	}
}
