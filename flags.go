package appctx

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type appFlag struct {
	names       []string
	description string
	def         any
	value       any
}

func (app *AppCtx[_, _]) Flag(name string, value, def any, description string) {
	app.newFlag([]string{name}, value, def, description)
}

func (app *AppCtx[_, _]) Flag2(shortName, longName string, value, def any, description string) {
	app.newFlag([]string{shortName, longName}, value, def, description)
}

func (app *AppCtx[_, _]) newFlag(names []string, value any, def any, description string) {
	app.flags = append(app.flags, appFlag{
		names:       names,
		description: description,
		def:         def,
		value:       value,
	})
}

func (app *AppCtx[_, _]) initFlags() error {
	if app.noFlags {
		return nil
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	for _, f := range app.flags {
		for _, name := range f.names {
			switch v := f.value.(type) {
			case *string:
				def, ok := f.def.(string)
				if !ok {
					return fmt.Errorf("invalid default value type of flag %v: %T (should be %T)", name, f.def, v)
				}

				fs.StringVar(v, name, def, "")
			case *int:
				def, ok := f.def.(int)
				if !ok {
					return fmt.Errorf("invalid default value type of flag %v: %T (should be %T)", name, f.def, v)
				}

				fs.IntVar(v, name, def, "")
			case *bool:
				def, ok := f.def.(bool)
				if !ok {
					return fmt.Errorf("invalid default value type of flag %v: %T (should be %T)", name, f.def, v)
				}

				fs.BoolVar(v, name, def, "")
			}
		}
	}

	fs.Usage = func() {
		fmt.Println(app.title + " v" + app.version + "\n" +
			"Usage:\n" +
			app.getFlagHelp())
	}

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return fmt.Errorf("parsing flags: %w", err)
	}

	return nil
}

func (app *AppCtx[_, _]) getFlagHelp() string {
	s := ""

	for _, f := range app.flags {
		prefixedNames := []string{}
		for _, name := range f.names {
			prefixedNames = append(prefixedNames, map[bool]string{false: "-", true: "--"}[len(name) > 1]+name)
		}

		s += "\n\t" + strings.Join(prefixedNames, ", ") + ": " + f.description
	}

	return s
}
