package controllers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
)

// 允许导出/导入的表名白名单
var allowedTables = map[string]bool{
	"users":           true,
	"roles":           true,
	"menus":           true,
	"permissions":     true,
	"operation_logs":  true,
	"notifications":   true,
	"scheduled_tasks": true,
	"uploads":         true,
	"settings":        true,
}

// validateTable 验证表名是否在白名单中
func validateTable(table string) error {
	if !allowedTables[table] {
		return fmt.Errorf("不允许操作表: %s", table)
	}
	return nil
}

// validateColumnName 验证列名是否安全（只允许字母、数字、下划线）
func validateColumnName(col string) error {
	for _, c := range col {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return fmt.Errorf("无效的列名: %s", col)
		}
	}
	return nil
}

// ExportController 导出控制器
type ExportController struct {
	DB *orm.DB
}

// ExportCSV 导出CSV
func (ec *ExportController) ExportCSV(c *ghttp.Context) {
	var req struct {
		Table    string   `json:"table"`     // 表名
		Fields   []string `json:"fields"`    // 导出字段
		Filters  string   `json:"filters"`   // JSON 筛选条件
		OrderBy  string   `json:"order_by"`  // 排序字段
		OrderDir string   `json:"order_dir"` // 排序方向 ASC/DESC
	}
	if err := c.BindJSON(&req); err != nil {
		c.Error(400, "参数错误")
		return
	}

	if req.Table == "" {
		c.Error(400, "请指定导出表")
		return
	}

	// 验证表名
	if err := validateTable(req.Table); err != nil {
		c.Error(400, err.Error())
		return
	}

	// 根据表名获取数据
	data, err := ec.getExportData(req.Table, req.Filters, req.OrderBy, req.OrderDir)
	if err != nil {
		c.Error(500, "数据查询失败: "+err.Error())
		return
	}

	// 生成CSV
	csvData := ec.generateCSV(data, req.Fields)

	// 设置响应头
	c.SetHeader("Content-Type", "text/csv; charset=utf-8")
	c.SetHeader("Content-Disposition", "attachment; filename="+req.Table+"_export.csv")

	// 添加UTF-8 BOM以支持Excel正确显示中文
	bom := []byte{0xEF, 0xBB, 0xBF}
	c.Writer.Write(bom)
	c.Writer.Write(csvData)
}

// ExportJSON 导出JSON
func (ec *ExportController) ExportJSON(c *ghttp.Context) {
	var req struct {
		Table    string   `json:"table"`
		Fields   []string `json:"fields"`
		Filters  string   `json:"filters"`
		OrderBy  string   `json:"order_by"`
		OrderDir string   `json:"order_dir"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.Error(400, "参数错误")
		return
	}

	if req.Table == "" {
		c.Error(400, "请指定导出表")
		return
	}

	// 验证表名
	if err := validateTable(req.Table); err != nil {
		c.Error(400, err.Error())
		return
	}

	data, err := ec.getExportData(req.Table, req.Filters, req.OrderBy, req.OrderDir)
	if err != nil {
		c.Error(500, "数据查询失败: "+err.Error())
		return
	}

	// 筛选字段
	filteredData := ec.filterFields(data, req.Fields)

	// 生成JSON
	jsonData, _ := json.MarshalIndent(filteredData, "", "  ")

	c.SetHeader("Content-Type", "application/json; charset=utf-8")
	c.SetHeader("Content-Disposition", "attachment; filename="+req.Table+"_export.json")
	c.Writer.Write(jsonData)
}

// ImportCSV 导入CSV
func (ec *ExportController) ImportCSV(c *ghttp.Context) {
	table := c.Query("table", "")
	if table == "" {
		c.Error(400, "请指定导入表")
		return
	}

	// 验证表名
	if err := validateTable(table); err != nil {
		c.Error(400, err.Error())
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(400, "请上传CSV文件")
		return
	}
	defer file.Close()

	// 检查文件大小（最大10MB）
	if header.Size > 10*1024*1024 {
		c.Error(400, "文件大小不能超过10MB")
		return
	}

	// 检查文件扩展名
	filename := header.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".csv") {
		c.Error(400, "请上传CSV格式文件")
		return
	}

	// 解析CSV
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	csvData := buf.Bytes()

	// 去除BOM
	if len(csvData) >= 3 && csvData[0] == 0xEF && csvData[1] == 0xBB && csvData[2] == 0xBF {
		csvData = csvData[3:]
	}

	csvReader := csv.NewReader(bytes.NewReader(csvData))
	records, err := csvReader.ReadAll()
	if err != nil {
		c.Error(400, "CSV解析失败: "+err.Error())
		return
	}

	if len(records) < 2 {
		c.Error(400, "CSV文件至少需要包含标题行和一行数据")
		return
	}

	// 第一行是字段名
	headers := records[0]
	dataRows := records[1:]

	// 导入数据
	successCount := 0
	failCount := 0
	errors := []string{}

	for i, row := range dataRows {
		if len(row) != len(headers) {
			failCount++
			errors = append(errors, "第"+strconv.Itoa(i+2)+"行数据列数不匹配")
			continue
		}

		// 构建数据映射
		rowData := make(map[string]any)
		for j, header := range headers {
			rowData[header] = row[j]
		}

		// 插入数据
		if err := ec.insertData(table, rowData); err != nil {
			failCount++
			errors = append(errors, "第"+strconv.Itoa(i+2)+"行导入失败: "+err.Error())
		} else {
			successCount++
		}
	}

	c.Success(map[string]any{
		"success": successCount,
		"failed":  failCount,
		"total":   len(dataRows),
		"errors":  errors,
	})
}

// ImportJSON 导入JSON
func (ec *ExportController) ImportJSON(c *ghttp.Context) {
	table := c.Query("table", "")
	if table == "" {
		c.Error(400, "请指定导入表")
		return
	}

	// 验证表名
	if err := validateTable(table); err != nil {
		c.Error(400, err.Error())
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(400, "请上传JSON文件")
		return
	}
	defer file.Close()

	if header.Size > 10*1024*1024 {
		c.Error(400, "文件大小不能超过10MB")
		return
	}

	filename := header.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".json") {
		c.Error(400, "请上传JSON格式文件")
		return
	}

	// 解析JSON
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	var jsonData []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &jsonData); err != nil {
		c.Error(400, "JSON解析失败: "+err.Error())
		return
	}

	if len(jsonData) == 0 {
		c.Error(400, "JSON数据为空")
		return
	}

	// 导入数据
	successCount := 0
	failCount := 0
	errors := []string{}

	for i, rowData := range jsonData {
		if err := ec.insertData(table, rowData); err != nil {
			failCount++
			errors = append(errors, "第"+strconv.Itoa(i+1)+"条数据导入失败: "+err.Error())
		} else {
			successCount++
		}
	}

	c.Success(map[string]any{
		"success": successCount,
		"failed":  failCount,
		"total":   len(jsonData),
		"errors":  errors,
	})
}

// getExportData 获取导出数据
func (ec *ExportController) getExportData(table, filters, orderBy, orderDir string) ([]map[string]any, error) {
	// 构建查询
	query := "SELECT * FROM " + table
	var args []any

	// 解析筛选条件
	if filters != "" {
		var filterMap map[string]any
		if err := json.Unmarshal([]byte(filters), &filterMap); err == nil {
			whereClause, whereArgs := ec.buildWhereClause(filterMap)
			if whereClause != "" {
				query += " WHERE " + whereClause
				args = append(args, whereArgs...)
			}
		}
	}

	// 排序（验证列名）
	if orderBy != "" {
		if err := validateColumnName(orderBy); err == nil {
			orderDir = strings.ToUpper(orderDir)
			if orderDir != "ASC" && orderDir != "DESC" {
				orderDir = "ASC"
			}
			query += " ORDER BY " + orderBy + " " + orderDir
		}
	}

	// 执行查询
	rows, err := ec.DB.SQL.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 获取列信息
	columns, _ := rows.Columns()

	// 读取数据
	var results []map[string]any
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}

		row := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}

// buildWhereClause 构建WHERE条件（使用参数化方式）
func (ec *ExportController) buildWhereClause(filters map[string]any) (string, []any) {
	var conditions []string
	var values []any
	for key, val := range filters {
		// 验证列名
		if err := validateColumnName(key); err != nil {
			continue // 跳过无效列名
		}
		if strVal, ok := val.(string); ok {
			conditions = append(conditions, key+" = ?")
			values = append(values, strVal)
		} else if numVal, ok := val.(float64); ok {
			conditions = append(conditions, key+" = ?")
			values = append(values, numVal)
		} else if boolVal, ok := val.(bool); ok {
			conditions = append(conditions, key+" = ?")
			values = append(values, boolVal)
		}
	}
	return strings.Join(conditions, " AND "), values
}

// generateCSV 生成CSV数据
func (ec *ExportController) generateCSV(data []map[string]any, fields []string) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	buf := new(bytes.Buffer)
	csvWriter := csv.NewWriter(buf)

	// 确定列名
	var headers []string
	if len(fields) > 0 {
		headers = fields
	} else {
		// 使用第一行数据的所有键
		for key := range data[0] {
			headers = append(headers, key)
		}
	}

	// 写入标题行
	csvWriter.Write(headers)

	// 写入数据行
	for _, row := range data {
		var values []string
		for _, field := range headers {
			val := row[field]
			values = append(values, ec.formatValue(val))
		}
		csvWriter.Write(values)
	}

	csvWriter.Flush()
	return buf.Bytes()
}

// filterFields 筛选字段
func (ec *ExportController) filterFields(data []map[string]any, fields []string) []map[string]any {
	if len(fields) == 0 {
		return data
	}

	var result []map[string]any
	for _, row := range data {
		filteredRow := make(map[string]any)
		for _, field := range fields {
			if val, exists := row[field]; exists {
				filteredRow[field] = val
			}
		}
		result = append(result, filteredRow)
	}
	return result
}

// formatValue 格式化值为字符串
func (ec *ExportController) formatValue(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return ""
	default:
		return ""
	}
}

// insertData 插入数据到表
func (ec *ExportController) insertData(table string, data map[string]any) error {
	// 构建INSERT语句
	var columns []string
	var placeholders []string
	var values []any

	for col, val := range data {
		// 验证列名安全性
		if err := validateColumnName(col); err != nil {
			continue // 跳过无效列名
		}
		columns = append(columns, col)
		placeholders = append(placeholders, "?")
		values = append(values, val)
	}

	if len(columns) == 0 {
		return nil
	}

	query := "INSERT INTO " + table + " (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ")"

	_, err := ec.DB.SQL.Exec(query, values...)
	return err
}

// ExportTemplate 导出模板下载
func (ec *ExportController) ExportTemplate(c *ghttp.Context) {
	table := c.Query("table", "")
	if table == "" {
		c.Error(400, "请指定表名")
		return
	}

	// 从数据库查询表结构
	columns, err := ec.getTableColumns(table)
	if err != nil {
		c.Error(500, "获取表结构失败")
		return
	}

	// 生成CSV模板
	buf := new(bytes.Buffer)
	csvWriter := csv.NewWriter(buf)
	csvWriter.Write(columns)
	csvWriter.Flush()

	c.SetHeader("Content-Type", "text/csv; charset=utf-8")
	c.SetHeader("Content-Disposition", "attachment; filename="+table+"_template.csv")

	// 添加BOM
	bom := []byte{0xEF, 0xBB, 0xBF}
	c.Writer.Write(bom)
	c.Writer.Write(buf.Bytes())
}

// getTableColumns 获取表的列名
func (ec *ExportController) getTableColumns(table string) ([]string, error) {
	// SQLite 查询表结构
	query := "PRAGMA table_info(" + table + ")"
	rows, err := ec.DB.SQL.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name string
		var dataType string
		var notNull int
		var dfltValue *string
		var pk int

		if err := rows.Scan(&cid, &name, &dataType, &notNull, &dfltValue, &pk); err != nil {
			continue
		}
		columns = append(columns, name)
	}

	return columns, nil
}
