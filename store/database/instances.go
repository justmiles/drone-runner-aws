package database

import (
	"context"

	"github.com/drone-runners/drone-runner-aws/store"
	"github.com/drone-runners/drone-runner-aws/types"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var _ store.InstanceStore = (*InstanceStore)(nil)

func NewInstanceStore(db *sqlx.DB) *InstanceStore {
	return &InstanceStore{db}
}

type InstanceStore struct {
	db *sqlx.DB
}

func (s InstanceStore) Find(_ context.Context, id string) (*types.Instance, error) {
	dst := new(types.Instance)
	err := s.db.Get(dst, instanceFindByID, id)
	return dst, err
}

func (s InstanceStore) List(_ context.Context, pool string, params *types.QueryParams) ([]*types.Instance, error) {
	dst := []*types.Instance{}
	var args []interface{}

	stmt := builder.Select("*").From("instances").Where(squirrel.Eq{"instance_pool": pool})
	args = append(args, pool)
	if params != nil {
		if params.Stage != "" {
			stmt = stmt.Where(squirrel.Eq{"instance_stage": params.Stage})
			args = append(args, params.Stage)
		}
		if params.Status != "" {
			stmt = stmt.Where(squirrel.Eq{"instance_state": params.Status})
			args = append(args, params.Status)
		}
	}
	stmt = stmt.OrderBy("instance_started " + "ASC")
	sql, _, _ := stmt.ToSql()
	var err = s.db.Select(&dst, sql, args...)
	return dst, err
}

func (s InstanceStore) Create(_ context.Context, instance *types.Instance) error {
	query, arg, err := s.db.BindNamed(instanceInsert, instance)
	if err != nil {
		return err
	}
	return s.db.QueryRow(query, arg...).Scan(&instance.ID)
}

func (s InstanceStore) Delete(ctx context.Context, id string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint
	if _, err := tx.Exec(instanceDelete, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s InstanceStore) Update(_ context.Context, instance *types.Instance) error {
	query, arg, err := s.db.BindNamed(instanceUpdate, instance)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(query, arg...)
	return err
}

func (s InstanceStore) Purge(ctx context.Context) error {
	panic("implement me")
}

const instanceBase = `
SELECT
 instance_name
,instance_id
,instance_address
,instance_provider
,instance_state
,instance_pool
,instance_image
,instance_region
,instance_zone
,instance_size
,instance_platform
,instance_arch
,instance_stage
,instance_ca_key
,instance_ca_cert
,instance_tls_key
,instance_tls_cert
,instance_started
,instance_updated
,is_hibernated
FROM instances
`

const instanceFindByID = instanceBase + `
WHERE instance_id = $1
`

const instanceInsert = `
INSERT INTO instances (
 instance_id
,instance_name
,instance_address
,instance_provider
,instance_state
,instance_pool
,instance_image
,instance_region
,instance_zone
,instance_size
,instance_platform
,instance_arch
,instance_stage
,instance_ca_key
,instance_ca_cert
,instance_tls_key
,instance_tls_cert
,instance_started
,instance_updated
,is_hibernated
) values (
 :instance_id
,:instance_name
,:instance_address
,:instance_provider
,:instance_state
,:instance_pool
,:instance_image
,:instance_region
,:instance_zone
,:instance_size
,:instance_platform
,:instance_arch
,:instance_stage
,:instance_ca_key
,:instance_ca_cert
,:instance_tls_key
,:instance_tls_cert
,:instance_started
,:instance_updated
,:is_hibernated
) RETURNING instance_id
`

const instanceDelete = `
DELETE FROM instances
WHERE instance_id = $1
`

const instanceUpdate = `
UPDATE instances
SET
  instance_state    = :instance_state
 ,instance_stage	= :instance_stage
 ,instance_updated  = :instance_updated
 ,is_hibernated 	= :is_hibernated
 ,instance_address  = :instance_address
WHERE instance_id   = :instance_id
`
