package permission

import (
	"fmt"
	"octopus/config"
)

var (
	// PSRMISSIONMONIT 查看监控界面
	PSRMISSIONMONIT = 1 << 0
	// PERMISSIONDEV dev 运维模式
	PERMISSIONDEV = 1 << 1
	// PERMISSIONEXEC 在线操作模式
	PERMISSIONEXEC = 1 << 2
)

var userGroup []*User

// User ...
type User struct {
	Username   string
	Password   string
	Permission int
	Token      string
}

func init() {
	userGroup = make([]*User, 0)
	for _, v := range config.C.Auth {
		if len(userGroup) <= 100 {
			var permission = 0
			for _, p := range v.Permission {
				switch p {
				case "monit":
					permission |= PSRMISSIONMONIT
				case "dev":
					permission |= PERMISSIONDEV
				case "exec":
					permission |= PERMISSIONEXEC
				}
			}
			set(&User{
				Username:   v.User,
				Password:   v.Password,
				Permission: permission,
			})
		}
	}
}
func set(u *User) {
	for _, v := range userGroup {
		if v.Username == u.Username {
			return
		}
	}
	key := config.C.AuthConfig.Key
	token := fmt.Sprintf("%x", key+"|"+u.Username)
	u.Token = token
	userGroup = append(userGroup, u)
}

// Get ...
func Get(token string) *User {
	for _, v := range userGroup {
		if v.Token == token {
			return v
		}
	}
	return nil
}
