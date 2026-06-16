package models

// User represents the users table.
type User struct {
	BaseModel
	Username    string  `json:"username" gai:"column:username;size:50"`
	Password    string  `json:"-"       gai:"column:password;size:255"`
	Email       string  `json:"email"   gai:"column:email;size:100;unique"`
	Phone       *string `json:"phone"   gai:"column:phone;size:20;nullable"`
	Avatar      *string `json:"avatar"  gai:"column:avatar;size:255;nullable"`
	RealName    *string `json:"real_name" gai:"column:real_name;size:50;nullable"`
	Status      string  `json:"status"  gai:"column:status"`
	LastLoginAt *string `json:"last_login_at" gai:"column:last_login_at;nullable"`
	RoleId      *int    `json:"role_id" gai:"column:role_id;nullable"`
	Role        *Role   `json:"role,omitempty" gai:"belongsTo"`
}

// TableName returns the database table name.
func (u *User) TableName() string {
	return "users"
}