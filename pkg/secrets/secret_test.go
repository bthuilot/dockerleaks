package secrets

import (
	"github.com/bthuilot/dockerleaks/internal/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecret_String(t *testing.T) {
	before := viper.GetBool(config.ViperUnmaskKey)
	t.Cleanup(func() {
		viper.Set(config.ViperUnmaskKey, before)
	})
	tests := []struct {
		name      string
		unmaskSet bool
		s         Secret
		want      string
	}{
		{
			name: "empty secret",
			s:    Secret{},
			want: "",
		},
		{
			name:      "unmasked secret < 6",
			unmaskSet: true,
			s: Secret{
				Value: "<6",
			},
			want: "<6",
		},
		{
			name:      "unmasked secret > 6",
			unmaskSet: true,
			s: Secret{
				Value: "unmasked-secret",
			},
			want: "unmasked-secret",
		},
		{
			name: "masked secret < 6",
			s: Secret{
				Value: "<6",
			},
			want: "*****",
		},
		{
			name: "masked secret > 6",
			s: Secret{
				Value: "masked-secret",
			},
			want: "mas*****ret",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(config.ViperUnmaskKey, tt.unmaskSet)
			assert.Equal(t, tt.want, tt.s.String())
		})
	}
}
