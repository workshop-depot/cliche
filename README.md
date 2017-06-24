# about cliche
> an opinionated cli tool to generate a cli app; `cliche` generates a cli app skeleton to start quickly with, based on awesome [cli](https://github.com/urfave/cli) package. Just call `. build.sh`.

For parameters, [argify](https://github.com/dc0d/argify) helps to bind a struct to command line parameters. For example in this struct:

```go
var conf struct {
	Info string `envvar:"APP_INFO" usage:"sample app info" value:"bare app structure"`

	Sample struct {
		SubCommand struct {
			Param string `envvar:"-"`
		}
	}
}
```

`Info` will get map to an app level parameter named `--info`. And `Sample` contains parameters for a command named `sample`. Same as `sample` command, here we can tell we also have a `subcommand` sub command which has a parameter named `--param`.

We just implement our commands inside the generated `app.go` after we created the initial skeleton. If we try to regenerate the app, somewhere that already contains some files (with same names), `cliche` will show an alert and stops, to protect the existing code.

The default logging package used here is [colog](https://github.com/comail/colog) which provides a simple toolset to printout nice logs. It allows different levels by just adding some text to the message as:

```go
log.Print("warn: that's all it takes!")
```

If you are using the standard `log` package, you don't have to change anything. And you can use your favorite logging package of-course.

Some app level variables introduced in generated code. Specifically `ctx` is an app level context which is used to notifying other go-routines when the app is shutting down. And go-routines can register themselves in `wg` wait group so the app can wait for them for a specified duration to end their work using the helper function `finit` which should get called in a command as `defer finit(time.Second)`.

```go
var (
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
)
```

Also `ctx` with get closed by calling `cancel` when an interrupt signal (like `Ctrl+C`) is received.

`cliche` will add build-time, commit hash, go version and last git tag to the app if they are present when calling `. build.sh`.