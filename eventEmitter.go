// antiguo nombre IEventEmitter.v
package event

import (
	"bytes"
	"crypto/rand"
	"ec-tsj/core"
	log "ec-tsj/logger"
	"errors"
	"fmt"
	"reflect"
)

// IEventEmitter event interface
type IEventEmitter interface {
	On(string, core.T /* TListener */) error          //
	Once(string, core.T /* TListener */) error        //
	AddListener(string, core.T /* TListener */) error //
	Emit(string, ...core.T /* Args*/) error           //
	EventNames() []string
	RemoveListener(string, TListener)
	RemoveAllListeners(...string)
	Off(string, TListener)
	ListenerCount(string) int
	Listeners(string) []TListener
	GetMaxListeners() int
	SetMaxListeners(int)
	PrependListener(string, core.T /* TListener */) error     //
	PrependOnceListener(string, core.T /* TListener */) error //
	RawListeners(string) []TListener
}

// Events
type event struct {
	count  int
	name   string
	max    int
	events map[string][]TListener
}

// TListener
type TListener struct {
	Id   string
	Fn   core.T
	Once bool
}

type _uuid_ [16]byte

//----------------------------------
// privado para los eventos
//----------------------------------
var (
	dispatcher IEventEmitter
	state      string
	Debug      bool
)

//--------------------
// type _uuid_
//--------------------
func uuid() (u *_uuid_) {
	u = new(_uuid_)
	// Set all bits to randomly (or pseudo-randomly) chosen values.
	rand.Read(u[:])
	u[8] = (u[8] | 0x40) & 0x7F    // setVariant - 0x40
	u[6] = (u[6] & 0xF) | (4 << 4) // setVersion - 4
	return
}

// Retorna desparseada versión de la UUID secuencia.
func (u *_uuid_) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

//--------------------------

func _dispatcher() IEventEmitter {
	// lo obtiene de reserva para hacer cosas con él
	_dispatcher := EventEmitter("IEventEmitter", false)
	_dispatcher.On("Loading.System", func(who, name string) error {
		s := fmt.Sprintf("Evento '%s.%s'. Cargado.", who, name)
		log.Info(s)
		state = "Mounted"
		return nil
	})
	_dispatcher.On("Creating.System", func(who, name string) error {
		s := fmt.Sprintf("Evento '%s.%s'. Dado de ALTA.", who, name)
		log.Warning(s)
		return nil
	})
	_dispatcher.On("Removing.System", func(who, name string) error {
		s := fmt.Sprintf("Evento '%s.%s'. Dado de BAJA.", who, name)
		log.Warning(s)
		return nil
	})
	_dispatcher.On("Accessing.System", func(who, name string) error {
		s := fmt.Sprintf("Evento '%s.%s'. AccedIdo.", who, name)
		log.Info(s)
		return nil
	})
	_dispatcher.Emit("Loading.System", "IEventEmitter", "Loading.System")

	return _dispatcher
}

func _events(v *event) {
	// eventos Count, UnCount, End e Init
	v.On("__Count__", func() error {
		v.count++

		return nil
	})
	v.On("__UnCount__", func() error {
		v.count--

		return nil
	})
	v.On("__End__", func() error {
		s := v.EventNames()
		for _, nameIn := range s {
			v.RemoveAllListeners(nameIn)
		}
		Debug = false
		v.count = 0

		return nil
	})
	v.On("__Init__", func(name string) error {
		if Debug {
			fmt.Printf("Inicio de handler '%s'\n", name)
		}

		return nil
	})

	v.count++
	v.Emit("__Init__", v.name)
}

// EventEmitter retorna un nuevo IEventEmitter
// @param {string} name
// @param {bool} dsp
// @return {IEventEmitter}
// @constructor
func EventEmitter(name string, dsp bool) IEventEmitter {
	v := &event{
		count:  0,
		name:   name,
		max:    _MAX_,
		events: make(map[string][]TListener),
	}
	_events(v)

	if dsp {
		dispatcher = _dispatcher()
	}

	return v
}

//-------------------
func stringToPtr(s string) *string {
	return &s
}

func (e *event) _emit(event string, name string) {
	dispatcher.Emit(event+".System", e.name, name)
	fmt.Print(&log.Buf)
	log.Buf.Reset()
}

// Retorna el número de TListeners escuchando a un evento dado
// @param {string} name
// @return {int}
func (e *event) ListenerCount(name string) int {
	return len(e.events[name])
}

// Obtiene el máximo número de TListener para un evento
// @return {int}
func (e *event) GetMaxListeners() int {
	return e.max
}

// Pone el máximo número de TListener para un evento
// @param {int}
func (e *event) SetMaxListeners(i int) {
	if i <= _MAX_ {
		e.max = _MAX_
	} else {
		e.max = i
	}
}

// Devuelve un array de los TListener de un evento
// @param {string} name
// @return {core.T}
func (e *event) Listeners(name string) []TListener {
	TListener := make([]TListener, 0, len(e.events))

	for _, v := range e.events[name] {
		TListener = append(TListener, v)
	}

	return TListener
}

// Devuelve un array de los TListener de un evento
// @param {string} name
// @return {core.T}
func (e *event) RawListeners(name string) []TListener {
	return e.Listeners(name)
}

// On coloca un nuevo Evento y TListener en dicho evento
// @param {string} name
// @param {core.T} Fn
// @return error
func (e *event) On(name string, Fn core.T /* TListener */) error {
	return e._on(name, false, Fn)
}

// PrependListener coloca un nuevo TListener en la lista de escucha
// de un evento, en la primera posición de la lista de oyentes.
// @param {string}
// @param {core.T}
// @return error
func (e *event) PrependListener(name string, Fn core.T /* TListener */) error {
	return e._on(name, true, Fn)
}

// PrependListener coloca un nuevo TListener en la lista de escucha
// de un evento, en la primera posición de la lista de oyentes.
// @param {string}
// @param {core.T}
// @return error
func (e *event) PrependOnceListener(name string, Fn core.T /* TListener */) error {
	err := e._on(name, true, Fn)
	e.events[name][0].Once = true

	return err
}

// helper para _on
func (e *event) _onEnd(name string, prepend bool, Fn core.T) error {
	lst := TListener{Fn: Fn, Once: false, Id: (uuid()).String()}
	if !e.has(name) {
		e.Emit("__Count__")
	}
	e.Emit("newListener", name, lst.Id)
	// avisa el warning del máximo de TListeners
	if e.ListenerCount(name) > (e.max - 1) {
		var bf bytes.Buffer
		log.Create(&bf, "-NOTA: ", log.I)("Se añaden más de 5 listeners (oyentes), al evento '" + name + "'")
		fmt.Print(&bf)
		bf.Reset()
	}
	//----------
	if prepend {
		// anteposición
		back := make([]TListener, 0, len(e.events[name]))
		back = append(back, e.events[name]...)
		delete(e.events, name)
		e.events[name] = append(e.events[name], lst)
		e.events[name] = append(e.events[name], back...)
	} else {
		// creación
		e.events[name] = append(e.events[name], lst)
	}
	if Debug {
		e._emit("Creating", name)
	}

	return nil
}

var messageA *string = stringToPtr("Fn: firma no es igual")

// helper para On y PrependTListener
func (e *event) _on(name string, prepend bool, Fn core.T) error {
	if Fn == nil {
		return errors.New("Fn: es nil")
	}
	if _, ok := Fn.(handle); ok {
		return e._onEnd(name, prepend, Fn) // hace la creación y anteposición
	}

	t := reflect.TypeOf(Fn)
	if t.Kind() != reflect.Func {
		return errors.New("Fn: no es una función")
	}
	if t.NumOut() != 1 {
		return errors.New("Fn: debe tener un valor de retorno")
	}
	if t.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return errors.New("Fn: debe retornar un mensaje de error")
	}
	if list, ok := e.events[name]; ok && len(list) > 0 {
		tt := reflect.TypeOf(list[0].Fn)
		if tt.NumIn() != t.NumIn() {
			return errors.New(*messageA)
		}
		for i := 0; i < tt.NumIn(); i++ {
			if tt.In(i) != t.In(i) {
				return errors.New(*messageA)
			}
		}
	}

	return e._onEnd(name, prepend, Fn) // hace la creación y anteposición
}

func (e *event) _des(lst TListener) (Id string, Fn core.T, Once bool) {
	return lst.Id, lst.Fn, lst.Once
}

//Once añade un TListener en un evento para ejecutarlo una sóla vez.
//Luego lo borra (borra el evento)
// @param {string} name
// @param {core.T} Fn
// @return error
func (e *event) Once(name string, fn core.T /* TListener */) error {
	v := e.On(name, fn)
	e.events[name][len(e.events[name])-1].Once = true

	return v
}

// AddTListener coloca un nuevo TListener para un evento
// @param {string}
// @param {core.T}
// @return error
func (e *event) AddListener(s string, Fn core.T /* TListener */) error {
	return e._on(s, false, Fn)
}

var preventDefault bool = true

// Emit lanza un evento. Si se ha definIdo como Once, lo
// borra después de ejecutarlo.
// @param {string} name
// @param {...core.T} parámetros
// @return error
func (e *event) Emit(name string, params ...core.T) error {
	listener /* TListener */ := e.events[name]
	for i := 0; i <= (len(listener) - 1); i++ { // lista por orden
		if Debug {
			varCopy := Debug
			Debug = false
			e._emit("Accessing", name)
			Debug = varCopy
		}

		var err error
		stopped := false
		err = nil
		if preventDefault {
			stopped, err = e.call(listener[i].Fn, params...)
		}

		// mira si es del tipo Once y lo borra
		if len(e.events[name]) > 0 && e.events[name][i].Once {
			e.RemoveListener(name, listener[i])
			i--
			listener = e.events[name]
		}

		//error/stopped
		if err != nil {
			return err
		}
		if stopped {
			break
		}
	}

	stopPropagation = false
	preventDefault = true

	return nil
}

var messageB *string = stringToPtr("parametros no coincIdentes")
var stopPropagation bool

// call helper to Emit
// llama de dos maneras:
// A) busca si es la manera normal, o sea, un evento tal cual
// B) busca si La firma es 'func(Event) error', esta descompuesto en una o dos  fn's (en su caso)
// una que gestiona la respuesta (maneja los datos y los prepara, etc) y otra que la manipula (handle).
// 	Asì se podrá regular la propagación del evento, cancelación, etc.
func (e *event) call(Fn core.T, params ...core.T) (stopped bool, err error) {
	if _, ok := Fn.(bool); ok {
		return
	}
	// si es fn handle. B)-------------
	if f, ok := Fn.(handle); ok {
		if len(params) != 1 {
			return stopped, errors.New(*messageB)
		}
		event, ok := (params[0]).(IEvent)
		if !ok {
			return stopped, errors.New(*messageB)
		}
		err = f(event) // llama al handle del evento
		return stopPropagation, err
	}

	// si es A) ------------------
	var (
		f     = reflect.ValueOf(Fn)
		t     = f.Type()
		numIn = t.NumIn()
		in    = make([]reflect.Value, 0, numIn)
	)

	// la fn es variadica
	if t.IsVariadic() { // varidica
		n := numIn - 1
		if len(params) < n {
			return stopped, errors.New(*messageB)
		}
		for _, param := range params[:n] {
			in = append(in, reflect.ValueOf(param))
		}
		s := reflect.MakeSlice(t.In(n), 0, len(params[n:]))
		for _, param := range params[n:] {
			s = reflect.Append(s, reflect.ValueOf(param))
		}
		in = append(in, s)

		// EJECUTA la Fn
		err, _ = f.CallSlice(in)[0].Interface().(error)
		return stopped, err
	}

	// la fn es normal
	if len(params) != numIn {
		return stopped, errors.New(*messageB)
	}
	for _, param := range params {
		in = append(in, reflect.ValueOf(param))
	}

	// EJECUTA la Fn
	err, _ = f.Call(in)[0].Interface().(error)

	return stopped, err
}

// Has retorna true si un evento existe
//
// @param {string} name
// @return bool
func (e *event) has(name string) bool {
	_, ok := e.events[name]

	return ok
}

// EventNames retorna lista de eventos
// @return []string
func (e *event) EventNames() []string {
	list := make([]string, 0, len(e.events))

	for name := range e.events {
		list = append(list, name)
	}

	return list
}

// RemoveTListener borra TListener de la lista de oyentes de un evento
// @param {...string}
func (e *event) RemoveListener(name string, iWho TListener) {
	var named string

	if v, ok := e.events[name]; ok {
		delete(e.events, name)
		for i := 0; i <= (len(v) - 1); i++ {
			if v[i].Id == iWho.Id {
				named = v[i].Id
			} else {
				e.events[name] = append(e.events[name], v[i])
			}
		}
	}

	e.Emit("removeListener", named)

	if len(e.events[name]) == 0 {
		e.RemoveAllListeners(name)
	}
}

// RemoveAllListeners borra evento de la lista de eventos.
// Si no se pasan argumentos remueve todos los eventos.
// @param {...string}
func (e *event) RemoveAllListeners(names ...string) {
	if len(names) > 0 {
		for _, name := range names {
			delete(e.events, name)
			if Debug {
				e._emit("Removing", name)
			}
			e.Emit("__UnCount__")
		}

		return
	}

	varCopy := Debug
	e.Emit("__End__")
	Debug = varCopy
	e.count = 0
	e.events = make(map[string][]TListener)
}

// Off borra listeners de la lista de listeners
// @param {...string}
func (e *event) Off(name string, idx TListener) {
	e.RemoveListener(name, idx)
}
