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
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this KubeadmConfig to the Hub version (v1beta1).
func (src *KubeadmConfig) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.KubeadmConfig)
	if err := Convert_v1alpha2_KubeadmConfig_To_v1beta1_KubeadmConfig(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &v1beta1.KubeadmConfig{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	if restored.Spec.JoinConfiguration != nil && restored.Spec.JoinConfiguration.NodeRegistration.IgnorePreflightErrors != nil {
		if dst.Spec.JoinConfiguration == nil {
			dst.Spec.JoinConfiguration = &v1beta1.JoinConfiguration{}
		}
		dst.Spec.JoinConfiguration.NodeRegistration.IgnorePreflightErrors = restored.Spec.JoinConfiguration.NodeRegistration.IgnorePreflightErrors
	}

	if restored.Spec.InitConfiguration != nil && restored.Spec.InitConfiguration.NodeRegistration.IgnorePreflightErrors != nil {
		if dst.Spec.InitConfiguration == nil {
			dst.Spec.InitConfiguration = &v1beta1.InitConfiguration{}
		}
		dst.Spec.InitConfiguration.NodeRegistration.IgnorePreflightErrors = restored.Spec.InitConfiguration.NodeRegistration.IgnorePreflightErrors
	}

	dst.Status.DataSecretName = restored.Status.DataSecretName
	dst.Status.ObservedGeneration = restored.Status.ObservedGeneration
	dst.Spec.Verbosity = restored.Spec.Verbosity
	dst.Spec.UseExperimentalRetryJoin = restored.Spec.UseExperimentalRetryJoin
	dst.Spec.DiskSetup = restored.Spec.DiskSetup
	dst.Spec.Mounts = restored.Spec.Mounts
	dst.Spec.Files = restored.Spec.Files
	dst.Status.Conditions = restored.Status.Conditions

	// Track files successfully up-converted. We need this to dedupe
	// restored files from user-updated files on up-conversion. We store
	// them as pointers for later modification without paying for second
	// lookup.
	dstPaths := make(map[string]*v1beta1.File, len(dst.Spec.Files))
	for i := range dst.Spec.Files {
		path := dst.Spec.Files[i].Path
		dstPaths[path] = &dst.Spec.Files[i]
	}

	// If we find a restored file matching the file path of a v1alpha2
	// file with no content, we should restore contentFrom to that file.
	for i := range restored.Spec.Files {
		restoredFile := restored.Spec.Files[i]
		dstFile, exists := dstPaths[restoredFile.Path]
		if exists && dstFile.Content == "" {
			if dstFile.ContentFrom == nil {
				dstFile.ContentFrom = new(v1beta1.FileSource)
			}
			*dstFile.ContentFrom = *restoredFile.ContentFrom
		}
	}

	return nil
}

// ConvertFrom converts from the KubeadmConfig Hub version (v1beta1) to this version.
func (dst *KubeadmConfig) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.KubeadmConfig)
	if err := Convert_v1beta1_KubeadmConfig_To_v1alpha2_KubeadmConfig(src, dst, nil); err != nil {
		return nil
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this KubeadmConfigList to the Hub version (v1beta1).
func (src *KubeadmConfigList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.KubeadmConfigList)
	return Convert_v1alpha2_KubeadmConfigList_To_v1beta1_KubeadmConfigList(src, dst, nil)
}

// ConvertFrom converts from the KubeadmConfigList Hub version (v1beta1) to this version.
func (dst *KubeadmConfigList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.KubeadmConfigList)
	return Convert_v1beta1_KubeadmConfigList_To_v1alpha2_KubeadmConfigList(src, dst, nil)
}

// ConvertTo converts this KubeadmConfigTemplate to the Hub version (v1beta1).
func (src *KubeadmConfigTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.KubeadmConfigTemplate)
	if err := Convert_v1alpha2_KubeadmConfigTemplate_To_v1beta1_KubeadmConfigTemplate(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &v1beta1.KubeadmConfigTemplate{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	if restored.Spec.Template.Spec.JoinConfiguration != nil && restored.Spec.Template.Spec.JoinConfiguration.NodeRegistration.IgnorePreflightErrors != nil {
		if dst.Spec.Template.Spec.JoinConfiguration == nil {
			dst.Spec.Template.Spec.JoinConfiguration = &v1beta1.JoinConfiguration{}
		}
		dst.Spec.Template.Spec.JoinConfiguration.NodeRegistration.IgnorePreflightErrors = restored.Spec.Template.Spec.JoinConfiguration.NodeRegistration.IgnorePreflightErrors
	}

	if restored.Spec.Template.Spec.InitConfiguration != nil && restored.Spec.Template.Spec.InitConfiguration.NodeRegistration.IgnorePreflightErrors != nil {
		if dst.Spec.Template.Spec.InitConfiguration == nil {
			dst.Spec.Template.Spec.InitConfiguration = &v1beta1.InitConfiguration{}
		}
		dst.Spec.Template.Spec.InitConfiguration.NodeRegistration.IgnorePreflightErrors = restored.Spec.Template.Spec.InitConfiguration.NodeRegistration.IgnorePreflightErrors
	}

	return nil
}

// ConvertFrom converts from the KubeadmConfigTemplate Hub version (v1beta1) to this version.
func (dst *KubeadmConfigTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.KubeadmConfigTemplate)
	if err := Convert_v1beta1_KubeadmConfigTemplate_To_v1alpha2_KubeadmConfigTemplate(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this KubeadmConfigTemplateList to the Hub version (v1beta1).
func (src *KubeadmConfigTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.KubeadmConfigTemplateList)
	return Convert_v1alpha2_KubeadmConfigTemplateList_To_v1beta1_KubeadmConfigTemplateList(src, dst, nil)
}

// ConvertFrom converts from the KubeadmConfigTemplateList Hub version (v1beta1) to this version.
func (dst *KubeadmConfigTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.KubeadmConfigTemplateList)
	return Convert_v1beta1_KubeadmConfigTemplateList_To_v1alpha2_KubeadmConfigTemplateList(src, dst, nil)
}

// Convert_v1alpha2_KubeadmConfigStatus_To_v1beta1_KubeadmConfigStatus converts this KubeadmConfigStatus to the Hub version (v1beta1).
func Convert_v1alpha2_KubeadmConfigStatus_To_v1beta1_KubeadmConfigStatus(in *KubeadmConfigStatus, out *v1beta1.KubeadmConfigStatus, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_KubeadmConfigStatus_To_v1beta1_KubeadmConfigStatus(in, out, s); err != nil {
		return err
	}

	// Manually convert the Error fields to the Failure fields
	out.FailureMessage = in.ErrorMessage
	out.FailureReason = in.ErrorReason

	// Note: Unable to convert between inline BootstrapData and data secret

	return nil
}

// Convert_v1beta1_KubeadmConfigStatus_To_v1alpha2_KubeadmConfigStatus converts from the Hub version (v1beta1) of the KubeadmConfigStatus to this version.
func Convert_v1beta1_KubeadmConfigStatus_To_v1alpha2_KubeadmConfigStatus(in *v1beta1.KubeadmConfigStatus, out *KubeadmConfigStatus, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_KubeadmConfigStatus_To_v1alpha2_KubeadmConfigStatus(in, out, s); err != nil {
		return err
	}

	// Manually convert the Failure fields to the Error fields
	out.ErrorMessage = in.FailureMessage
	out.ErrorReason = in.FailureReason

	// Note: Unable to convert between inline BootstrapData and data secret

	return nil
}

// Convert_v1alpha2_KubeadmConfigSpec_To_v1beta1_KubeadmConfigSpec converts this KubeadmConfigSpec to the Hub version (v1beta1).
func Convert_v1alpha2_KubeadmConfigSpec_To_v1beta1_KubeadmConfigSpec(in *KubeadmConfigSpec, out *v1beta1.KubeadmConfigSpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha2_KubeadmConfigSpec_To_v1beta1_KubeadmConfigSpec(in, out, s)
}

// Convert_v1beta1_KubeadmConfigSpec_To_v1alpha2_KubeadmConfigSpec converts from the Hub version (v1beta1) of the KubeadmConfigSpec to this version.
func Convert_v1beta1_KubeadmConfigSpec_To_v1alpha2_KubeadmConfigSpec(in *v1beta1.KubeadmConfigSpec, out *KubeadmConfigSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta1_KubeadmConfigSpec_To_v1alpha2_KubeadmConfigSpec(in, out, s)
}

// Convert_v1beta1_File_To_v1alpha2_File converts from the Hub version (v1beta1) of the File to this version.
func Convert_v1beta1_File_To_v1alpha2_File(in *v1beta1.File, out *File, s apiconversion.Scope) error {
	return autoConvert_v1beta1_File_To_v1alpha2_File(in, out, s)
}
