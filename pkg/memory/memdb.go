package memory

/*
import (
	"github.com/hashicorp/go-memdb"
	"golinkcut/internal/entity"
)

type Storage struct {
	db *memdb.MemDB
}

func NewStorage() *Storage {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"link": &memdb.TableSchema{
				Name: "link",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Short"},
					},
					"link_id": &memdb.IndexSchema{
						Name:    "link_id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Original"},
					},
				},
			},
		},
	}
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	return &Storage{db}
}

func (s *Storage) Add(link entity.Link) error {
	txn := s.db.Txn(true)
	if err := txn.Insert("link", link); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storage) Get(shortLink string) string {
	txn := s.db.Txn(false)
	defer txn.Abort()
	raw, err := txn.First("link", "id", shortLink)
	if err != nil {
		return ""
	}
	return raw.(entity.Link).Original
}
*/
