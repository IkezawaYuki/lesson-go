package flag_mirror

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

var ErrHelp = errors.New("flag help requested")

var errParse = errors.New("parse error")

var errRange = errors.New("value out of range")

func numError(err error) error {
	ne, ok := err.(*strconv.NumError)
	if !ok {
		return err
	}
	if ne.Err == strconv.ErrSyntax {
		return errParse
	}
	if ne.Err == strconv.ErrRange {
		return errRange
	}
	return err
}

type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) Get() interface{} {
	return bool(*b)
}

func (b *boolValue) String() string {
	return strconv.FormatBool(bool(*b))
}

func (b *boolValue) IsBoolFlag() bool {
	return true
}

type boolFlag interface {
	flag.Value
	IsBoolFlag() bool
}

type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		err = numError(err)
	}
	*i = intValue(v)
	return err
}

func (i *intValue) Get() interface{} {
	return int(*i)
}

func (i *intValue) String() string {
	return strconv.Itoa(int(*i))
}

type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		err = numError(err)
	}
	*i = int64Value(v)
	return err
}

func (i *int64Value) Get() interface{} {
	return int64(*i)
}

func (i *int64Value) String() string {
	return strconv.FormatInt(int64(*i), 10)
}

type uint64Value uint

func newUint64Value(val uint, p *uint) *uint64Value {
	*p = val
	return (*uint64Value)(p)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		err = numError(err)
	}
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Get() interface{} { return uint64(*i) }

func (i *uint64Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() interface{} {
	return string(*s)
}

func (s *stringValue) String() string {
	return string(*s)
}

type float64Value float64

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		err = numError(err)
	}
	*f = float64Value(v)
	return err
}

func (f *float64Value) Get() interface{} {
	return float64(*f)
}

func (f *float64Value) String() string {
	return strconv.FormatFloat(float64(*f), 'g', -1, 64)
}

type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		err = errParse
	}
	*d = durationValue(v)
	return err
}

func (d *durationValue) Get() interface{} { return time.Duration(*d) }

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

type uintValue uint

func newUintValue(val uint, p *uint) *uintValue {
	*p = val
	return (*uintValue)(p)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	if err != nil {
		err = numError(err)
	}
	*i = uintValue(v)
	return err
}

func (i *uintValue) Get() interface{} {
	return uint(*i)
}

func (i *uintValue) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}

type Value interface {
	String() string
	Set(string) error
}

type Getter interface {
	Value
	Get() interface{}
}

type ErrorHandling int

const (
	ContinueOnError ErrorHandling = iota
	ExitOnError
	PanicOnError
)

type FlagSet struct {
	Usage         func()
	name          string
	parsed        bool
	actual        map[string]*Flag
	formal        map[string]*Flag
	args          []string
	errorHandling ErrorHandling
	output        io.Writer
}

type Flag struct {
	Name     string
	Usage    string
	Value    Value
	DefValue string
}

func sortFlags(flags map[string]*Flag) []*Flag {
	result := make([]*Flag, len(flags))
	i := 0
	for _, f := range flags {
		result[i] = f
		i++
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func (f *FlagSet) Output() io.Writer {
	if f.output == nil {
		return os.Stderr
	}
	return f.output
}

func (f *FlagSet) Name() string {
	return f.name
}

func (f *FlagSet) ErrorHandling() ErrorHandling {
	return f.errorHandling
}

func (f *FlagSet) SetOutput(output io.Writer) {
	f.output = output
}

func (f *FlagSet) VisitAll(fn func(*Flag)) {
	for _, flag := range sortFlags(f.formal) {
		fn(flag)
	}
}

func (f *FlagSet) PrintDefaults() {
	f.VisitAll(func(flag *Flag) {
		s := fmt.Sprintf("-%s", flag.Name)
		name, usage := UnquoteUsage(flag)
		if len(name) > 0 {
			s += " " + name
		}

		// todo
	})
}

func UnquoteUsage(flag *Flag) (name string, usage string) {
	usage = flag.Usage
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name = usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				}
			}
			break
		}
	}
	name = "value"
	switch flag.Value.(type) {
	case boolFlag:
		name = ""
	case *durationValue:
		name = "duration"
	case *float64Value:
		name = "float"
	case *intValue, *int64Value:
		name = "init"
	case *stringValue:
		name = "string"
	case *uintValue, *uint64Value:
		name = "unit"
	}
	return
}

func (f *FlagSet) defaultUsage() {
	if f.name == "" {
		fmt.Fprintf(f.Output(), "Usage:\n")
	} else {
		fmt.Fprintf(f.Output(), "Usage of %s:˜\n", f.name)
	}
	f.PrintDefaults()
}

func VisitAll(fn func(*Flag)) {
	CommandLine.VisitAll(fn)
}

var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

func init() {
	CommandLine.Usage = commandLineUsage
}

var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
}

func commandLineUsage() {
	Usage()
}

func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	f := &FlagSet{
		name:          name,
		errorHandling: errorHandling,
	}
	f.Usage = f.defaultUsage
	return f
}