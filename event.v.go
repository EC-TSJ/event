// Package event is a simple event system.
package event

type (
	// Eventer interface
	IEvent interface {
		StopPropagation()
		StopInmediatePropagation()
		Cancelable(bool)
		PreventDefault()
		DefaultPrevented() bool
		ReturnValue() bool
	}

	// Event es la clase base para las clases que contienen datos de eventos.
	TEventer struct {
		stop      bool
		cancel    bool
		prevented bool
	}

	// handle aliase
	handle = func(IEvent) error
)

// true si no cancelado
//
// @return {bool}
func (e *TEventer) ReturnValue() bool {
	return !e.cancel
}

// StopPropagation Detiene la propagación del evento a más oyentes de eventos.
func (e *TEventer) StopPropagation() {
	e.stop = true
	stopPropagation = true
}

// Lo mismo que StopPropagation
func (e *TEventer) StopInmediatePropagation() {
	e.stop = true
	stopPropagation = true
}

// Cancelable pone el evento actual a true ó false
//
// @param {bool}
func (e *TEventer) Cancelable(f bool) {
	e.prevented = false
	e.cancel = f
}

// PreventDefault cancela el evento (da órdenes para cancelar el evento. true ó false).
// Detiene la ejecución del evento
func (e *TEventer) PreventDefault() {
	if e.cancel {
		preventDefault = false
		e.prevented = true
	} else {
		preventDefault = true
		e.prevented = false
	}
}

// DefaultPrevented nos dice si es camcelable y se ha llamado a PreventDefault
//
// @return {bool}
func (e *TEventer) DefaultPrevented() bool {
	return e.prevented
}
