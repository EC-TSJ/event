# Event

[![Home](https://godoc.org/github.com/gookit/event?status.svg)](file:///D:/EC-TSJ/Documents/CODE/SOURCE/Go/pkg/lib/cli)
[![Build Status](https://travis-ci.org/gookit/event.svg?branch=master)](https://travis-ci.org/)
[![Coverage Status](https://coveralls.io/repos/github/gookit/event/badge.svg?branch=master)](https://coveralls.io/github/)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/event)](https://goreportcard.com/report/github.com/)

> **[EN README](README.md)**

Event es una librería para gestionar los eventos.

## GoDoc

- [godoc for github](https://godoc.org/github.com/)

## Funciones Principales
--- 


> type  ***`event`***  struct {

		count  int
		name   string
		events map[string][]interface{}
}

> Interface ***INodeEventTarget***, y objeto ***NodeEventTarget*** con métodos:

- `NodeEventTarget(string) INodeEventTarget`
-	`On(string, interface{}) error`
- `Once(string, interface{}) error`
-	`AddListener(string, interface{}) error`
-	`Emit(string, ...interface{}) error`
-	`EventNames() []string`
-	`RemoveListener(string, int)`
-	`RemoveAllListeners(...string)`
-	`Off(string, int)`
-	`ListenerCount(string) int`

> Interface ***IEventEmitter***, y objeto ***EventEmitter*** con métodos:

- `EventEmitter(string) IEventEmitter`
- evento `newListener`
- evento `removeListener`
- `On(string, interface{}) error `
- `Once(string, interface{}) error `
- `AddListener(string, interface{}) error `
- `Emit(string, ...interface{}) error`
- `Has(string) bool `
- `EventNames() []string `
- `RemoveListener(string)` 
- `RemoveEvent(...string)`
- `Off(...string) `
- `ListenerCount(string) int`
- `GetMaxListeners() int`
- `SetMaxListeners(int)`
- `Listeners(string) []interface{}`
- `PrependListener(string, interface{})`

> Objeto ***event***, con métodos:

- `DefaultMaxListeners(i int)`
- `GetEventListeners(d Dispatcher, eventName string) []TListener` 
- `ListenerId(d Dispatcher, eventName string) (lista []string)` 
- `ListenerCount(d Dispatcher, eventName string) int `
- `SetMaxListeners(n int, dsp ...EventEmitter)` 

> Interface ***IEventTarget***, con métodos:

-	`AddEventListener(string, interface{}, bool) error`
-	`DispatchEvent(string, ...interface{})`
-	`RemoveEventListener(string, int)`


> Interface ***IEvent*** y objeto ***Eventer***, con métodos:

- `StopPropagation()`
- `IsPropagationStopped() bool`
-	`Cancelable(bool)`
-	`PreventDefault()`
-	`DefaultPrevented() bool`
-	`ReturnValue() bool`
-	`Target() IEventTarget`


## Ejemplos

```go
Ejemplo 1

type (
	Listener interface {
		Name() string
		Handle(e event.Eventer) error
	}
)

type (
	fooEvent struct {
		event.Event

		i, j int
	}

	fooListener struct {
		name string
	}
)

func NewFooEvent(i, j int) event.Eventer {
	fmt.Println("Hola 1")
	return &fooEvent{i: i, j: j}
}

func NewFooListener() Listener {
	return &fooListener{
		name: "my.foo.event",
	}
}

func (l *fooListener) Name() string {
	return l.name
}

func (l *fooListener) Handle(e event.Eventer) error {
	ev := e.(*fooEvent)
	ev.StopPropagation()

	fmt.Println("Fire", ev.i+ev.j)
	return nil
}

func main() {
	e := event.New()
	e.Once("my.event.name.1", func() error {
		fmt.Println("Fire event")
		return nil
	})

	e.On("my.event.name.2", func(text string) error {
		fmt.Println("Fire", text)
		return nil
	})

	e.On("my.event.name.2", func(text string) error {
		fmt.Println("Fire YPP", text)
		return nil
	})

	e.On("my.event.name.3", func(i, j int) error {
		fmt.Println("Fire", i+j)
		return nil
	})

	e.On("my.event.name.4", func(name string, params ...string) error {
		fmt.Println(name, params)
		return nil
	})

	fmt.Println(handler.ListenerCount("my.event.name.2"))
	fmt.Println(handler.ListenerCount("my.event.name.1"))
	fmt.Println(handler.GetMaxListeners())
	handler.SetMaxListeners(7)
	fmt.Println(handler.GetMaxListeners())
	newlistener := handler.Listeners("my.event.name.2")
	fmt.Println(newlistener)
	newlistener[0].(func(string) error)("Joer")


	e.Emit("my.event.name.1")                           // Print: Fire event
	e.Emit("my.event.name.1", "joder")                  // Print: Fire event
	e.Emit("my.event.name.2", "some event")             // Print: Fire some event
	e.Emit("my.event.name.3", 1, 2)                     // Print: Fire 3
	e.Emit("my.event.name.4", "params:", "a", "b", "c") // Print: params: [a b c]
	e.Emit("my.event.name.2", "some event")             // Print: Fire some event
	e.Emit("my.event.name.3", 1, 2)                     // Print: Fire 3
	e.Emit("my.event.name.4", "params:", "a", "b", "c") // Print: params: [a b c]
	fmt.Println(e.EventNames())

	collect := []Listener{
		NewFooListener(),
	}

	// Registration
	for _, l := range collect {
		e.On(l.Name(), l.Handle)
	}

	// Call
	e.Emit("my.foo.event", NewFooEvent(1, 2))
}


```


``` go
Ejemplo 2

	//event.Debug = true


	handler := event.New("Handler")

	handler.On("newListener", func(event, fn string) error {
		s := fmt.Sprintf("Evento: '%s', Listener: '%s'. Dado de ALTA.", event, fn)
		logr.Info(s)
		fmt.Println(s)
		return nil
	})

	handler.Once("my.event.name.1", func() error {
		fmt.Println("Fire event")
		return errors.New("mmm")
	})

	handler.On("my.event.name.2", func(text string) error {
		fmt.Println("Fire", text)
		return nil
	})

	handler.Once("my.event.name.2", func(text string) error {
		fmt.Println("Fire YEXY", text)
		return nil
	})

	handler.PrependListener("my.event.name.2", func(text string) error {
		fmt.Println("Fire TERCERA", text)
		return nil
	})

	handler.Once("my.event.name.2", func(text string) error {
		fmt.Println("Fire CUARTA", text)
		return nil
	})

	handler.On("my.event.name.3", func(i, j int) error {
		fmt.Println("Fire", i+j)
		return nil
	})

	handler.On("my.event.name.4", func(name string, params ...string) error {
		fmt.Println("Fire", name, params)
		return nil
	})

	handler.On("removeListener", func(fn string) error {
		s := fmt.Sprintf("Listener: '%s'. Dado de BAJA.", fn)
		logr.Info(s)
		fmt.Println(s)
		return nil
	})

	fmt.Println(handler.ListenersCount("my.event.name.2"))
	fmt.Println(handler.ListenersCount("my.event.name.1"))
	fmt.Println(handler.GetMaxListeners())
	handler.SetMaxListeners(7)
	fmt.Println(handler.GetMaxListeners())
	newlistener := handler.Listeners("my.event.name.2")
	fmt.Println(newlistener)
	newlistener[0].Fn.(func(string) error)("Joer")

	if handler.Has("my.event.name.1") {
		handler.Emit("my.event.name.1") // Print: Fire even and bye
	} // Print: Fire even)
	if handler.Has("my.event.name.1") {
		handler.Emit("my.event.name.1")
	}
	handler.Emit("my.event.name.1")
	handler.Emit("my.event.name.2", "ya some event")          // Print: Fire some event
	handler.Emit("my.event.name.3", 1, 2)                     // Print: Fire 3
	handler.Emit("my.event.name.4", "params:", "a", "b", "v") // Print: params: [a b c]
	fmt.Println(handler.EventNames())
	handler.RemoveEvent()
	handler = event.New("Handler")
	handler.Emit("__End__")
```

``` go
Ejemplo 3

type (
	Listener interface {
		Name() string
		Handle(event.Eventer) error
	}
)

type (
	fooEvent struct {
		event.Event

		i, j int
	}

	fooListener struct {
		name string
	}
)

func NewFooEvent(i, j int) event.Eventer {
	fmt.Println("FooEvent")
	return &fooEvent{i: i, j: j}
}

func NewFooListener() Listener {
	return &fooListener{
		name: "my.foo.event",
	}
}

func (l *fooListener) Name() string {
	return l.name
}

func (l *fooListener) Handle(e event.Eventer) error {
	ev := e.(*fooEvent)
	ev.Cancelable(true)
	if ev.PreventDefault() {
		return nil
	}
	//ev.StopPropagation()

	fmt.Println("Fire", ev.i+ev.j)
	return nil
}

func Handle(e event.Eventer) error {
	ev := e.(*fooEvent)
	//ev.StopPropagation()

	fmt.Println("LaHostia", ev.i+ev.j)
	return nil
}

func main() {
	handler := event.New("Handler")
	l := NewFooListener()
	handler.On(l.Name(), l.Handle)
	handler.On(l.Name(), Handle)
	handler.On(l.Name(), func(e event.Eventer) error {
		fmt.Print("Joder")
		return nil
	})

	// Call
	handler.Emit("my.foo.event", &fooEvent{i: 2, j: 5})
}
```

``` go
Ejemplo 4.

type (
	Listener interface {
		Name() string
		Handle(event.Eventer) error
	}
)

type (
	fooEvent struct {
		event.Event

		i, j int
	}

	fooListener struct {
		name string
	}
)


func (l *fooListener) Name() string {
	return l.name
}

func (l *fooListener) Handle(e event.Eventer) error {
	ev := e.(*fooEvent)
	//ev.StopPropagation()

	fmt.Println("Fire", ev.i+ev.j)
	return nil
}

func Handle(e event.Eventer) error {
	ev := e.(*fooEvent)

	fmt.Println("LaHostia", ev.i+ev.j)
	return nil
}

func main() {
	handler := event.New("Handler")
	l := &fooListener{name: "my.foo.event"}

	handler.On(l.Name(), l.Handle)
	handler.On(l.Name(), Handle)
	handler.On(l.Name(), func(e event.Eventer) error {
		fmt.Println("Joder")
		return nil
	})
	handler.On("Magnolia", func(e event.Eventer) error {
		fmt.Println("YAVA")
		return nil
	})

	kal := &fooEvent{i: 4, j: 6}
	kal.Cancelable(false)
	kal.PreventDefault()
	kk := &event.Event{}

	// Call
	handler.Emit("my.foo.event", kal)
	handler.Emit("my.foo.event", kk)
	handler.Emit("Magnolia", &fooEvent{i: 2, j: 5})
```
## Notas

<!-- - [gookit/ini](https://github.com/gookit/ini) INI -->
## LICENSE

**[MIT](LICENSE)**
