package model

import (
	"fmt"
	"time"

	"github.com/kyma-incubator/reconciler/pkg/db"
)

const tblCacheDeps string = "config_cachedeps"

type CacheDependencyEntity struct {
	Bucket  string    `db:"notNull"`
	Key     string    `db:"notNull"`
	Label   string    `db:"notNull"`
	Cluster string    `db:"notNull"`
	CacheID int64     `db:"notNull"`
	Created time.Time `db:"readOnly"`
}

func (cde *CacheDependencyEntity) String() string {
	return fmt.Sprintf("CacheDependencyEntity [Bucket=%s,Key=%s,Label=%s,Cluster=%s,CacheID=%d]",
		cde.Bucket, cde.Key, cde.Label, cde.Cluster, cde.CacheID)
}

func (cde *CacheDependencyEntity) New() db.DatabaseEntity {
	return &CacheDependencyEntity{}
}

func (cde *CacheDependencyEntity) Marshaller() *db.EntityMarshaller {
	marshaller := db.NewEntityMarshaller(&cde)
	marshaller.AddUnmarshaller("Created", convertTimestampToTime)
	return marshaller
}

func (cde *CacheDependencyEntity) Table() string {
	return tblCacheDeps
}

func (cde *CacheDependencyEntity) Equal(other db.DatabaseEntity) bool {
	if other == nil {
		return false
	}
	otherDep, ok := other.(*CacheDependencyEntity)
	if ok {
		return cde.Bucket == otherDep.Bucket &&
			cde.Key == otherDep.Key &&
			cde.Label == otherDep.Label &&
			cde.Cluster == otherDep.Cluster
	}
	return false
}
