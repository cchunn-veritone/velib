package velib

import "net/http"

//     _____      _   _
//    / ____|    | | | |
//   | (___   ___| |_| |_ ___ _ __ ___
//    \___ \ / _ \ __| __/ _ \ '__/ __|
//    ____) |  __/ |_| ||  __/ |  \__ \
//   |_____/ \___|\__|\__\___|_|  |___/
//
//   Externally exposed functions for setting an option

func SetPort(o int) { v.setPort(o) }

// Setters needed only to change a state to true
func Ready()   { v.setReady(true) }
func Failed()  { v.setFailed(true) }
func Verbose() { v.setVerbose(true) }

// Setters without getters don't need Set prefix
func Process(f func(responseWriter http.ResponseWriter, request *http.Request)) { v.setProcess(f) }

//     _____      _   _
//    / ____|    | | | |
//   | |  __  ___| |_| |_ ___ _ __ ___
//   | | |_ |/ _ \ __| __/ _ \ '__/ __|
//   | |__| |  __/ |_| ||  __/ |  \__ \
//    \_____|\___|\__|\__\___|_|  |___/
//
//   Externally Exposed functions to get an option value

func GetReady() bool  { return v.getReady() }
func GetFailed() bool { return v.getFailed() }

//     _____ _                   _     __  __      _   _               _
//    / ____| |                 | |   |  \/  |    | | | |             | |
//   | (___ | |_ _ __ _   _  ___| |_  | \  / | ___| |_| |__   ___   __| |___
//    \___ \| __| '__| | | |/ __| __| | |\/| |/ _ \ __| '_ \ / _ \ / _` / __|
//    ____) | |_| |  | |_| | (__| |_  | |  | |  __/ |_| | | | (_) | (_| \__ \
//   |_____/ \__|_|   \__,_|\___|\__| |_|  |_|\___|\__|_| |_|\___/ \__,_|___/
//
//    Used by the getters to access the options stored in the struct

func (v *Velib) getPort() int  { return v.port }
func (v *Velib) setPort(o int) { v.port = o }

func (v *Velib) getReady() bool  { return v.ready }
func (v *Velib) setReady(o bool) { v.ready = o }

func (v *Velib) getFailed() bool  { return v.failed }
func (v *Velib) setFailed(o bool) { v.failed = o }

func (v *Velib) getVerbose() bool  { return v.verbose }
func (v *Velib) setVerbose(o bool) { v.verbose = o }

func (v *Velib) setProcess(f func(responseWriter http.ResponseWriter, request *http.Request)) {
	v.process = f
}
