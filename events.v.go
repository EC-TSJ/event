// antiguo nombre default.v
// Package event is a simple event system.
package event

var (
	_MAX_ int = 5
)

/*

func Once(d IEventEmitter, eventName string, options ...core.T) (av []string) {
	v := d.(*event)
	if z, ok := v.events[eventName]; ok {
		for a := 0; a < len(z); a++ {
			av = append(av, z[a].Id)
		}
	}

	return
}

func On(d IEventEmitter, eventName string, options ...core.T) (lista []string) {
	v := d.(*event)
	if z, ok := v.events[eventName]; ok {
		for a := 0; a < len(z); a++ {
			lista = append(lista, z[a].Id)
		}
	}

	return
}

*/

// Pone el maxListeners en general
// @param {int}
func DefaultMaxListeners(i int) {
	if i <= 5 {
		_MAX_ = 5
	} else {
		_MAX_ = i
	}
}

// Devuelve un array de los Listener de un evento
// @param {IEventEmitter}
// @param {string}
// @return {[]TListener}
func GetEventListeners(d IEventEmitter, eventName string) []TListener {
	return d.Listeners(eventName)
}

// Obtiene una lista de eventos@id (identificadores) de listeners
// @param {IEventEmitter}
// @param {string}
// @return {[]string}
func ListenerId(d IEventEmitter, eventName string) (lista []string) {
	if z, ok := d.(*event).events[eventName]; ok {
		for a := 0; a < len(z); a++ {
			lista = append(lista, eventName+"@"+z[a].Id)
		}
	}

	return
}

// Nos dice si existe como un evento
// @param {IEventEmitter}
// @param {string}
// @return {bool}
func Has(d IEventEmitter, eventName string) bool {
	return d.(*event).has(eventName)
}

// Obtiene el numero de listeners por evento
// @param {IEventEmitter}
// @param {string}
// @return {int}
func ListenerCount(d IEventEmitter, eventName string) int {
	return d.(*event).ListenerCount(eventName)
}

// Pone setMaxListeners para cada EventEmitter
// @param {int}
// @param {...IEventEmitter}
func SetMaxListeners(n int, dsp ...IEventEmitter) {
	if n <= 5 {
		n = 5
	}

	for _, v := range dsp {
		v.(*event).max = n
	}
}
