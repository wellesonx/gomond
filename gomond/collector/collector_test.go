package collector

import (
	"github.com/gelleson/gomond/gomond/parser"
	"reflect"
	"sync"
	"testing"
	"time"
)

func Test_isActive(t *testing.T) {
	type args struct {
		ttl  time.Time
		meta logMeta
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true case",
			args: args{
				ttl: time.Now().Add(-time.Second * 10),
				meta: logMeta{
					created: time.Now(),
				},
			},
			want: true,
		},
		{
			name: "false case",
			args: args{
				ttl: time.Now().Add(-time.Second * 10),
				meta: logMeta{
					created: time.Now().Add(-time.Second * 12),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isActive(tt.args.ttl, tt.args.meta); got != tt.want {
				t.Errorf("isActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryLogCollector_filterByTTL(t *testing.T) {
	type fields struct {
		logs   []logMeta
		viewed int64
		option MemoryOption
		mutex  *sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		total  int
	}{
		{
			name: "filtered",
			fields: fields{
				logs: []logMeta{
					{
						app:     "1",
						created: time.Now(),
					},
					{
						app:     "2",
						created: time.Now(),
					},
					{
						app:     "3",
						created: time.Now().Add(-time.Second * 12),
					},
				},
				viewed: 0,
				option: MemoryOption{
					TTL: time.Second * 10,
				},
				mutex: &sync.Mutex{},
			},
			total: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryLogCollector{
				logs:   tt.fields.logs,
				viewed: tt.fields.viewed,
				option: tt.fields.option,
				mutex:  tt.fields.mutex,
			}

			m.filterByTTL()

			if tt.total != len(m.logs) {
				t.Errorf("expexted %d but actual value %d", tt.total, len(m.logs))
			}
		})
	}
}

func TestMemoryLogCollector_Get(t *testing.T) {
	firstCall := []logMeta{
		{
			app:     "1",
			created: time.Now(),
		},
		{
			app:     "2",
			created: time.Now(),
		},
		{
			app:     "3",
			created: time.Now().Add(-time.Second * 12),
		},
	}

	type fields struct {
		logs   []logMeta
		viewed int64
		option MemoryOption
		mutex  *sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "filtered",
			fields: fields{
				logs:   firstCall,
				viewed: 0,
				option: MemoryOption{
					TTL: time.Second * 10,
				},
				mutex: &sync.Mutex{},
			},
			want: len(firstCall),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryLogCollector{
				logs:   tt.fields.logs,
				viewed: tt.fields.viewed,
				option: tt.fields.option,
				mutex:  tt.fields.mutex,
			}
			if got := m.Get(); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("Get() = %v, want %v", len(got), tt.want)
			}

			m.Push(parser.Log{})

			if got := m.Get(); !reflect.DeepEqual(len(got), 1) {
				t.Errorf("Get() = %v, want %v", len(got), 1)
			}
		})
	}
}
