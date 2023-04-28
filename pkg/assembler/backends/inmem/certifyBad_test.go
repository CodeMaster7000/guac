//
// Copyright 2023 The GUAC Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inmem_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/guacsec/guac/internal/testing/ptrfrom"
	"github.com/guacsec/guac/pkg/assembler/backends/inmem"
	"github.com/guacsec/guac/pkg/assembler/graphql/model"
	"golang.org/x/exp/slices"
)

func TestCertifyBad(t *testing.T) {
	type call struct {
		Sub   model.PackageSourceOrArtifactInput
		Match *model.MatchFlags
		CB    *model.CertifyBadInputSpec
	}
	tests := []struct {
		Name         string
		InPkg        []*model.PkgInputSpec
		InSrc        []*model.SourceInputSpec
		InArt        []*model.ArtifactInputSpec
		Calls        []call
		Query        *model.CertifyBadSpec
		ExpCB        []*model.CertifyBad
		ExpIngestErr bool
		ExpQueryErr  bool
	}{
		{
			Name:  "HappyPath",
			InPkg: []*model.PkgInputSpec{p1},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p1,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeSpecificVersion,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				Justification: ptrfrom.String("test justification"),
			},
			ExpCB: []*model.CertifyBad{
				&model.CertifyBad{
					Subject:       p1out,
					Justification: "test justification",
				},
			},
		},
		{
			Name:  "HappyPath All Version",
			InPkg: []*model.PkgInputSpec{p1},
			Calls: []call{
				call{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p1,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeAllVersions,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				Justification: ptrfrom.String("test justification"),
			},
			ExpCB: []*model.CertifyBad{
				&model.CertifyBad{
					Subject:       p1outName,
					Justification: "test justification",
				},
			},
		},
		{
			Name:  "Ingest same twice",
			InPkg: []*model.PkgInputSpec{p1},
			Calls: []call{
				call{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p1,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeSpecificVersion,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				call{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p1,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeSpecificVersion,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				Justification: ptrfrom.String("test justification"),
			},
			ExpCB: []*model.CertifyBad{
				&model.CertifyBad{
					Subject:       p1out,
					Justification: "test justification",
				},
			},
		},
		{
			Name:  "Query on Justification",
			InPkg: []*model.PkgInputSpec{p1},
			Calls: []call{
				call{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p1,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeSpecificVersion,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification one",
					},
				},
				call{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p1,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeSpecificVersion,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification two",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				Justification: ptrfrom.String("test justification one"),
			},
			ExpCB: []*model.CertifyBad{
				&model.CertifyBad{
					Subject:       p1out,
					Justification: "test justification one",
				},
			},
		},
		{
			Name:  "Query on Package",
			InPkg: []*model.PkgInputSpec{p1, p2},
			InSrc: []*model.SourceInputSpec{s1},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p1,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeSpecificVersion,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p2,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeSpecificVersion,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Source: s1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				Subject: &model.PackageSourceOrArtifactSpec{
					Package: &model.PkgSpec{
						Version: ptrfrom.String("2.11.1"),
					},
				},
			},
			ExpCB: []*model.CertifyBad{
				{
					Subject:       p2out,
					Justification: "test justification",
				},
			},
		},
		{
			Name:  "Query on Source",
			InPkg: []*model.PkgInputSpec{p1},
			InSrc: []*model.SourceInputSpec{s1, s2},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p1,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeSpecificVersion,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Source: s1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Source: s2,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				Subject: &model.PackageSourceOrArtifactSpec{
					Source: &model.SourceSpec{
						Name: ptrfrom.String("bobsrepo"),
					},
				},
			},
			ExpCB: []*model.CertifyBad{
				{
					Subject:       s2out,
					Justification: "test justification",
				},
			},
		},
		{
			Name:  "Query on Artifact",
			InSrc: []*model.SourceInputSpec{s1},
			InArt: []*model.ArtifactInputSpec{a1, a2},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Artifact: a1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Artifact: a2,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Source: s1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				Subject: &model.PackageSourceOrArtifactSpec{
					Artifact: &model.ArtifactSpec{
						Algorithm: ptrfrom.String("sha1"),
					},
				},
			},
			ExpCB: []*model.CertifyBad{
				{
					Subject:       a2out,
					Justification: "test justification",
				},
			},
		},
		{
			Name:  "Query none",
			InArt: []*model.ArtifactInputSpec{a1, a2},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Artifact: a1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Artifact: a2,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				Subject: &model.PackageSourceOrArtifactSpec{
					Artifact: &model.ArtifactSpec{
						Algorithm: ptrfrom.String("asdf"),
					},
				},
			},
			ExpCB: nil,
		},
		{
			Name:  "Query multiple",
			InSrc: []*model.SourceInputSpec{s1, s2},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Source: s1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Source: s2,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				Justification: ptrfrom.String("test justification"),
			},
			ExpCB: []*model.CertifyBad{
				{
					Subject:       s1out,
					Justification: "test justification",
				},
				{
					Subject:       s2out,
					Justification: "test justification",
				},
			},
		},
		{
			Name:  "Query Packages",
			InPkg: []*model.PkgInputSpec{p1, p2},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p1,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeSpecificVersion,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p2,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeSpecificVersion,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p2,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeAllVersions,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				Subject: &model.PackageSourceOrArtifactSpec{
					Package: &model.PkgSpec{
						Version: ptrfrom.String("2.11.1"),
					},
				},
			},
			ExpCB: []*model.CertifyBad{
				{
					Subject:       p2out,
					Justification: "test justification",
				},
				{
					Subject:       p1outName,
					Justification: "test justification",
				},
			},
		},
		{
			Name:  "Query ID",
			InArt: []*model.ArtifactInputSpec{a1, a2},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Artifact: a1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Artifact: a2,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				ID: ptrfrom.String("4"),
			},
			ExpCB: []*model.CertifyBad{
				{
					Subject:       a1out,
					Justification: "test justification",
				},
			},
		},
		{
			Name: "Ingest without subject",
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Artifact: a1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			ExpIngestErr: true,
		},
		{
			Name:  "Ingest with two subjects",
			InSrc: []*model.SourceInputSpec{s1},
			InArt: []*model.ArtifactInputSpec{a1},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Source:   s1,
						Artifact: a1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			ExpIngestErr: true,
		},
		{
			Name:  "Query with two subjects",
			InSrc: []*model.SourceInputSpec{s1},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Source: s1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				Subject: &model.PackageSourceOrArtifactSpec{
					Package: &model.PkgSpec{
						Version: ptrfrom.String("2.11.1"),
					},
					Artifact: &model.ArtifactSpec{
						Algorithm: ptrfrom.String("asdf"),
					},
				},
			},
			ExpQueryErr: true,
		},
		{
			Name:  "Query bad ID",
			InSrc: []*model.SourceInputSpec{s1},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Source: s1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			Query: &model.CertifyBadSpec{
				ID: ptrfrom.String("asdf"),
			},
			ExpQueryErr: true,
		},
	}
	ignoreID := cmp.FilterPath(func(p cmp.Path) bool {
		return strings.Compare(".ID", p[len(p)-1].String()) == 0
	}, cmp.Ignore())
	ctx := context.Background()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			b, err := inmem.GetBackend(nil)
			if err != nil {
				t.Fatalf("Could not instantiate testing backend: %v", err)
			}
			for _, p := range test.InPkg {
				if _, err := b.IngestPackage(ctx, *p); err != nil {
					t.Fatalf("Could not ingest package: %v", err)
				}
			}
			for _, s := range test.InSrc {
				if _, err := b.IngestSource(ctx, *s); err != nil {
					t.Fatalf("Could not ingest source: %v", err)
				}
			}
			for _, a := range test.InArt {
				if _, err := b.IngestArtifact(ctx, a); err != nil {
					t.Fatalf("Could not ingest artifact: %v", err)
				}
			}
			for _, o := range test.Calls {
				_, err := b.IngestCertifyBad(ctx, o.Sub, o.Match, *o.CB)
				if (err != nil) != test.ExpIngestErr {
					t.Fatalf("did not get expected ingest error, want: %v, got: %v", test.ExpIngestErr, err)
				}
				if err != nil {
					return
				}
			}
			got, err := b.CertifyBad(ctx, test.Query)
			if (err != nil) != test.ExpQueryErr {
				t.Fatalf("did not get expected query error, want: %v, got: %v", test.ExpQueryErr, err)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(test.ExpCB, got, ignoreID); diff != "" {
				t.Errorf("Unexpected results. (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCertifyBadNeighbors(t *testing.T) {
	type call struct {
		Sub   model.PackageSourceOrArtifactInput
		Match *model.MatchFlags
		CB    *model.CertifyBadInputSpec
	}
	tests := []struct {
		Name         string
		InPkg        []*model.PkgInputSpec
		InSrc        []*model.SourceInputSpec
		InArt        []*model.ArtifactInputSpec
		Calls        []call
		ExpNeighbors map[string][]string
	}{
		{
			Name:  "HappyPath",
			InPkg: []*model.PkgInputSpec{p1},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p1,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeSpecificVersion,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			ExpNeighbors: map[string][]string{
				"5": []string{"2", "6"}, // pkg version
				"6": []string{"2"},      // certify bad
			},
		},
		{
			Name:  "Pkg Name Src and Artifact",
			InPkg: []*model.PkgInputSpec{p1},
			InSrc: []*model.SourceInputSpec{s1},
			InArt: []*model.ArtifactInputSpec{a1},
			Calls: []call{
				{
					Sub: model.PackageSourceOrArtifactInput{
						Package: p1,
					},
					Match: &model.MatchFlags{
						Pkg: model.PkgMatchTypeAllVersions,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Source: s1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
				{
					Sub: model.PackageSourceOrArtifactInput{
						Artifact: a1,
					},
					CB: &model.CertifyBadInputSpec{
						Justification: "test justification",
					},
				},
			},
			ExpNeighbors: map[string][]string{
				"2":  []string{"2"},
				"3":  []string{"2", "2"},
				"4":  []string{"2", "2", "10"}, // pkg name
				"5":  []string{"2"},            // pkg version
				"6":  []string{"6"},
				"7":  []string{"6", "6"},
				"8":  []string{"6", "11"}, // src name
				"9":  []string{"12"},      // art
				"10": []string{"2"},       // cb 1 -> pkg name
				"11": []string{"6"},       // cb 2 -> src name
				"12": []string{"9"},       // cb 3 -> art
			},
		},
	}
	ctx := context.Background()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			b, err := inmem.GetBackend(nil)
			if err != nil {
				t.Fatalf("Could not instantiate testing backend: %v", err)
			}
			for _, p := range test.InPkg {
				if _, err := b.IngestPackage(ctx, *p); err != nil {
					t.Fatalf("Could not ingest package: %v", err)
				}
			}
			for _, s := range test.InSrc {
				if _, err := b.IngestSource(ctx, *s); err != nil {
					t.Fatalf("Could not ingest source: %v", err)
				}
			}
			for _, a := range test.InArt {
				if _, err := b.IngestArtifact(ctx, a); err != nil {
					t.Fatalf("Could not ingest artifact: %v", err)
				}
			}
			for _, o := range test.Calls {
				if _, err := b.IngestCertifyBad(ctx, o.Sub, o.Match, *o.CB); err != nil {
					t.Fatalf("Could not ingest CertifyBad: %v", err)
				}
			}
			for q, r := range test.ExpNeighbors {
				got, err := b.Neighbors(ctx, q, nil)
				if err != nil {
					t.Fatalf("Could not query neighbors: %s", err)
				}
				gotIDs := convNodes(got)
				slices.Sort(r)
				slices.Sort(gotIDs)
				if diff := cmp.Diff(r, gotIDs); diff != "" {
					t.Errorf("Unexpected results. (-want +got):\n%s", diff)
				}
			}
		})
	}
}