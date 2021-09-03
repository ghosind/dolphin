package dolphin

type App struct {
}

func New() *App {
	return &App{}
}

func Default() *App {
	app := New()

	return app
}

func (app *App) Run() {
	// TODO: implement
}
