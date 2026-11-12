// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package network

import (
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/fieldpath"
	"github.com/crossplane/upjet/v2/pkg/config/conversion"
	"github.com/google/go-cmp/cmp"
)

const apiGroup = "network.vultr.m.upbound.io/"

func loadBalancer(version string, annotations map[string]any, forProvider map[string]any) map[string]any {
	metadata := map[string]any{"name": "example"}
	if annotations != nil {
		metadata["annotations"] = annotations
	}
	return map[string]any{
		"apiVersion": apiGroup + version,
		"kind":       "LoadBalancer",
		"metadata":   metadata,
		"spec":       map[string]any{"forProvider": forProvider},
	}
}

func TestLoadBalancerGlobalRegions(t *testing.T) {
	vpcAnnotation := map[string]any{
		conversion.AnnotationKey: `{"spec.forProvider.globalRegions[1].vpcId":"vpc-1"}`,
	}

	type args struct {
		src           map[string]any
		targetVersion string
	}
	type want struct {
		target map[string]any
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"RegionIDsBecomeObjects": {
			args: args{
				src:           loadBalancer(versionV1Beta1, nil, map[string]any{"globalRegions": []any{"ams", "fra"}}),
				targetVersion: versionV1Beta2,
			},
			want: want{
				target: loadBalancer(versionV1Beta2, nil, map[string]any{"globalRegions": []any{
					map[string]any{"regionId": "ams"},
					map[string]any{"regionId": "fra"},
				}}),
			},
		},
		"ObjectsBecomeRegionIDsAndVPCIDsMoveToAnnotation": {
			args: args{
				src: loadBalancer(versionV1Beta2, nil, map[string]any{"globalRegions": []any{
					map[string]any{"regionId": "ams"},
					map[string]any{"regionId": "fra", "vpcId": "vpc-1"},
				}}),
				targetVersion: versionV1Beta1,
			},
			want: want{
				target: loadBalancer(versionV1Beta1, vpcAnnotation, map[string]any{"globalRegions": []any{"ams", "fra"}}),
			},
		},
		"VPCIDsAreRestoredFromAnnotation": {
			args: args{
				src:           loadBalancer(versionV1Beta1, vpcAnnotation, map[string]any{"globalRegions": []any{"ams", "fra"}}),
				targetVersion: versionV1Beta2,
			},
			want: want{
				target: loadBalancer(versionV1Beta2, nil, map[string]any{"globalRegions": []any{
					map[string]any{"regionId": "ams"},
					map[string]any{"regionId": "fra", "vpcId": "vpc-1"},
				}}),
			},
		},
		"MissingGlobalRegionsIsLeftUnset": {
			args: args{
				src:           loadBalancer(versionV1Beta1, nil, map[string]any{"region": "ewr"}),
				targetVersion: versionV1Beta2,
			},
			want: want{
				target: loadBalancer(versionV1Beta2, nil, map[string]any{}),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			src := fieldpath.Pave(tc.args.src)
			target := fieldpath.Pave(loadBalancer(tc.args.targetVersion, nil, map[string]any{}))
			for _, c := range loadBalancerGlobalRegionsConversions() {
				pc, ok := c.(conversion.PavedConversion)
				if !ok {
					continue
				}
				if _, err := pc.ConvertPaved(src, target); err != nil {
					t.Fatalf("ConvertPaved(...): unexpected error: %v", err)
				}
			}
			if diff := cmp.Diff(tc.want.target, target.UnstructuredContent()); diff != "" {
				t.Errorf("ConvertPaved(...): -want, +got:\n%s", diff)
			}
		})
	}
}
