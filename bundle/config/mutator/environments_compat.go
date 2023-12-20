package mutator

import (
	"context"
	"fmt"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/libs/config"
	cv "github.com/databricks/cli/libs/config"
)

type environmentsToTargets struct{}

func EnvironmentsToTargets() bundle.Mutator {
	return &environmentsToTargets{}
}

func (m *environmentsToTargets) Name() string {
	return "EnvironmentsToTargets"
}

func (m *environmentsToTargets) Apply(ctx context.Context, b *bundle.Bundle) error {
	// Short circuit if the "environments" key is not set.
	// This is the common case.
	if b.Config.Environments == nil {
		return nil
	}

	// The "environments" key is set; validate and rewrite it to "targets".
	return b.Config.Mutate(func(v config.Value) (config.Value, error) {
		environments := v.Get("environments")
		targets := v.Get("targets")

		// Return an error if both "environments" and "targets" are set.
		if environments != cv.NilValue && targets != cv.NilValue {
			return cv.NilValue, fmt.Errorf(
				"both 'environments' and 'targets' are specified; only 'targets' should be used. "+
					"Instance of 'environments' found at %s.",
				environments.Location().String(),
			)
		}

		// Rewrite "environments" to "targets".
		if environments != cv.NilValue && targets == cv.NilValue {
			return v.SetKey("targets", environments), nil
		}

		return v, nil
	})
}