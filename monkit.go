package main

import (
	"io"
	"net/http"

	"gopkg.in/spacemonkeygo/monkit.v2"
	"gopkg.in/spacemonkeygo/monkit.v2/environment"
	"gopkg.in/spacemonkeygo/monkit.v2/present"
)

var mon = monkit.Package()

func serve(addr string) error {
	environment.Register(monkit.Default)
	http.Handle("/monkit/", http.StripPrefix("/monkit/", present.HTTP(monkit.Default)))

	http.HandleFunc("/", handler)
	return http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer mon.Task()(&ctx, w, r)(nil)
	r = r.WithContext(ctx)

	work(w)
}

// FromRequest takes a registry (usually the Default registry), an incoming
// path, and optional query parameters, and returns a Result if possible.
//
// FromRequest understands the following paths:
//  * /ps, /ps/text       - returns the result of SpansText
//  * /ps/dot             - returns the result of SpansDot
//  * /ps/json            - returns the result of SpansJSON
//  * /funcs, /funcs/text - returns the result of FuncsText
//  * /funcs/dot          - returns the result of FuncsDot
//  * /funcs/json         - returns the result of FuncsJSON
//  * /stats, /stats/text - returns the result of StatsText
//  * /stats/json         - returns the result of StatsJSON
//  * /trace/svg          - returns the result of TraceQuerySVG
//  * /trace/json         - returns the result of TraceQueryJSON
//
// The last two paths are worth discussing in more detail, as they take
// query parameters. All trace endpoints require at least one of the following
// two query parameters:
//  * regex    - If provided, the very next Span that crosses a Func that has
//               a name that matches this regex will start a trace until that
//               triggering Span ends, provided the trace_id matches.
//  * trace_id - If provided, the very next Span on a trace with the given
//               trace id will start a trace until the triggering Span ends,
//               provided the regex matches. NOTE: the trace_id will be parsed
//               in hex.
// By default, regular expressions are matched ahead of time against all known
// Funcs, but perhaps the Func you want to trace hasn't been observed by the
// process yet, in which case the regex will fail to match anything. You can
// turn off this preselection behavior by providing preselect=false as an
// additional query param. Be advised that until a trace completes, whether
// or not it has started, it adds a small amount of overhead (a comparison or
// two) to every monitored function.

func printMetrics(addr string, w io.Writer) error {
	// stats provides runtime statistics
	w.Write([]byte("/monkit/stats:\n"))
	resp, err := http.Get("http://" + addr + "/monkit/stats")
	if err != nil {
		return err
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	// funcs provides stats on handler itself (including count)
	w.Write([]byte("\n/monkit/funcs:\n"))
	resp, err = http.Get("http://" + addr + "/monkit/funcs")
	if err != nil {
		return err
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}
