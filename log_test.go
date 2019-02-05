package ech0

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"

	gommon "github.com/labstack/gommon/log"
)

// zerolog output json
type msgZero struct {
	Time   *time.Time `json:"time"`
	File   string     `json:"file"`
	Level  string     `json:"level"`
	Prefix string     `json:"prefix"`
	Line   int        `json:"line"`
	Msg    string     `json:"message"`
	// fields below are only present when calling the '{Debug|Info|Warn|Error}'j methods
	Foo string                 `json:"foo"`
	Baz []string               `json:"baz"`
	Zaz map[string]interface{} `json:"z"`
}

// gommon output JSON - notice how line number is a string for some reason....
type msgGommon struct {
	Time   *time.Time `json:"time"`
	File   string     `json:"file"`
	Level  string     `json:"level"`
	Prefix string     `json:"prefix"`
	Line   string     `json:"line"`
	Msg    string     `json:"message"`
	// fields below are only present when calling the '{Debug|Info|Warn|Error}'j methods
	Foo string                 `json:"foo"`
	Baz []string               `json:"baz"`
	Zaz map[string]interface{} `json:"z"`
}

func TestLogLevels(t *testing.T) {

	type tc struct {
		name    string
		methods []string
	}

	gom := gommon.New("")
	gomb := &bytes.Buffer{}
	gom.SetOutput(gomb)
	gom.SetLevel(gommon.DEBUG)

	z := New("")
	zb := &bytes.Buffer{}
	z.SetOutput(zb)
	z.SetLevel(gommon.DEBUG)

	gt := reflect.TypeOf(gom)
	zt := reflect.TypeOf(z)

	for _, testcase := range []tc{
		{name: "Debug", methods: []string{"Debug", "Debugf", "Debugj"}},
		{name: "Info", methods: []string{"Info", "Infof", "Infoj"}},
		{name: "Warn", methods: []string{"Warn", "Warnf", "Warnj"}},
		{name: "Error", methods: []string{"Error", "Errorf", "Errorj"}},
		{name: "NoLevel", methods: []string{"Print", "Printf", "Printj"}},
	} {
		t.Run(testcase.name, func(t *testing.T) {

			for i, j := range [][]reflect.Value{
				[]reflect.Value{reflect.ValueOf("hello")},
				[]reflect.Value{reflect.ValueOf("%s"), reflect.ValueOf("bbq")},
				[]reflect.Value{reflect.ValueOf(map[string]interface{}{"foo": "bar", "baz": []string{"a", "b"}, "z": map[string]interface{}{"1": "2"}})},
			} {

				var zbb, gombb []byte
				var gotz msgZero
				var gotm msgGommon
				var zm, gm reflect.Method

				// pre-pend method receiver
				args := make([]reflect.Value, len(j)+1)
				copy(args[1:], j)

				// call the log function using reflect....
				meth := testcase.methods[i]
				t.Logf("testing method %s", meth)

				args[0] = reflect.ValueOf(gom)
				gm, _ = gt.MethodByName(meth)
				gm.Func.Call(args)

				args[0] = reflect.ValueOf(z)
				zm, _ = zt.MethodByName(meth)
				zm.Func.Call(args)

				// ... output is JSON, so parse it...
				gombb = gomb.Bytes()
				t.Logf("gommon output: %s", gombb)
				if err := json.Unmarshal(gombb, &gotm); err != nil {
					t.Error(err)
					goto next
				}

				zbb = zb.Bytes()
				t.Logf("zerolog output: %s", zbb)
				if err := json.Unmarshal(zbb, &gotz); err != nil {
					t.Error(err)
					goto next
				}

				// .... compare outputs (loosely)
				if !equal(t, gotz, gotm) {
					t.Errorf("%s: gommon: %s, zerolog: %s", meth, gombb, zbb)
				}

			next:
				gomb.Reset()
				zb.Reset()

			}
		})

	}
}

func equal(t *testing.T, a msgZero, b msgGommon) bool {
	// this is a 'loose' comparison of log outputs. Some values like timestamps and line numbers won't match exactly
	t.Helper()
	t.Logf("zerolog: %#v, gommon: %#v", a, b)

	bl, err := strconv.Atoi(b.Line)
	if err != nil {
		t.Fatalf("unexpected: '%s' not an int", b.Line)
	}

	if a.Time == nil || b.Time == nil {
		t.Fatal("unexpected: time value is nil")
	}

	linesEq := math.Abs(float64(bl-a.Line)) <= 5
	timesEq := math.Abs(float64(a.Time.Sub(*b.Time))) <= float64(time.Second*2)

	eq := a.File == b.File &&
		a.Prefix == b.Prefix &&
		strings.EqualFold(a.Level, b.Level) &&
		a.Msg == b.Msg &&
		linesEq &&
		timesEq

	if b.Foo != "" {
		eq = eq && a.Foo == b.Foo
	}
	if b.Zaz != nil {
		eq = eq && reflect.DeepEqual(a.Zaz, b.Zaz)
	}
	if b.Baz != nil {
		eq = eq && reflect.DeepEqual(b.Baz, a.Baz)
	}

	return eq
}

func TestMisc(t *testing.T) {
	l := New("")
	l.SetLevel(gommon.INFO)
	l.SetPrefix("hello")
	l.Info("hello")
	l.SetLevel(gommon.WARN)
	l.Warn("hello", "again")

	g := gommon.New("")
	g.Warn("hello", "again")
}

func BenchmarkZeroFormat(b *testing.B) {
	benchFormat(New(""), b)
}
func BenchmarkZeroJSON(b *testing.B) {
	benchJSON(New(""), b)
}

func BenchmarkZero(b *testing.B) {
	bench(New(""), b)
}
func BenchmarkGommonFormat(b *testing.B) {
	benchFormat(gommon.New(""), b)

}
func BenchmarkGommonJSON(b *testing.B) {
	benchJSON(gommon.New(""), b)
}

func BenchmarkGommon(b *testing.B) {
	bench(gommon.New(""), b)
}

func benchJSON(l echo.Logger, b *testing.B) {
	b.Helper()
	b.ReportAllocs()
	l.SetOutput(ioutil.Discard)
	j := map[string]interface{}{"foo": "bar"}

	for i := 0; i < b.N; i++ {
		l.Infoj(j)
	}
}
func benchFormat(l echo.Logger, b *testing.B) {
	b.Helper()
	b.ReportAllocs()
	l.SetOutput(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		l.Infof("%s", "hello")
	}
}

func bench(l echo.Logger, b *testing.B) {
	b.Helper()
	b.ReportAllocs()
	l.SetOutput(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		l.Info("hello")
	}
}
