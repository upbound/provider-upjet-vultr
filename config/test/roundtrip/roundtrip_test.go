// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

// Package roundtrip contains API roundtrip tests for the multi-version
// managed resource APIs.
package roundtrip

import (
	"maps"
	"reflect"
	"testing"

	"github.com/crossplane/upjet/v2/pkg/apitesting/roundtrip"
	"github.com/crossplane/upjet/v2/pkg/config/conversion"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	clusterapis "github.com/upbound/provider-vultr/apis/cluster"
	namespacedapis "github.com/upbound/provider-vultr/apis/namespaced"
	"github.com/upbound/provider-vultr/config"
)

// equateFieldConversionAnnotation ignores the Upjet field conversion
// annotation when comparing object annotations. The roundtrip library only
// drops the map entry, so an object that started with nil annotations would
// still differ from its round-tripped copy, whose annotations hold just that
// key. Fields without a counterpart in an older API version, such as the
// LoadBalancer global region VPC IDs, are carried in this annotation.
func equateFieldConversionAnnotation() cmp.Option {
	objectMetaType := reflect.TypeFor[metav1.ObjectMeta]()
	return cmp.FilterPath(func(p cmp.Path) bool {
		sf, ok := p.Last().(cmp.StructField)
		return ok && sf.Name() == "Annotations" && p.Index(-2).Type() == objectMetaType
	}, cmp.FilterValues(func(x, y map[string]string) bool {
		// Two empty maps are already equated by the library defaults, and
		// handling them here as well would make the options ambiguous.
		return len(x) != 0 || len(y) != 0
	}, cmp.Comparer(func(x, y map[string]string) bool {
		x, y = maps.Clone(x), maps.Clone(y)
		delete(x, conversion.AnnotationKey)
		delete(y, conversion.AnnotationKey)
		return maps.Equal(x, y)
	})))
}

func TestRoundTrip(t *testing.T) {
	provider, err := config.GetProvider(t.Context())
	if err != nil {
		t.Fatalf("GetProvider: %s", err)
	}
	providerNamespaced, err := config.GetProviderNamespaced(t.Context())
	if err != nil {
		t.Fatalf("GetProviderNamespaced: %s", err)
	}

	testScheme := runtime.NewScheme()
	if err := clusterapis.AddToScheme(testScheme); err != nil {
		t.Fatalf("cluster-scoped apis AddToScheme: %s", err)
	}
	if err := namespacedapis.AddToScheme(testScheme); err != nil {
		t.Fatalf("namespaced apis AddToScheme: %s", err)
	}

	rt, err := roundtrip.NewRoundTripTest(provider, providerNamespaced, testScheme,
		roundtrip.WithFuzzerConfig(
			roundtrip.FuzzerIterations(10),
			roundtrip.FuzzerNilChance(0)),
		roundtrip.WithFuzzerConfig(
			roundtrip.FuzzerIterations(30),
			roundtrip.FuzzerNilChance(0.3)),
		roundtrip.WithComparisonOptions(
			roundtrip.EquateEmptyAndSingleZeroSlice(),
			roundtrip.EquateNilAndZeroValuePtr(),
			equateFieldConversionAnnotation(),
		),
	)
	if err != nil {
		t.Fatalf("NewRoundTripTest: %s", err)
	}

	t.Run("TestSerializationRoundtrip", func(t *testing.T) {
		rt.TestSerializationRoundtrip(t)
	})
	t.Run("TestConversionRoundtrip", func(t *testing.T) {
		rt.TestConversionRoundtrip(t)
	})
}
