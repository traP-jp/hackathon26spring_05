package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
)

// 事前計算済みプロフィールを取得する。
func (r *repositoryImpl) FindPrecomputedProfileByUsername(username string) (*domain.UserOverride, error) {
	var dataJSON []byte
	err := r.db.Get(&dataJSON, `SELECT data_json FROM precompiled_users WHERE username = ?`, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var profile domain.UserOverride
	if err := json.Unmarshal(dataJSON, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
