package xrest

import (
	"github.com/rs/cors"
)

type CORSOptions = cors.Options

func CORS() Middleware {
	return cors.Default().Handler
}

func CORSWithOptions(options CORSOptions) Middleware {
	return cors.New(options).Handler
}
