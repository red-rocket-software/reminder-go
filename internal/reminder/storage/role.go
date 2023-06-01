package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/red-rocket-software/reminder-go/pkg/postgresql"
)

// GetUserPermissions return users permissions to role form DB PostgresSQL
func GetUserPermissions(ctx context.Context, role string, cfg config.Config) ([]string, error) {
	postgresClient, err := postgresql.NewClient(ctx, 5, cfg)
	if err != nil {
		log.Fatalf("Error create new db client:%v\n", err)
	}
	defer postgresClient.Close()

	var routes []string
	const sql = `SELECT sf.name
FROM role.role_permissions rp
JOIN role.permissions p ON p.id = ANY (rp.permissions)
JOIN role.sub_features sf ON sf.id = ANY (p.sub_features)
WHERE rp.role = $1`
	rows, err := postgresClient.Query(ctx, sql, role)

	defer rows.Close()

	if err != nil {
		return []string{}, fmt.Errorf("error get user permissions: %v", err)
	}

	for rows.Next() {
		var route string
		if err := rows.Scan(&route); err != nil {
			return []string{}, fmt.Errorf("error get user permissions: %v", err)
		}
		routes = append(routes, route)
	}

	if err := rows.Err(); err != nil {
		return []string{}, fmt.Errorf("error get user permissions: %v", err)
	}

	return routes, nil
}
