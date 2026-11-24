// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package config

import (
	"context"
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"github.com/google/go-cmp/cmp"
)

// Upjet treats an unset ControllerReconcileVersion as a version other than
// the generated one and merges the field conversion annotation into the
// Terraform state, which fails for the indexed paths of list fields.
func TestMultiVersionResourcesReconcileGeneratedVersion(t *testing.T) {
	type args struct {
		getProvider func(context.Context) (*ujconfig.Provider, error)
	}
	type want struct {
		reconcileVersions map[string]string
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"Cluster": {
			args: args{getProvider: GetProvider},
			want: want{reconcileVersions: map[string]string{"vultr_load_balancer": "v1beta2"}},
		},
		"Namespaced": {
			args: args{getProvider: GetProviderNamespaced},
			want: want{reconcileVersions: map[string]string{"vultr_load_balancer": "v1beta2"}},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			p, err := tc.args.getProvider(t.Context())
			if err != nil {
				t.Fatalf("getProvider: %s", err)
			}

			got := map[string]string{}
			for name, r := range p.Resources {
				if len(r.PreviousVersions) == 0 {
					continue
				}
				got[name] = r.ControllerReconcileVersion //nolint:staticcheck // upjet does not default it at runtime
			}

			if diff := cmp.Diff(tc.want.reconcileVersions, got); diff != "" {
				t.Errorf("reconcile versions of multi-version resources: -want, +got:\n%s", diff)
			}
		})
	}
}
