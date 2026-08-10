package clients

import (
	"strings"
	"testing"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	"github.com/google/go-cmp/cmp"

	namespacedv1beta1 "github.com/upbound/provider-vultr/apis/namespaced/v1beta1"
)

func TestResolveNamespacedSpec(t *testing.T) {
	type args struct {
		spec      namespacedv1beta1.NamespacedProviderConfigSpec
		namespace string
	}

	cases := map[string]struct {
		args args
		want namespacedv1beta1.ProviderConfigSpec
	}{
		"SecretRefResolvesToReferencerNamespace": {
			args: args{
				namespace: "team-a",
				spec: namespacedv1beta1.NamespacedProviderConfigSpec{
					Credentials: namespacedv1beta1.NamespacedProviderCredentials{
						Source: xpv2.CredentialsSourceSecret,
						SecretRef: &xpv2.LocalSecretKeySelector{
							LocalSecretReference: xpv2.LocalSecretReference{Name: "creds"},
							Key:                  "credentials",
						},
					},
				},
			},
			want: namespacedv1beta1.ProviderConfigSpec{
				Credentials: namespacedv1beta1.ProviderCredentials{
					Source: xpv2.CredentialsSourceSecret,
					CommonCredentialSelectors: xpv2.CommonCredentialSelectors{
						SecretRef: &xpv2.SecretKeySelector{
							SecretReference: xpv2.SecretReference{Name: "creds", Namespace: "team-a"},
							Key:             "credentials",
						},
					},
				},
			},
		},
		"NilSecretRefStaysNil": {
			args: args{
				namespace: "team-a",
				spec: namespacedv1beta1.NamespacedProviderConfigSpec{
					Credentials: namespacedv1beta1.NamespacedProviderCredentials{
						Source: xpv2.CredentialsSourceInjectedIdentity,
					},
				},
			},
			want: namespacedv1beta1.ProviderConfigSpec{
				Credentials: namespacedv1beta1.ProviderCredentials{
					Source: xpv2.CredentialsSourceInjectedIdentity,
				},
			},
		},
		"FsEnvAndScalarFieldsPassThrough": {
			args: args{
				namespace: "team-a",
				spec: namespacedv1beta1.NamespacedProviderConfigSpec{
					RateLimit:  new(int64(600)),
					RetryLimit: new(int64(3)),
					Credentials: namespacedv1beta1.NamespacedProviderCredentials{
						Source: xpv2.CredentialsSourceFilesystem,
						Fs:     &xpv2.FsSelector{Path: "/creds"},
						Env:    &xpv2.EnvSelector{Name: "VULTR_CREDS"},
					},
				},
			},
			want: namespacedv1beta1.ProviderConfigSpec{
				RateLimit:  new(int64(600)),
				RetryLimit: new(int64(3)),
				Credentials: namespacedv1beta1.ProviderCredentials{
					Source: xpv2.CredentialsSourceFilesystem,
					CommonCredentialSelectors: xpv2.CommonCredentialSelectors{
						Fs:  &xpv2.FsSelector{Path: "/creds"},
						Env: &xpv2.EnvSelector{Name: "VULTR_CREDS"},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := resolveNamespacedSpec(tc.args.spec, tc.args.namespace)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("resolveNamespacedSpec() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestBuildConfiguration(t *testing.T) {
	type args struct {
		creds      map[string]string
		rateLimit  *int64
		retryLimit *int64
	}

	cases := map[string]struct {
		args            args
		want            map[string]any
		wantErrContains string
	}{
		"APIKeySuccess": {
			args: args{creds: map[string]string{"api_key": "vultr-key"}},
			want: map[string]any{"api_key": "vultr-key"},
		},
		"EmptyAPIKey": {
			args:            args{creds: map[string]string{"api_key": ""}},
			wantErrContains: `credentials secret has no "api_key" key`,
		},
		"MissingAPIKey": {
			args:            args{creds: map[string]string{}},
			wantErrContains: `credentials secret has no "api_key" key`,
		},
		"ExtraKeysIgnored": {
			args: args{creds: map[string]string{"api_key": "vultr-key", "unrelated": "x"}},
			want: map[string]any{"api_key": "vultr-key"},
		},
		"RateAndRetryLimitsSet": {
			args: args{
				creds:      map[string]string{"api_key": "vultr-key"},
				rateLimit:  new(int64(600)),
				retryLimit: new(int64(3)),
			},
			want: map[string]any{"api_key": "vultr-key", "rate_limit": 600, "retry_limit": 3},
		},
		"RetryLimitZeroIsPassed": {
			args: args{
				creds:      map[string]string{"api_key": "vultr-key"},
				retryLimit: new(int64(0)),
			},
			want: map[string]any{"api_key": "vultr-key", "retry_limit": 0},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := buildConfiguration(tc.args.creds, tc.args.rateLimit, tc.args.retryLimit)
			if tc.wantErrContains != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErrContains)
				}
				if !strings.Contains(err.Error(), tc.wantErrContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tc.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("buildConfiguration() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
