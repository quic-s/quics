package common

import (
	"github.com/dgraph-io/badger/v3"
)

type Command struct {
	Key   []byte
	Value []byte
}

type Repository interface {
	GetAll() ([]Command, error)
	GetValue([]byte) []byte
	SetValue(key, value []byte)
	EditValue([]byte)
	DeleteValue([]byte) bool
}

type CommandsRepository struct {
	DB *badger.DB
}

func (c *CommandsRepository) GetAll() ([]Command, error) {
	var cmds []Command
	err := c.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				cmds = append(cmds, Command{Key: k, Value: v})
				return nil
			})

			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cmds, nil
}

func (c *CommandsRepository) SetValue(k, v []byte) error {
	err := c.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(k, v)
		return err
	})
	return err
}

func (c *CommandsRepository) GetValue(k []byte) ([]byte, error) {
	var v []byte
	err := c.DB.View(func(txn *badger.Txn) error {
		i, err := txn.Get(k)

		if err != nil {
			v = []byte("This command does not exist")
			return nil
		}

		v, err = i.ValueCopy(v)
		return err
	})

	return v, err
}

func (c *CommandsRepository) DeleteValue(k []byte) error {
	err := c.DB.Update(func(txn *badger.Txn) error {
		err := txn.Delete(k)
		return err
	})
	return err
}
