package bucket

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Upsert(t *testing.T) {
	type args struct {
		ctx    context.Context
		typ    string
		id     string
		oldDoc interface{}
		newDoc interface{}
		ttl    uint32
	}
	tests := []struct {
		name       string
		args       args
		withInsert bool
		wantErr    bool
	}{
		{
			name: "Standard",
			args: args{
				ctx:    context.Background(),
				typ:    "webshop",
				oldDoc: generate(),
				newDoc: generate(),
				ttl:    0,
			},
			withInsert: true,
			wantErr:    false,
		},
		{
			name: "Upsert new create",
			args: args{
				ctx:    context.Background(),
				typ:    "webshop",
				newDoc: generate(),
				ttl:    0,
			},
			withInsert: false,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.withInsert {
				var err error
				_, tt.args.id, err = th.EInsert(tt.args.ctx, tt.args.typ, "", tt.args.oldDoc, tt.args.ttl)
				assert.NoErrorf(t, err, "insert error")

				got := &webshop{}
				assert.NoError(t, th.Get(tt.args.ctx, tt.args.typ, tt.args.id, got))
				assert.Equal(t, tt.args.oldDoc, *got)
			}

			var err error
			if _, tt.args.id, err = th.Upsert(tt.args.ctx, tt.args.typ, tt.args.id, tt.args.newDoc, tt.args.ttl); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			got := &webshop{}
			assert.NoError(t, th.Get(tt.args.ctx, tt.args.typ, tt.args.id, got))
			assert.Equal(t, tt.args.newDoc, *got)
		})
	}
}
