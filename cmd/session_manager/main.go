package main

import (
	"context"
	"session_manager/internal/server"
)

func main() {
	ctx := context.Background()

	// envorinments [db and etc...]
	env := server.NewEnv(ctx)
	defer env.Stop(ctx)

	// server
	srv := server.NewServer(env)
	srv.Run(ctx)
	defer srv.Stop(ctx)
}
