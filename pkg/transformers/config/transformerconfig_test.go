/*
Copyright 2018 The Kubernetes Authors.

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

package config

import (
	"testing"

	"reflect"

	"sigs.k8s.io/kustomize/pkg/gvk"
)

func TestAddNamereferenceFieldSpec(t *testing.T) {
	cfg := &TransformerConfig{}

	nbrs := NameBackReferences{
		Gvk: gvk.Gvk{
			Kind: "KindA",
		},
		FieldSpecs: []FieldSpec{
			{
				Gvk: gvk.Gvk{
					Kind: "KindB",
				},
				Path:               "path/to/a/field",
				CreateIfNotPresent: false,
			},
		},
	}

	cfg.AddNamereferenceFieldSpec(nbrs)
	if len(cfg.NameReference) != 1 {
		t.Fatal("failed to add namerefence FieldSpec")
	}
}

func TestAddFieldSpecs(t *testing.T) {
	cfg := &TransformerConfig{}

	fieldSpec := FieldSpec{
		Gvk:                gvk.Gvk{Group: "GroupA", Kind: "KindB"},
		Path:               "path/to/a/field",
		CreateIfNotPresent: true,
	}

	cfg.AddPrefixFieldSpec(fieldSpec)
	if len(cfg.NamePrefix) != 1 {
		t.Fatalf("failed to add nameprefix FieldSpec")
	}
	cfg.AddLabelFieldSpec(fieldSpec)
	if len(cfg.CommonLabels) != 1 {
		t.Fatalf("failed to add nameprefix FieldSpec")
	}
	cfg.AddAnnotationFieldSpec(fieldSpec)
	if len(cfg.CommonAnnotations) != 1 {
		t.Fatalf("failed to add nameprefix FieldSpec")
	}
}

func TestMerge(t *testing.T) {
	nameReference := []NameBackReferences{
		{
			Gvk: gvk.Gvk{
				Kind: "KindA",
			},
			FieldSpecs: []FieldSpec{
				{
					Gvk: gvk.Gvk{
						Kind: "KindB",
					},
					Path:               "path/to/a/field",
					CreateIfNotPresent: false,
				},
			},
		},
		{
			Gvk: gvk.Gvk{
				Kind: "KindA",
			},
			FieldSpecs: []FieldSpec{
				{
					Gvk: gvk.Gvk{
						Kind: "KindC",
					},
					Path:               "path/to/a/field",
					CreateIfNotPresent: false,
				},
			},
		},
	}
	fieldSpecs := []FieldSpec{
		{
			Gvk:                gvk.Gvk{Group: "GroupA", Kind: "KindB"},
			Path:               "path/to/a/field",
			CreateIfNotPresent: true,
		},
		{
			Gvk:                gvk.Gvk{Group: "GroupA", Kind: "KindC"},
			Path:               "path/to/a/field",
			CreateIfNotPresent: true,
		},
	}
	cfga := &TransformerConfig{}
	cfga.AddNamereferenceFieldSpec(nameReference[0])
	cfga.AddPrefixFieldSpec(fieldSpecs[0])

	cfgb := &TransformerConfig{}
	cfgb.AddNamereferenceFieldSpec(nameReference[1])
	cfgb.AddPrefixFieldSpec(fieldSpecs[1])

	actual := cfga.Merge(cfgb)

	if len(actual.NamePrefix) != 2 {
		t.Fatal("merge failed for namePrefix FieldSpec")
	}

	if len(actual.NameReference) != 1 {
		t.Fatal("merge failed for namereference FieldSpec")
	}

	expected := &TransformerConfig{}
	expected.AddNamereferenceFieldSpec(nameReference[0])
	expected.AddNamereferenceFieldSpec(nameReference[1])
	expected.AddPrefixFieldSpec(fieldSpecs[0])
	expected.AddPrefixFieldSpec(fieldSpecs[1])

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %v\n but got: %v\n", expected, actual)
	}
}
