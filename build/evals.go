package build

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"reflect"
	"strings"

	shellquote "github.com/kballard/go-shellquote"
)

func evalVars(b *Build, err error) (*Build, error) {
	if err != nil {
		return b, err
	}

	if b.Variables == nil {
		return b, nil
	}

	cfg := b.Variables
	for k, v := range cfg.Eval {
		cfg.Eval[k] = evalCommandSafely(v)
	}

	for k, v := range cfg.Env {
		cfg.Env[k] = evalEnvVarSafely(v)
	}

	b.Variables = cfg

	return evalTemplatable(b, nil)
}

func evalCommandSafely(v string) string {
	cmdWithArgs, err := shellquote.Split(v)
	if err != nil {
		fmt.Println("error processing command, will return empty string, ignored: ", err.Error())
		return ""
	}
	if output, err := exec.Command(cmdWithArgs[0], cmdWithArgs[1:]...).Output(); err == nil {
		return strings.TrimSpace(string(output))
	}
	return ""
}

func evalEnvVarSafely(v string) string {
	return os.Getenv(v)
}

// eval only relevant for DockerBuild
func evalTemplatable(b *Build, err error) (*Build, error) {
	if err != nil {
		return b, err
	}

	D := b.Docker
	if D == nil {
		return b, err
	}

	nv := recurseEvalTemplatable(reflect.ValueOf(D), *b.Variables)

	if d, ok := nv.Interface().(*DockerBuild); ok {
		b.Docker = d
		return b, nil
	}

	return b, fmt.Errorf("cannot get DockerBuild back after eval")
}

func recurseEvalTemplatable(V reflect.Value, vars Variables) reflect.Value {
	rv := reflect.Indirect(V)
	rt := rv.Type()
	if rv.Kind() != reflect.Struct {
		return V
	}

	for i := 0; i < rt.NumField(); i++ {
		cf := rt.Field(i)
		cv := rv.Field(i)
		if isTemplatable(cf) {
			extrapolate(cv, vars)
		}
		if isPtrToStruct(cv) {
			cv.Set(recurseEvalTemplatable(cv, vars))
		}
	}

	return V
}

func isTemplatable(sf reflect.StructField) bool {
	_, ok := sf.Tag.Lookup("templatable")
	isString := sf.Type.Kind() == reflect.String
	isStringSlice := sf.Type.Kind() == reflect.Slice && sf.Type.Elem().Kind() == reflect.String
	return (isString || isStringSlice) && ok
}

func extrapolate(v reflect.Value, vars Variables) {
	var (
		err  error
		bb   bytes.Buffer
		tmpl *template.Template
	)

	if tpl, ok := v.Interface().(string); ok {
		if tmpl, err = template.New("extrapolate").Parse(tpl); err == nil {
			err = tmpl.Execute(&bb, vars)
		}

		if err == nil {
			v.Set(reflect.ValueOf(bb.String()))
			return
		}
	}

	if sliceString, ok := v.Interface().([]string); ok {
		var ns []string
		for _, tpl := range sliceString {
			if tmpl, err = template.New("extrapolate").Parse(tpl); err == nil {
				err = tmpl.Execute(&bb, vars)
			}

			if err == nil {
				ns = append(ns, bb.String())
				bb = bytes.Buffer{}
			} else {
				fmt.Println("error extrapolating template, ignored: ", err.Error())
			}
		}
		v.Set(reflect.ValueOf(ns))
		return
	}

	fmt.Println("error extrapolating template, ignored: ", err.Error())
}

func isPtrToStruct(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct
}
