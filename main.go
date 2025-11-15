package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Nouveau-Printemps/webring/backend"
	"github.com/joho/godotenv"
)

var (
	//go:embed dist
	embeds embed.FS
)

var (
	configFile = "config.toml"
	dev        = false
	port       = 8000
)

func init() {
	err := godotenv.Load(".env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Error("loading .env", "error", err)
	}

	if v := os.Getenv("CONFIG_FILE"); v != "" {
		configFile = v
	}
	flag.StringVar(&configFile, "config", configFile, "config file")

	if v := os.Getenv("PORT"); v != "" {
		port, err = strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
	}
	flag.IntVar(&port, "port", port, "server port")
	flag.BoolVar(&dev, "dev", false, "development mode")

}

func main() {
	flag.Parse()

	backend.SetupLogger(dev)

	cfg, ok := backend.LoadConfig(configFile)
	if !ok {
		slog.Info("exiting")
		os.Exit(1)
	}

	assetsFS := backend.UsableEmbedFS("dist", embeds)
	if dev {
		assetsFS = os.DirFS("dist")
	}

	r := backend.NewRouter(dev, cfg, assetsFS)

	r.Get("/", new(backend.Data).HandleGenericLater(
		"home",
		cfg.Translation.HomePage,
		strings.Join(cfg.Description, " ")),
	)
	r.Get("/join", new(backend.Data).HandleGenericLater("join", cfg.Translation.JoinTheRingPage, ""))
	r.Get("/legal", new(backend.Data).HandleGenericLater("legal", cfg.Translation.LegalInformationPage, ""))
	r.Get("/random", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, cfg.Websites[rand.Intn(len(cfg.Websites))].URL, http.StatusFound)
	})
	r.NotFound(new(backend.Data).HandleGenericLater("404", "Not found", "HTTP 404"))

	backend.HandleStaticFiles(r, "/assets", assetsFS)
	backend.HandleStaticFiles(r, "/static", os.DirFS(cfg.PublicFolder))

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r}

	errChan := make(chan error)
	go startServer(server, errChan)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	ok = true
	for ok {
		select {
		case err := <-errChan:
			slog.Error("http server running", "error", err)
			slog.Info("restarting the server")
			go startServer(server, errChan)
		case <-sc:
			err := server.Shutdown(context.Background())
			if err != nil {
				slog.Error("closing http server", "error", err)
			}
			ok = false
		}
	}
	slog.Info("http server stopped")
}

func startServer(server *http.Server, errChan chan error) {
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		errChan <- err
	}
}
