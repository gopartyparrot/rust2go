package rust2go

import "testing"

func TestParseU64F64(t *testing.T) {
	type args struct {
		u128 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "0.05 apy",
			args:    args{"14623560433"},
			want:    "0.00000000079274479955",
			wantErr: false,
		},
		{
			name:    "123.4",
			args:    args{"2276328218695758774272"},
			want:    "123.40000000000000568434",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseU64F64(tt.args.u128)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseU64F64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseU64F64() = %v, want %v", got, tt.want)
			}
		})
	}
}
