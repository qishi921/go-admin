package controllers

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/user/admin-system/app/models"
)

// AuditController 审计日志控制器（增强版）
type AuditController struct {
	DB *orm.DB
}

// List 日志列表
func (ac *AuditController) List(c *ghttp.Context) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	// 构建查询
	query := orm.Query[models.OperationLog](ac.DB)

	// 筛选条件
	userID := c.Query("user_id", "")
	if userID != "" {
		query = query.Where("user_id", "=", userID)
	}

	action := c.Query("action", "")
	if action != "" {
		query = query.Where("action", "like", "%"+action+"%")
	}

	status := c.Query("status", "")
	if status != "" {
		query = query.Where("status", "=", status)
	}

	// 数据变更筛选
	resourceTable := c.Query("resource_table", "")
	if resourceTable != "" {
		query = query.Where("resource_table", "=", resourceTable)
	}

	changeType := c.Query("change_type", "")
	if changeType != "" {
		query = query.Where("change_type", "=", changeType)
	}

	// 时间范围
	startDate := c.Query("start_date", "")
	if startDate != "" {
		query = query.Where("created_at", ">=", startDate+" 00:00:00")
	}
	endDate := c.Query("end_date", "")
	if endDate != "" {
		query = query.Where("created_at", "<=", endDate+" 23:59:59")
	}

	query = query.OrderBy("created_at", "DESC")
	result, err := orm.Paginate[models.OperationLog](query, page, pageSize)
	if err != nil {
		c.Error(500, "查询失败")
		return
	}

	c.Success(result)
}

// Detail 日志详情（包含数据快照）
func (ac *AuditController) Detail(c *ghttp.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(400, "缺少日志ID")
		return
	}

	query := orm.Query[models.OperationLog](ac.DB).Where("id", "=", id)
	log, err := orm.First[models.OperationLog](query)
	if err != nil || log == nil {
		c.Error(404, "日志不存在")
		return
	}

	// 解析数据快照
	var oldData, newData map[string]any
	if log.OldData != "" {
		json.Unmarshal([]byte(log.OldData), &oldData)
	}
	if log.NewData != "" {
		json.Unmarshal([]byte(log.NewData), &newData)
	}

	// 计算变更差异
	diff := ac.calculateDiff(oldData, newData)

	c.Success(map[string]any{
		"log":      log,
		"old_data": oldData,
		"new_data": newData,
		"diff":     diff,
	})
}

// calculateDiff 计算数据差异
func (ac *AuditController) calculateDiff(oldData, newData map[string]any) []map[string]any {
	var diff []map[string]any

	if oldData == nil && newData == nil {
		return diff
	}

	// 创建操作
	if oldData == nil && newData != nil {
		for key, val := range newData {
			diff = append(diff, map[string]any{
				"field":     key,
				"action":    "added",
				"old_value": nil,
				"new_value": val,
			})
		}
		return diff
	}

	// 删除操作
	if oldData != nil && newData == nil {
		for key, val := range oldData {
			diff = append(diff, map[string]any{
				"field":     key,
				"action":    "removed",
				"old_value": val,
				"new_value": nil,
			})
		}
		return diff
	}

	// 更新操作
	// 检查旧数据中的字段
	for key, oldVal := range oldData {
		newVal, exists := newData[key]
		if !exists {
			diff = append(diff, map[string]any{
				"field":     key,
				"action":    "removed",
				"old_value": oldVal,
				"new_value": nil,
			})
		} else if !ac.equalValues(oldVal, newVal) {
			diff = append(diff, map[string]any{
				"field":     key,
				"action":    "changed",
				"old_value": oldVal,
				"new_value": newVal,
			})
		}
	}

	// 检查新增的字段
	for key, newVal := range newData {
		if _, exists := oldData[key]; !exists {
			diff = append(diff, map[string]any{
				"field":     key,
				"action":    "added",
				"old_value": nil,
				"new_value": newVal,
			})
		}
	}

	return diff
}

// equalValues 比较两个值是否相等
func (ac *AuditController) equalValues(a, b any) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}

// Stats 统计信息
func (ac *AuditController) Stats(c *ghttp.Context) {
	// 今日操作数
	today := time.Now().Format("2006-01-02")
	todayQuery := orm.Query[models.OperationLog](ac.DB).Where("created_at", "like", today+"%")
	todayLogs, _ := orm.Get[models.OperationLog](todayQuery)
	todayCount := len(todayLogs)

	// 按操作类型统计
	actionStats := make(map[string]int)
	for _, log := range todayLogs {
		actionStats[log.Action]++
	}

	// 按状态统计
	successCount := 0
	failCount := 0
	for _, log := range todayLogs {
		if log.Status == "success" {
			successCount++
		} else {
			failCount++
		}
	}

	// 按变更类型统计
	changeStats := make(map[string]int)
	for _, log := range todayLogs {
		if log.ChangeType != "" {
			changeStats[log.ChangeType]++
		}
	}

	// 活跃用户数
	userSet := make(map[int]bool)
	for _, log := range todayLogs {
		if log.UserId != nil {
			userSet[*log.UserId] = true
		}
	}

	c.Success(map[string]any{
		"today_count":   todayCount,
		"success_count": successCount,
		"fail_count":    failCount,
		"action_stats":  actionStats,
		"change_stats":  changeStats,
		"active_users":  len(userSet),
	})
}

// DataChangeHistory 数据变更历史
func (ac *AuditController) DataChangeHistory(c *ghttp.Context) {
	resourceTable := c.Query("resource_table", "")
	recordID := c.Query("record_id", "")

	if resourceTable == "" || recordID == "" {
		c.Error(400, "请指定表名和记录ID")
		return
	}

	query := orm.Query[models.OperationLog](ac.DB).
		Where("resource_table", "=", resourceTable).
		Where("record_id", "=", recordID).
		OrderBy("created_at", "DESC")

	logs, err := orm.Get[models.OperationLog](query)
	if err != nil {
		c.Error(500, "查询失败")
		return
	}

	// 构建变更历史
	var history []map[string]any
	for _, log := range logs {
		var oldData, newData map[string]any
		if log.OldData != "" {
			json.Unmarshal([]byte(log.OldData), &oldData)
		}
		if log.NewData != "" {
			json.Unmarshal([]byte(log.NewData), &newData)
		}

		history = append(history, map[string]any{
			"id":          log.ID,
			"change_type": log.ChangeType,
			"old_data":    oldData,
			"new_data":    newData,
			"diff":        ac.calculateDiff(oldData, newData),
			"user_id":     log.UserId,
			"username":    log.Username,
			"created_at":  log.CreatedAt,
		})
	}

	c.Success(map[string]any{
		"resource_table": resourceTable,
		"record_id":      recordID,
		"history":        history,
	})
}

// Export 导出日志
func (ac *AuditController) Export(c *ghttp.Context) {
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")

	query := orm.Query[models.OperationLog](ac.DB)

	if startDate != "" {
		query = query.Where("created_at", ">=", startDate+" 00:00:00")
	}
	if endDate != "" {
		query = query.Where("created_at", "<=", endDate+" 23:59:59")
	}

	query = query.OrderBy("created_at", "DESC")
	logs, err := orm.Get[models.OperationLog](query)
	if err != nil {
		c.Error(500, "查询失败")
		return
	}

	// 生成 CSV
	var csvData strings.Builder
	csvData.WriteString("ID,用户ID,用户名,操作,方法,路径,IP,状态,耗时(ms),创建时间,表名,记录ID,变更类型\n")

	for _, log := range logs {
		recordID := ""
		if log.RecordId != nil {
			recordID = strconv.FormatUint(*log.RecordId, 10)
		}
		csvData.WriteString(strings.Join([]string{
			strconv.FormatUint(log.ID, 10),
			intToStr(log.UserId),
			log.Username,
			log.Action,
			log.Method,
			log.Path,
			log.Ip,
			log.Status,
			strconv.Itoa(log.Duration),
			log.CreatedAt,
			log.ResourceTable,
			recordID,
			log.ChangeType,
		}, ",") + "\n")
	}

	c.SetHeader("Content-Type", "text/csv; charset=utf-8")
	c.SetHeader("Content-Disposition", "attachment; filename=audit_logs.csv")

	// 添加 BOM
	bom := []byte{0xEF, 0xBB, 0xBF}
	c.Writer.Write(bom)
	c.Writer.Write([]byte(csvData.String()))
}

func intToStr(val *int) string {
	if val == nil {
		return ""
	}
	return strconv.Itoa(*val)
}
