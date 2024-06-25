package git

import (
	"reflect"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/google/go-github/v45/github"
)

func Test_filterRemoteTags(t *testing.T) {
	tagFormat := "v%major%.%minor%.%patch%"
	tests := []struct {
		name      string
		refs      []*github.Reference
		tagFormat string
		tagPrefix string
		want      semver.Version
	}{
		{
			name:      "sample",
			tagFormat: tagFormat,
			tagPrefix: "",
			refs: []*github.Reference{
				{Ref: github.String("refs/tags/v0.5.5")},
				{Ref: github.String("refs/tags/v1.2.3")},
			},
			want: semver.Version{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:      "sampleA",
			tagFormat: tagFormat,
			tagPrefix: "moduleA-",
			refs: []*github.Reference{
				{Ref: github.String("refs/tags/moduleA-v0.5.5")},
				{Ref: github.String("refs/tags/moduleB-v1.2.3")},
			},
			want: semver.Version{Major: 0, Minor: 5, Patch: 5},
		},
		{
			name:      "sampleB",
			tagFormat: tagFormat,
			tagPrefix: "moduleB-",
			refs: []*github.Reference{
				{Ref: github.String("refs/tags/moduleA-v0.5.5")},
				{Ref: github.String("refs/tags/moduleB-v1.2.3")},
			},
			want: semver.Version{Major: 1, Minor: 2, Patch: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterRemoteTags(tt.refs, tt.tagFormat, tt.tagPrefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterRemoteTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
