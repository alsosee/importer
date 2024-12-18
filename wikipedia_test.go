package main

import (
	"testing"
)

func TestParseInfobox(t *testing.T) {
	tests := []struct {
		name string
		text string
		want map[string]string
	}{
		{
			name: "empty",
			text: "",
			want: map[string]string{},
		},
		{
			name: "simple",
			text: `
<table class="infobox vevent"><tbody>
<tr>
    <th scope="row" class="infobox-label" style="white-space: nowrap; padding-right: 0.65em;">Directed by</th>
    <td class="infobox-data"><a href="/wiki/Paul_McGuigan_(filmmaker)" class="mw-redirect" title="Paul McGuigan (filmmaker)">Paul McGuigan</a></td>
</tr>
</tbody></table>`,
			want: map[string]string{
				"Directed by": "Paul McGuigan",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInfobox(tt.text)
			if err != nil {
				t.Fatalf("%q ParseInfobox() error = %v", tt.name, err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("%q ParseInfobox() = %v, want %v", tt.name, got, tt.want)
			}

			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("%q ParseInfobox() = %v, want %v", tt.name, got, tt.want)
				}
			}
		})
	}
}
