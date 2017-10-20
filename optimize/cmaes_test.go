// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"errors"
	"math"
	"testing"

	"gonum.org/v1/gonum/optimize/functions"
)

func TestCmaEsChol(t *testing.T) {
	for i, test := range []struct {
		dim      int
		problem  Problem
		method   *CmaEsChol
		settings *Settings
		good     func(*Result, error) error
	}{
		{
			// Test that can find a small value.
			dim: 10,
			problem: Problem{
				Func: functions.ExtendedRosenbrock{}.Func,
			},
			method: &CmaEsChol{
				StopLogDet: math.NaN(),
			},
			settings: &Settings{
				FunctionThreshold: 0.01,
			},
			good: func(result *Result, err error) error {
				if result.Status != FunctionThreshold {
					return errors.New("result not function threshold")
				}
				if result.F > 0.01 {
					return errors.New("result not sufficiently small")
				}
				return nil
			},
		},
		{
			// Test that can stop when the covariance gets small.
			// For this case, also test that it is really at a minimum.
			dim: 2,
			problem: Problem{
				Func: functions.ExtendedRosenbrock{}.Func,
			},
			method: &CmaEsChol{},
			settings: &Settings{
				FunctionThreshold: math.Inf(-1),
			},
			good: func(result *Result, err error) error {
				if result.Status != MethodConverge {
					return errors.New("result not method converge")
				}
				if result.F > 1e-12 {
					return errors.New("minimum not found")
				}
				return nil
			},
		},
		{
			// Test that population works properly and it stops after a certain
			// number of iterations.
			dim: 3,
			problem: Problem{
				Func: functions.ExtendedRosenbrock{}.Func,
			},
			method: &CmaEsChol{
				Population: 100,
			},
			settings: &Settings{
				FunctionThreshold: math.Inf(-1),
				MajorIterations:   10,
			},
			good: func(result *Result, err error) error {
				if result.Status != IterationLimit {
					return errors.New("result not iteration limit")
				}
				if result.FuncEvaluations != 1000 {
					return errors.New("wrong number of evaluations")
				}
				return nil
			},
		},
		{
			// Test that works properly in parallel, and stops with some
			// number of function evaluations.
			dim: 5,
			problem: Problem{
				Func: functions.ExtendedRosenbrock{}.Func,
			},
			method: &CmaEsChol{
				Population: 100,
			},
			settings: &Settings{
				Concurrent:        5,
				FunctionThreshold: math.Inf(-1),
				FuncEvaluations:   250, // Somewhere in the middle of an iteration.
			},
			good: func(result *Result, err error) error {
				if result.Status != FunctionEvaluationLimit {
					return errors.New("result not function evaluations")
				}
				if result.FuncEvaluations < 250 {
					return errors.New("too few function evaluations")
				}
				if result.FuncEvaluations > 250+4 { // can't guarantee exactly, because could grab extras in parallel first.
					return errors.New("too many function evaluations")
				}
				return nil
			},
		},
	} {
		// Run and check that the expected termination occurs.
		result, err := Global(test.problem, test.dim, test.settings, test.method)
		if testErr := test.good(result, err); testErr != nil {
			t.Errorf("cas %d: %v", i, testErr)
		}

		// Run a second time to make sure there are no residual effects
		result, err = Global(test.problem, test.dim, test.settings, test.method)
		if testErr := test.good(result, err); testErr != nil {
			t.Errorf("cas %d second: %v", i, testErr)
		}
	}
}
