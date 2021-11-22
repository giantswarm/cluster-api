/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/cluster-api/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

func TestFuzzyConversion(t *testing.T) {
	t.Run("for Cluster", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.Cluster{},
		Spoke:       &Cluster{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{ClusterJSONFuzzFuncs},
	}))
	t.Run("for Machine", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.Machine{},
		Spoke:       &Machine{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{BootstrapFuzzFuncs, MachineStatusFuzzFunc, CustomObjectMetaFuzzFunc},
	}))
	t.Run("for MachineSet", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.MachineSet{},
		Spoke:       &MachineSet{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{BootstrapFuzzFuncs, CustomObjectMetaFuzzFunc},
	}))
	t.Run("for MachineDeployment", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:         &v1beta1.MachineDeployment{},
		Spoke:       &MachineDeployment{},
		FuzzerFuncs: []fuzzer.FuzzerFuncs{BootstrapFuzzFuncs, CustomObjectMetaFuzzFunc},
	}))
}

func TestConvertCluster(t *testing.T) {
	t.Run("to hub", func(t *testing.T) {
		t.Run("should convert the first value in Status.APIEndpoints to Spec.ControlPlaneEndpoint", func(t *testing.T) {
			g := NewWithT(t)

			src := &Cluster{
				Status: ClusterStatus{
					APIEndpoints: []APIEndpoint{
						{
							Host: "example.com",
							Port: 6443,
						},
					},
				},
			}
			dst := &v1beta1.Cluster{}

			g.Expect(src.ConvertTo(dst)).To(Succeed())
			g.Expect(dst.Spec.ControlPlaneEndpoint.Host).To(Equal("example.com"))
			g.Expect(dst.Spec.ControlPlaneEndpoint.Port).To(BeEquivalentTo(6443))
		})
	})

	t.Run("from hub", func(t *testing.T) {
		t.Run("preserves fields from hub version", func(t *testing.T) {
			g := NewWithT(t)

			src := &v1beta1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "hub",
				},
				Spec: v1beta1.ClusterSpec{
					ControlPlaneRef: &corev1.ObjectReference{
						Name: "controlplane-1",
					},
				},
				Status: v1beta1.ClusterStatus{
					ControlPlaneReady: true,
				},
			}
			dst := &Cluster{}

			g.Expect(dst.ConvertFrom(src)).To(Succeed())
			restored := &v1beta1.Cluster{}
			g.Expect(dst.ConvertTo(restored)).To(Succeed())

			// Test field restored fields.
			g.Expect(restored.Name).To(Equal(src.Name))
			g.Expect(restored.Spec.ControlPlaneRef).To(Equal(src.Spec.ControlPlaneRef))
			g.Expect(restored.Status.ControlPlaneReady).To(Equal(src.Status.ControlPlaneReady))
		})

		t.Run("should convert Spec.ControlPlaneEndpoint to Status.APIEndpoints[0]", func(t *testing.T) {
			g := NewWithT(t)

			src := &v1beta1.Cluster{
				Spec: v1beta1.ClusterSpec{
					ControlPlaneEndpoint: v1beta1.APIEndpoint{
						Host: "example.com",
						Port: 6443,
					},
				},
			}
			dst := &Cluster{}

			g.Expect(dst.ConvertFrom(src)).To(Succeed())
			g.Expect(dst.Status.APIEndpoints[0].Host).To(Equal("example.com"))
			g.Expect(dst.Status.APIEndpoints[0].Port).To(BeEquivalentTo(6443))
		})
	})
}

func TestConvertMachine(t *testing.T) {
	t.Run("to hub", func(t *testing.T) {
		t.Run("should convert the Spec.ClusterName from label", func(t *testing.T) {
			g := NewWithT(t)

			src := &Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						MachineClusterLabelName: "test-cluster",
					},
				},
			}
			dst := &v1beta1.Machine{}

			g.Expect(src.ConvertTo(dst)).To(Succeed())
			g.Expect(dst.Spec.ClusterName).To(Equal("test-cluster"))
		})
	})

	t.Run("from hub", func(t *testing.T) {
		t.Run("preserves fields from hub version", func(t *testing.T) {
			g := NewWithT(t)

			failureDomain := "my failure domain"
			src := &v1beta1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "hub",
				},
				Spec: v1beta1.MachineSpec{
					ClusterName: "test-cluster",
					Bootstrap: v1beta1.Bootstrap{
						DataSecretName: pointer.StringPtr("secret-data"),
					},
					FailureDomain: &failureDomain,
				},
			}
			dst := &Machine{}

			g.Expect(dst.ConvertFrom(src)).To(Succeed())
			restored := &v1beta1.Machine{}
			g.Expect(dst.ConvertTo(restored)).To(Succeed())

			// Test field restored fields.
			g.Expect(restored.Name).To(Equal(src.Name))
			g.Expect(restored.Spec.Bootstrap.DataSecretName).To(Equal(src.Spec.Bootstrap.DataSecretName))
			g.Expect(restored.Spec.ClusterName).To(Equal(src.Spec.ClusterName))
			g.Expect(restored.Spec.FailureDomain).To(Equal(src.Spec.FailureDomain))
		})
	})
}

func TestConvertMachineSet(t *testing.T) {
	t.Run("to hub", func(t *testing.T) {
		t.Run("should convert the Spec.ClusterName from label", func(t *testing.T) {
			g := NewWithT(t)

			src := &MachineSet{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						MachineClusterLabelName: "test-cluster",
					},
				},
			}
			dst := &v1beta1.MachineSet{}

			g.Expect(src.ConvertTo(dst)).To(Succeed())
			g.Expect(dst.Spec.ClusterName).To(Equal("test-cluster"))
			g.Expect(dst.Spec.Template.Spec.ClusterName).To(Equal("test-cluster"))
		})
	})

	t.Run("from hub", func(t *testing.T) {
		t.Run("preserves field from hub version", func(t *testing.T) {
			g := NewWithT(t)

			src := &v1beta1.MachineSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "hub",
				},
				Spec: v1beta1.MachineSetSpec{
					ClusterName: "test-cluster",
					Template: v1beta1.MachineTemplateSpec{
						Spec: v1beta1.MachineSpec{
							ClusterName: "test-cluster",
						},
					},
				},
			}
			dst := &MachineSet{}

			g.Expect(dst.ConvertFrom(src)).To(Succeed())
			restored := &v1beta1.MachineSet{}
			g.Expect(dst.ConvertTo(restored)).To(Succeed())

			// Test field restored fields.
			g.Expect(restored.Name).To(Equal(src.Name))
			g.Expect(restored.Spec.ClusterName).To(Equal(src.Spec.ClusterName))
			g.Expect(restored.Spec.Template.Spec.ClusterName).To(Equal(src.Spec.Template.Spec.ClusterName))
		})
	})
}

func TestConvertMachineDeployment(t *testing.T) {
	t.Run("to hub", func(t *testing.T) {
		t.Run("should convert the Spec.ClusterName from label", func(t *testing.T) {
			g := NewWithT(t)

			src := &MachineDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						MachineClusterLabelName: "test-cluster",
					},
				},
			}
			dst := &v1beta1.MachineDeployment{}

			g.Expect(src.ConvertTo(dst)).To(Succeed())
			g.Expect(dst.Spec.ClusterName).To(Equal("test-cluster"))
			g.Expect(dst.Spec.Template.Spec.ClusterName).To(Equal("test-cluster"))
		})

		t.Run("should convert the annotations", func(t *testing.T) {
			g := NewWithT(t)

			src := &MachineDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						RevisionAnnotation:        "test",
						RevisionHistoryAnnotation: "test",
						DesiredReplicasAnnotation: "test",
						MaxReplicasAnnotation:     "test",
					},
				},
			}
			dst := &v1beta1.MachineDeployment{}

			g.Expect(src.ConvertTo(dst)).To(Succeed())
			g.Expect(dst.Annotations).To(HaveKey(v1beta1.RevisionAnnotation))
			g.Expect(dst.Annotations).To(HaveKey(v1beta1.RevisionHistoryAnnotation))
			g.Expect(dst.Annotations).To(HaveKey(v1beta1.DesiredReplicasAnnotation))
			g.Expect(dst.Annotations).To(HaveKey(v1beta1.MaxReplicasAnnotation))
		})
	})

	t.Run("from hub", func(t *testing.T) {
		t.Run("preserves fields from hub version", func(t *testing.T) {
			g := NewWithT(t)

			src := &v1beta1.MachineDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "hub",
				},
				Spec: v1beta1.MachineDeploymentSpec{
					ClusterName: "test-cluster",
					Paused:      true,
					Template: v1beta1.MachineTemplateSpec{
						Spec: v1beta1.MachineSpec{
							ClusterName: "test-cluster",
						},
					},
				},
			}
			src.Status.SetTypedPhase(v1beta1.MachineDeploymentPhaseRunning)
			dst := &MachineDeployment{}

			g.Expect(dst.ConvertFrom(src)).To(Succeed())
			restored := &v1beta1.MachineDeployment{}
			g.Expect(dst.ConvertTo(restored)).To(Succeed())

			// Test field restored fields.
			g.Expect(restored.Name).To(Equal(src.Name))
			g.Expect(restored.Spec.ClusterName).To(Equal(src.Spec.ClusterName))
			g.Expect(restored.Spec.Paused).To(Equal(src.Spec.Paused))
			g.Expect(restored.Spec.Template.Spec.ClusterName).To(Equal(src.Spec.Template.Spec.ClusterName))
		})

		t.Run("should convert the annotations", func(t *testing.T) {
			g := NewWithT(t)

			src := &v1beta1.MachineDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						v1beta1.RevisionAnnotation:        "test",
						v1beta1.RevisionHistoryAnnotation: "test",
						v1beta1.DesiredReplicasAnnotation: "test",
						v1beta1.MaxReplicasAnnotation:     "test",
					},
				},
			}
			dst := &MachineDeployment{}

			g.Expect(dst.ConvertFrom(src)).To(Succeed())
			g.Expect(dst.Annotations).To(HaveKey(RevisionAnnotation))
			g.Expect(dst.Annotations).To(HaveKey(RevisionHistoryAnnotation))
			g.Expect(dst.Annotations).To(HaveKey(DesiredReplicasAnnotation))
			g.Expect(dst.Annotations).To(HaveKey(MaxReplicasAnnotation))
		})
	})
}

func MachineStatusFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		MachineStatusFuzzer,
	}
}

func MachineStatusFuzzer(in *MachineStatus, c fuzz.Continue) {
	c.FuzzNoCustom(in)

	// These fields have been removed in v1beta1
	// data is going to be lost, so we're forcing zero values to avoid round trip errors.
	in.Version = nil
}

func BootstrapFuzzFuncs(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		BootstrapFuzzer,
		MachineSpecFuzzer,
	}
}

func BootstrapFuzzer(obj *Bootstrap, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// Bootstrap.Data has been removed in v1alpha4, so setting it to nil in order to avoid v1alpha3 --> <hub> --> v1alpha3 round trip errors.
	obj.Data = nil
}

func ClusterJSONFuzzFuncs(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		ClusterVariableFuzzer,
	}
}

func ClusterVariableFuzzer(in *v1beta1.ClusterVariable, c fuzz.Continue) {
	c.FuzzNoCustom(in)

	// Not every random byte array is valid JSON, e.g. a string without `""`,so we're setting a valid value.
	in.Value = apiextensionsv1.JSON{Raw: []byte("\"test-string\"")}
}

func MachineSpecFuzzer(in *v1beta1.MachineSpec, c fuzz.Continue) {
	c.FuzzNoCustom(in)

	// ClusterName is stored as a label in v1alpha2 but in v1beta1 it is available on both the top-level spec and within the template.
	// We need to manually set it here to ensure the restored value is checked against, otherwise an empty string it expected.
	in.ClusterName = "cluster-name"
}

func CustomObjectMetaFuzzFunc(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		CustomObjectMetaFuzzer,
		ObjectMetaFuzzer,
	}
}

func CustomObjectMetaFuzzer(in *ObjectMeta, c fuzz.Continue) {
	c.FuzzNoCustom(in)

	// These fields have been removed in v1alpha4
	// data is going to be lost, so we're forcing zero values here.
	in.Name = ""
	in.GenerateName = ""
	in.Namespace = ""
	in.OwnerReferences = nil
	in.Annotations = make(map[string]string)
	in.Labels = make(map[string]string)
}

func ObjectMetaFuzzer(in *metav1.ObjectMeta, c fuzz.Continue) {
	c.FuzzNoCustom(in)

	in.Annotations = make(map[string]string)
	in.Labels = make(map[string]string)
}
