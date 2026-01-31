package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// MovingClouds is the main component of the application.
// A component is a customizable, independent, and reusable UI element.
// It is created by embedding app.Compo into a struct.
type MovingClouds struct {
	app.Compo
}

type draggableButton struct {
	app.Compo
	left        int
	top         int
	dragging    bool
	offsetX     int
	offsetY     int
	Image       string
	onMouseMove app.Func
	onMouseUp   app.Func
}

func (b *draggableButton) Render() app.UI {
	btn := app.Button().
		Style("position", "absolute").
		Style("left", strconv.Itoa(b.left)+"px").
		Style("top", strconv.Itoa(b.top)+"px").
		Style("cursor", "move").
		OnMouseDown(b.startDrag)

	if b.Image != "" {
		btn = btn.Style("background-image", "url('"+b.Image+"')").
			Style("background-size", "cover").
			Style("background-position", "center").
			Style("width", "100px").
			Style("height", "100px").
			Style("background-color", "transparent"). // Make background transparent
			Style("border", "none").                  // Remove border
			Text("")
	} else {
		btn = btn.Text("Drag Me")
	}

	return btn
}

func (b *draggableButton) startDrag(ctx app.Context, e app.Event) {
	b.dragging = true
	ev := e.JSValue()
	clientX := ev.Get("clientX").Int()
	clientY := ev.Get("clientY").Int()
	b.offsetX = clientX - b.left
	b.offsetY = clientY - b.top

	// Define callbacks
	b.onMouseMove = app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		if !b.dragging {
			return nil
		}
		event := args[0]
		clientX := event.Get("clientX").Int()
		clientY := event.Get("clientY").Int()

		ctx.Dispatch(func(ctx app.Context) {
			b.left = clientX - b.offsetX
			b.top = clientY - b.offsetY
			// Trigger update
			ctx.Update() // Calling Update() on the component itself
		})
		return nil
	})

	b.onMouseUp = app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		b.dragging = false
		ctx.Dispatch(func(ctx app.Context) {
			app.Window().JSValue().Call("removeEventListener", "mousemove", b.onMouseMove)
			app.Window().JSValue().Call("removeEventListener", "mouseup", b.onMouseUp)
			b.onMouseMove.Release()
			b.onMouseUp.Release()
		})
		return nil
	})

	// Attach to window
	app.Window().JSValue().Call("addEventListener", "mousemove", b.onMouseMove)
	app.Window().JSValue().Call("addEventListener", "mouseup", b.onMouseUp)
}

// The Render method is where the component appearance is defined.
func (mc *MovingClouds) Render() app.UI {
	return app.Div().
		Style("background-image", "url('/web/moving-clouds.png')").
		Style("background-size", "cover").
		Style("background-position", "center").
		Style("min-height", "100vh").
		Style("position", "relative").
		Body(
			&draggableButton{
				Image: "/web/cloud.png",
				left:  50,
				top:   50,
			},
			&draggableButton{
				Image: "/web/cloud.png",
				left:  100,
				top:   50,
			},
			&draggableButton{
				Image: "/web/cloud.png",
				left:  150,
				top:   50,
			},
			&draggableButton{
				Image: "/web/cloud.png",
				left:  200,
				top:   50,
			},
		)
}

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	// The first thing to do is to associate the main component with a path.
	//
	// This is done by calling the Route() function, which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.Route("/", func() app.Composer { return &MovingClouds{} })

	// Once the routes set up, the next thing to do is to either launch the app
	// or the server that serves the app.
	//
	// When executed on the client-side, the RunWhenOnBrowser() function
	// launches the app,  starting a loop that listens for app events and
	// executes client instructions. Since it is a blocking call, the code below
	// it will never be executed.
	//
	// When executed on the server-side, RunWhenOnBrowser() does nothing, which
	// lets room for server implementation without the need for precompiling
	// instructions.
	app.RunWhenOnBrowser()

	// Define a flag to check if we should generate the static website
	genStatic := flag.Bool("static", false, "Generate static website")
	flag.Parse()

	if *genStatic {
		err := app.GenerateStaticWebsite(".", &app.Handler{
			Name:        "Moving Clouds Publishing",
			Description: "A Moving Clouds Web Application",
		})

		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Finally, launching the server that serves the app is done by using the Go
	// standard HTTP package.
	//
	// The Handler is an HTTP handler that serves the client and all its
	// required resources to make it work into a web browser. Here it is
	// configured to handle requests with a path that starts with "/".

	http.Handle("/", &app.Handler{
		Name:        "Moving Clouds Publishing",
		Description: "A Moving Clouds Web Application",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
