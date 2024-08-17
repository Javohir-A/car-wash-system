package casbin

// import (
// 	"gateway/config"
// 	"log"

// 	pgadapter "github.com/Blank-Xu/sql-adapter"
// 	"github.com/casbin/casbin/v2"
// )

// func CasbinEnforcer(cfg *config.Config) (*casbin.Enforcer, error) {
// 	db, err := ConnectDB(cfg)
// 	if err != nil {
// 		log.Printf("failed to connect to the database: %v", err)
// 		return nil, err
// 	}
// 	defer db.Close()

// 	a, err := pgadapter.NewAdapter(db, "postgres", "casbin_rule")
// 	if err != nil {
// 		log.Printf("failed to construct adapter: %v", err)
// 		return nil, err
// 	}

// 	enforcer, err := casbin.NewEnforcer("casbin/model.conf", a)
// 	if err != nil {
// 		log.Printf("failed to construct enforcer: %v", err)
// 		return nil, err
// 	}

// 	err = enforcer.LoadPolicy()
// 	if err != nil {
// 		log.Printf("failed to load policy: %v", err)
// 		return nil, err
// 	}

// 	enforcer.ClearPolicy()

// 	policies := [][]string{
// 		{"admin", "/car-wash/*", "*", "allow"},
// 		{"provider", "/car-wash/*", "*", "allow"},

// 		{"customer", "/car-wash/providers", "POST", "deny"},
// 		{"customer", "/car-wash/providers/*", "PUT", "deny"},
// 		{"customer", "/car-wash/providers/*", "DELETE", "deny"},

// 		{"customer", "/car-wash/services", "POST", "deny"},
// 		{"customer", "/car-wash/services/*", "PUT", "deny"},
// 		{"customer", "/car-wash/services/*", "DELETE", "deny"},

// 		{"customer", "/car-wash/notifications", "POST", "deny"},

// 		{"customer", "/car-wash/*", "*", "allow"},
// 	}

// 	_, err = enforcer.AddPolicies(policies)
// 	if err != nil {
// 		log.Printf("failed to add policies: %v", err)
// 		return nil, err
// 	}

// 	err = enforcer.SavePolicy()
// 	if err != nil {
// 		log.Printf("failed to save policy: %v", err)
// 		return nil, err
// 	}

// 	return enforcer, nil
// }
