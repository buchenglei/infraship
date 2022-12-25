package app

// import "runtime"

type ApplicationOption func(*Application) error

// func AppWithBallast(defaultAlloc ...int) CliAppOptionFunc {
// 	defaultSize := 10 * 1024 * 1024 * 1024
// 	if len(defaultAlloc) == 1 {
// 		defaultSize = defaultAlloc[0]
// 	}

// 	return func(c *App) error {
// 		ballast := make([]byte, defaultSize)
// 		runtime.KeepAlive(ballast)

// 		return nil
// 	}
// }
