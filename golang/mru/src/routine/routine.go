package routine

import (
	"assertion"
	"context"
	"dut"
	"errors"
	"fmt"
	"rut"
	"strings"
)

type Routine struct {
	Name        string `json:"name"`
	Function    string
	Paramerters []string
	Expected    string
	IsAssert    bool
	Assertions  []*assertion.Assertion `json:"assertions"`
	Description string                 `json:"description"`
	Dut         *dut.DUT
}

func (r *Routine) GetParams() (int, []string) {
	res := make([]string, 0, 1)

	if len(r.Paramerters) == 0 {
		return 0, nil
	}

	for _, p := range r.Paramerters {
		p = strings.TrimSpace(p)
		if p != "" {
			res = append(res, p)
		}
	}

	return len(res), res
}

func (r *Routine) String() string {
	var res string

	if r.IsAssert {
		res += fmt.Sprintf("ASSERT ")
	}
	res += fmt.Sprintf("%s.%s(", r.Dut.Name, r.Function)
	for _, p := range r.Paramerters {
		res += fmt.Sprintf("%s, ", p)
	}

	res = strings.Trim(res, ", ")

	if r.IsAssert {
		res += fmt.Sprintf(") = %s", r.Expected)
	} else {
		res += fmt.Sprintf(")")
	}

	return res
}

func (r *Routine) Run(ctx context.Context, db *rut.DB) error {
	for _, a := range r.Assertions {
		msg, ok := a.Do(ctx, db)
		if !ok {
			return errors.New(fmt.Sprintf("Routine: %s failed with: %s", r.Name, msg))
		}
	}

	return nil
}

func (r *Routine) Do() bool {
	_, parms := r.GetParams()
	if r.IsAssert {
		return r.Dut.Assert(r.Function, r.Expected, parms...)
	} else {
		r.Dut.Call(r.Function, parms...)
	}
	return true
}
