package auth

type SignInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DepartmentInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PositionInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type EmployeeInfo struct {
	Id           string          `json:"id"`
	Name         string          `json:"name"`
	DepartmentId *string         `json:"department_id"`
	Department   *DepartmentInfo `json:"department"`
	PositionId   *string         `json:"position_id"`
	Position     *PositionInfo   `json:"position"`
}

type ExternalUser struct {
	Id             string        `json:"id"`
	Name           string        `json:"name"`
	Username       string        `json:"username"`
	Email          string        `json:"email"`
	Role           string        `json:"role"`
	IsActive       bool          `json:"is_active"`
	CompanyId      *string       `json:"company_id"`
	BranchId       *string       `json:"branch_id"`
	EmployeeId     *string       `json:"employee_id"`
	Employee       *EmployeeInfo `json:"employee"`
	ProfilePicture *string       `json:"profile_picture"`
}

type ExternalAuthData struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         ExternalUser `json:"user"`
}

type ExternalAuthResponse struct {
	Data    *ExternalAuthData `json:"data"`
	Error   string            `json:"error,omitempty"`
	Message string            `json:"message,omitempty"`
}

type SignInResponse struct {
	AccessToken    string       `json:"access_token"`
	RefreshToken   string       `json:"refresh_token"`
	User           ExternalUser `json:"user"`
	CitramascoRole string       `json:"citramasco_role"`
	Department     string       `json:"department,omitempty"`
	Position       string       `json:"position,omitempty"`
	Permissions    []string     `json:"permissions"`
}

type ExternalShiftMaster struct {
	Id          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	IsOvernight bool   `json:"is_overnight"`
}

type ExternalShiftRoster struct {
	Id            string               `json:"id"`
	EmployeeId    string               `json:"employee_id"`
	Date          string               `json:"date"`
	IsDayOff      bool                 `json:"is_day_off"`
	ShiftMasterId *string              `json:"shift_master_id"`
	ShiftMaster   *ExternalShiftMaster `json:"shift_master"`
}

type ExternalRosterResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Data []ExternalShiftRoster `json:"data"`
	} `json:"data"`
	Message string `json:"message"`
}
