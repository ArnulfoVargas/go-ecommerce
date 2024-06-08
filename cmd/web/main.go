package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const version string = "1.0.0"
const cssVersion string = "1"

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

type Application struct {
	config
	infoLog *log.Logger
  errorLog *log.Logger
  templateCache map[string]*template.Template
  version string
  idleTimeOut time.Duration
  readTimeOut time.Duration
  headerTimeOut time.Duration
  writeTimeOut time.Duration
}

func (app *Application) Serve() error {
  srv := &http.Server{
    Addr: fmt.Sprintf(":%d", app.config.port),
    Handler: app.Routes(),
    IdleTimeout: app.idleTimeOut,
    ReadTimeout: app.readTimeOut,
    ReadHeaderTimeout: app.headerTimeOut,
    WriteTimeout: app.writeTimeOut,
  }

  app.infoLog.Printf("Starting HTTP server in %s mode on port: %d\n", app.config.env, app.config.port)

  return srv.ListenAndServe()
}

type Option func(app *Application)

func NewApplication(opts ...Option) *Application {
  var config config

  flag.IntVar(&config.port, "port", 3001, "Server port to listen on")
  flag.StringVar(&config.env, "env", "development", "Application environment {application | production}")
  flag.StringVar(&config.api, "api", "http://localhost:3000", "URL to api")

  flag.Parse()

  config.stripe.key = os.Getenv("STRIPE_KEY")
  config.stripe.secret = os.Getenv("STRIPE_SECRET")

  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)
  errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)
  templateCache := make(map[string]*template.Template);

  app := &Application{}

  app.infoLog = infoLog
  app.errorLog = errorLog
  app.templateCache = templateCache
  app.config = config
  app.version = version
  app.idleTimeOut = time.Duration(5) * time.Second
  app.headerTimeOut = time.Duration(5) * time.Second
  app.writeTimeOut = time.Duration(5) * time.Second
  app.readTimeOut = time.Duration(5) * time.Second

  for _, option := range opts {
    option(app)
  }

  return app
}

func SetIdleTimeOut(timeOut int) Option {
  return func(app *Application) {
    app.idleTimeOut = time.Duration(timeOut) * time.Second
  }
}

func SetWriteTimeOut(timeOut int) Option {
  return func(app *Application) {
    app.writeTimeOut = time.Duration(timeOut) * time.Second
  }
}

func SetReadTimeOut(timeOut int) Option {
  return func(app *Application) {
    app.readTimeOut = time.Duration(timeOut) * time.Second
  }
}

func SetHeaderTimeOut(timeOut int) Option {
  return func(app *Application) {
    app.headerTimeOut = time.Duration(timeOut) * time.Second
  }
}

func main() {
  app := NewApplication(
    SetIdleTimeOut(30),
    SetReadTimeOut(10),
    SetHeaderTimeOut(5),
    SetWriteTimeOut(5),
  )

  err := app.Serve()
  if err != nil {
    app.errorLog.Println(err)
    log.Fatal(err)
  }
}
