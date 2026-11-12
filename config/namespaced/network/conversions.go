// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package network

import (
	"fmt"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/fieldpath"
	"github.com/crossplane/upjet/v2/pkg/config/conversion"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	versionV1Beta1 = "v1beta1"
	versionV1Beta2 = "v1beta2"

	fieldGlobalRegions = "globalRegions"
	fieldRegionID      = "regionId"
	fieldVPCID         = "vpcId"
)

type globalRegionsConversion struct {
	sourceVersion string
	targetVersion string
	convert       func(regions []any) ([]any, error)
}

func (c *globalRegionsConversion) Applicable(src, dst runtime.Object) bool {
	return src.GetObjectKind().GroupVersionKind().Version == c.sourceVersion &&
		dst.GetObjectKind().GroupVersionKind().Version == c.targetVersion
}

func (c *globalRegionsConversion) ConvertPaved(src, target *fieldpath.Paved) (bool, error) {
	if !c.Applicable(&unstructured.Unstructured{Object: src.UnstructuredContent()},
		&unstructured.Unstructured{Object: target.UnstructuredContent()}) {
		return false, nil
	}
	converted := false
	for _, prefix := range conversion.DefaultPathPrefixes() {
		path := fmt.Sprintf("%s.%s", prefix, fieldGlobalRegions)
		v, err := src.GetValue(path)
		if fieldpath.IsNotFound(err) {
			continue
		}
		if err != nil {
			return false, errors.Wrapf(err, "cannot get %q from the conversion source object", path)
		}
		regions, ok := v.([]any)
		if !ok {
			return false, errors.Errorf("%q of the conversion source object must be a list, got %T", path, v)
		}
		out, err := c.convert(regions)
		if err != nil {
			return false, errors.Wrapf(err, "cannot convert %q", path)
		}
		if err := target.SetValue(path, out); err != nil {
			return false, errors.Wrapf(err, "cannot set %q on the conversion target object", path)
		}
		converted = true
	}
	return converted, nil
}

func regionIDsToObjects(regions []any) ([]any, error) {
	out := make([]any, 0, len(regions))
	for _, r := range regions {
		id, ok := r.(string)
		if !ok {
			return nil, errors.Errorf("region must be a string, got %T", r)
		}
		out = append(out, map[string]any{fieldRegionID: id})
	}
	return out, nil
}

// objectsToRegionIDs keeps one entry per global region, even when it has no
// region ID, so that the VPC IDs preserved in the field conversion annotation
// stay aligned with their regions on the way back.
func objectsToRegionIDs(regions []any) ([]any, error) {
	out := make([]any, 0, len(regions))
	for _, r := range regions {
		m, ok := r.(map[string]any)
		if !ok {
			return nil, errors.Errorf("region must be an object, got %T", r)
		}
		id, _ := m[fieldRegionID].(string)
		out = append(out, id)
	}
	return out, nil
}

// loadBalancerGlobalRegionsConversions returns the conversions between the v1beta1
// LoadBalancer API, whose globalRegions is a set of region IDs, and the
// v1beta2 API, whose globalRegions is a list of objects holding a region ID
// and a VPC ID. The VPC IDs have no v1beta1 counterpart and are carried in
// the field conversion annotation while an object is served as v1beta1.
func loadBalancerGlobalRegionsConversions() []conversion.Conversion {
	prefixes := conversion.DefaultPathPrefixes()
	convs := make([]conversion.Conversion, 0, 3+2*len(prefixes))
	convs = append(convs,
		conversion.NewIdentityConversionExpandPaths(conversion.AllVersions, conversion.AllVersions, prefixes, fieldGlobalRegions),
		&globalRegionsConversion{sourceVersion: versionV1Beta1, targetVersion: versionV1Beta2, convert: regionIDsToObjects},
		&globalRegionsConversion{sourceVersion: versionV1Beta2, targetVersion: versionV1Beta1, convert: objectsToRegionIDs},
	)
	for _, prefix := range prefixes {
		path := fmt.Sprintf("%s.%s[*].%s", prefix, fieldGlobalRegions, fieldVPCID)
		convs = append(convs,
			conversion.NewNewlyIntroducedFieldConversion(versionV1Beta1, versionV1Beta2, path, conversion.FromAnnotation),
			conversion.NewNewlyIntroducedFieldConversion(versionV1Beta2, versionV1Beta1, path, conversion.ToAnnotation),
		)
	}
	return convs
}
