package controllers

import (
	"net/http"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/user/admin-system/app/models"
)

// OperationLogController handles CRUD operations for OperationLog.
type OperationLogController struct {
	DB *orm.DB
}

// NewOperationLogController creates a new controller instance.
func NewOperationLogController(db *orm.DB) *OperationLogController {
	return &OperationLogController{DB: db}
}

// Index lists all operation logs with pagination and search.
func (ctrl *OperationLogController) Index(c *ghttp.Context) {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	search := c.Query("search")

	q := orm.Query[models.OperationLog](ctrl.DB)
	if search != "" {
		q = q.Where("action", "LIKE", "%"+search+"%")
	}
	q = q.OrderBy("created_at", "DESC")

	result, err := orm.Paginate[models.OperationLog](q, page, perPage)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.Success(result)
}

// Show returns a single operationlog by ID.
func (ctrl *OperationLogController) Show(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.OperationLog](
		orm.Query[models.OperationLog](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "OperationLog not found")
		return
	}
	c.Success(item)
}

// Ensure OperationLogController satisfies the ResourceController interface.
var _ interface {
	Index(c *ghttp.Context)
	Show(c *ghttp.Context)
} = (*OperationLogController)(nil)
