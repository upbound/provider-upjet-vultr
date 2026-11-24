package clusterauth

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/vultr/govultr/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/upbound/provider-vultr/apis/namespaced/vke/v1beta1"
)

const (
	clusterID = "0f1e2d3c-example"
	endpoint  = "https://example.vultr-k8s.com:6443"

	// The certificate fields hold base64 of "ca", "cert" and "key".
	kubeconfig = `apiVersion: v1
kind: Config
current-context: vke-example
clusters:
- name: vke-example
  cluster:
    server: https://example.vultr-k8s.com:6443
    certificate-authority-data: Y2E=
contexts:
- name: vke-example
  context:
    cluster: vke-example
    user: admin
users:
- name: admin
  user:
    client-certificate-data: Y2VydA==
    client-key-data: a2V5
`
	kubeconfigNoContext = `apiVersion: v1
kind: Config
current-context: missing
`
)

type fakeGetter struct {
	kubeconfig string
	err        error
	gotID      string
}

func (f *fakeGetter) GetKubeConfig(_ context.Context, vkeID string) (*govultr.KubeConfig, *http.Response, error) {
	f.gotID = vkeID
	if f.err != nil {
		return nil, nil, f.err
	}
	return &govultr.KubeConfig{KubeConfig: f.kubeconfig}, nil, nil
}

func encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func clusterAuth(id *string, deleted bool) *v1beta1.ClusterAuth {
	cr := &v1beta1.ClusterAuth{
		Spec: v1beta1.ClusterAuthSpec{ForProvider: v1beta1.ClusterAuthParameters{ClusterID: id}},
	}
	if deleted {
		cr.DeletionTimestamp = ptr.To(metav1.Now())
	}
	return cr
}

func TestObserve(t *testing.T) {
	errBoom := errors.New("boom")

	type args struct {
		getter *fakeGetter
		cr     *v1beta1.ClusterAuth
	}
	type want struct {
		observation managed.ExternalObservation
		endpoint    *string
		requestedID string
		err         error
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"PublishesKubeconfig": {
			args: args{
				getter: &fakeGetter{kubeconfig: encode(kubeconfig)},
				cr:     clusterAuth(ptr.To(clusterID), false),
			},
			want: want{
				observation: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: true,
					ConnectionDetails: managed.ConnectionDetails{
						keyKubeconfig: []byte(kubeconfig),
						keyEndpoint:   []byte(endpoint),
						keyClusterCA:  []byte("ca"),
						keyClientCert: []byte("cert"),
						keyClientKey:  []byte("key"),
					},
				},
				endpoint:    ptr.To(endpoint),
				requestedID: clusterID,
			},
		},
		"DeletedResourceDoesNotExist": {
			args: args{
				getter: &fakeGetter{err: errBoom},
				cr:     clusterAuth(ptr.To(clusterID), true),
			},
			want: want{
				observation: managed.ExternalObservation{ResourceExists: false},
			},
		},
		"MissingClusterID": {
			args: args{
				getter: &fakeGetter{kubeconfig: encode(kubeconfig)},
				cr:     clusterAuth(nil, false),
			},
			want: want{
				err: errors.New(errNoClusterID),
			},
		},
		"GetKubeconfigFails": {
			args: args{
				getter: &fakeGetter{err: errBoom},
				cr:     clusterAuth(ptr.To(clusterID), false),
			},
			want: want{
				requestedID: clusterID,
				err:         errors.Wrap(errBoom, errGetKubeconfig),
			},
		},
		"KubeconfigIsNotBase64": {
			args: args{
				getter: &fakeGetter{kubeconfig: "not base64!"},
				cr:     clusterAuth(ptr.To(clusterID), false),
			},
			want: want{
				requestedID: clusterID,
				err:         errors.Wrap(base64.CorruptInputError(3), errDecodeKubeconfig),
			},
		},
		"KubeconfigWithoutCurrentContext": {
			args: args{
				getter: &fakeGetter{kubeconfig: encode(kubeconfigNoContext)},
				cr:     clusterAuth(ptr.To(clusterID), false),
			},
			want: want{
				requestedID: clusterID,
				err:         errors.New(errNoCurrentContext),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: tc.args.getter}
			got, err := e.Observe(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("Observe(...): -want error, +got error:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.observation, got); diff != "" {
				t.Errorf("Observe(...): -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.endpoint, tc.args.cr.Status.AtProvider.Endpoint); diff != "" {
				t.Errorf("Observe(...): -want endpoint, +got endpoint:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.requestedID, tc.args.getter.gotID); diff != "" {
				t.Errorf("Observe(...): -want requested cluster ID, +got requested cluster ID:\n%s", diff)
			}
		})
	}
}
