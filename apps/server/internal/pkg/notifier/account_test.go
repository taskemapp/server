package notifier

import (
	"github.com/google/uuid"
	"testing"
)

func Test_buildConfirmLink(t *testing.T) {
	type args struct {
		host      string
		confirmID string
	}
	tt := struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{

		"Basic link",
		args{
			"localhost",
			uuid.Nil.String(),
		},
		"localhost/verify?id=00000000-0000-0000-0000-000000000000",
		false,
	}

	t.Run(tt.name, func(t *testing.T) {
		t.Parallel()
		got, err := buildConfirmLink(tt.args.host, tt.args.confirmID)
		if (err != nil) != tt.wantErr {
			t.Errorf("buildConfirmLink() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if got != tt.want {
			t.Errorf("buildConfirmLink() got = %v, want %v", got, tt.want)
		}
	})
}
