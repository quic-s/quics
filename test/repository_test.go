package test

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quic-s/quics/pkg/common"
	"reflect"
	"testing"
)

func TestGetAll(t *testing.T) {
	tests := []struct {
		want []common.Command
	}{
		{
			want: []common.Command{
				{[]byte("hello"), []byte("world")},
				{[]byte("milk"), []byte("cereal")},
			},
		},
	}

	for _, tc := range tests {
		db, err := InitInMemoryDB()
		if err != nil {
			t.Fatalf("unable to init in-memory DB: %v", err)
		}
		defer db.Close()
		c := common.CommandsRepository{db}

		for _, cmd := range tc.want {
			if err := c.SetValue(cmd.Key, cmd.Value); err != nil {
				t.Fatalf("unable to set value: %v", err)
			}
		}

		got, err := c.GetAll()
		if err != nil {
			t.Fatalf("unable to get commands: %v", err)
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("got != want;\n%v != %v", got, tc.want)
		}
	}
}

func TestGetValue(t *testing.T) {
	tests := []struct {
		input []common.Command
		key   []byte
		want  []byte
	}{
		{
			input: []common.Command{
				{[]byte("hello"), []byte("world")},
				{[]byte("milk"), []byte("cereal")},
			},
			key:  []byte("hello"),
			want: []byte("world"),
		},
		{
			input: []common.Command{
				{[]byte("hello"), []byte("world")},
				{[]byte("milk"), []byte("cereal")},
			},
			key:  []byte("i don't exist"),
			want: []byte("this command does not exist"),
		},
	}

	for _, tc := range tests {
		db, err := InitInMemoryDB()
		if err != nil {
			t.Fatalf("unable to init in-memory DB: %v", err)
		}
		defer db.Close()
		c := common.CommandsRepository{db}

		for _, cmd := range tc.input {
			if err := c.SetValue(cmd.Key, cmd.Value); err != nil {
				t.Fatalf("unable to set value: %v", err)
			}
		}

		got, err := c.GetValue(tc.key)
		if err != nil {
			t.Fatalf("unable to get commands: %v", err)
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("got != want;\n%v != %v", got, tc.want)
		}
	}
}

func TestSetValue(t *testing.T) {
	tests := []struct {
		input    []common.Command
		key      []byte
		newValue []byte
	}{
		{
			input: []common.Command{
				{[]byte("hello"), []byte("world")},
				{[]byte("milk"), []byte("cereal")},
			},
			key:      []byte("hello"),
			newValue: []byte("chat"),
		},
		{
			input: []common.Command{
				{[]byte("hello"), []byte("world")},
				{[]byte("milk"), []byte("cereal")},
			},
			key:      []byte("new"),
			newValue: []byte("value"),
		},
	}

	for _, tc := range tests {
		db, err := InitInMemoryDB()
		if err != nil {
			t.Fatalf("unable to init in-memory DB: %v", err)
		}
		defer db.Close()
		c := common.CommandsRepository{db}

		for _, cmd := range tc.input {
			if err := c.SetValue(cmd.Key, cmd.Value); err != nil {
				t.Fatalf("unable to set value: %v", err)
			}
		}
		if err := c.SetValue(tc.key, tc.newValue); err != nil {
			t.Fatalf("unable to set value: %v", err)
		}
		got, err := c.GetValue(tc.key)
		if err != nil {
			t.Fatalf("unable to get commands: $v, err")
		}
		if !reflect.DeepEqual(got, tc.newValue) {
			t.Fatalf("got != want;\n%v != %v", got, tc.newValue)
		}
	}
}

func TestDeleteValue(t *testing.T) {
	tests := []struct {
		input []common.Command
		key   []byte
	}{
		{
			input: []common.Command{
				{[]byte("hello"), []byte("world")},
				{[]byte("milk"), []byte("cereal")},
			},
			key: []byte("hello"),
		},
		{
			input: []common.Command{
				{[]byte("hello"), []byte("world")},
				{[]byte("milk"), []byte("cereal")},
			},
			key: []byte("new"),
		},
	}

	for _, tc := range tests {
		db, err := InitInMemoryDB()
		if err != nil {
			t.Fatalf("unable to init in-memory DB: %v", err)
		}
		defer db.Close()
		c := common.CommandsRepository{db}

		for _, cmd := range tc.input {
			if err := c.SetValue(cmd.Key, cmd.Value); err != nil {
				t.Fatalf("unable to set value: %v", err)
			}
		}
		if err := c.DeleteValue(tc.key); err != nil {
			t.Fatalf("unable to delete value: %v", err)
		}
		if v, _ := c.GetValue(tc.key); string(v) != "this command does not exist" {
			t.Fatalf("command still exists: %v", err)
		}
	}
}

func InitInMemoryDB() (*badger.DB, error) {
	opt := badger.DefaultOptions("").WithInMemory(true).WithLogger(nil)
	return badger.Open(opt)
}
