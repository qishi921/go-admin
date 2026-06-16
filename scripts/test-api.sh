#!/bin/bash

# Admin System API 测试脚本
# 用法: ./test-api.sh [base_url]

BASE_URL="${1:-http://localhost:8080}"
TOKEN=""

echo "=== Admin System API 测试 ==="
echo "基础 URL: $BASE_URL"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

pass() {
    echo -e "${GREEN}✓ $1${NC}"
}

fail() {
    echo -e "${RED}✗ $1${NC}"
    exit 1
}

# 1. 健康检查
echo "--- 健康检查 ---"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/health" 2>/dev/null)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
BODY=$(echo "$RESPONSE" | head -n -1)

if [ "$HTTP_CODE" = "200" ]; then
    pass "健康检查通过"
else
    fail "健康检查失败 (HTTP $HTTP_CODE)"
fi

# 2. 登录测试
echo ""
echo "--- 登录测试 ---"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}' 2>/dev/null)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
BODY=$(echo "$RESPONSE" | head -n -1)

if [ "$HTTP_CODE" = "200" ]; then
    TOKEN=$(echo "$BODY" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    if [ -n "$TOKEN" ]; then
        pass "登录成功，获取到 Token"
    else
        fail "登录响应缺少 Token"
    fi
else
    fail "登录失败 (HTTP $HTTP_CODE)"
fi

# 3. 用户列表测试
echo ""
echo "--- 用户列表测试 ---"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/users" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
BODY=$(echo "$RESPONSE" | head -n -1)

if [ "$HTTP_CODE" = "200" ]; then
    CODE=$(echo "$BODY" | grep -o '"code":[0-9]*' | cut -d':' -f2)
    if [ "$CODE" = "0" ]; then
        pass "获取用户列表成功"
    else
        fail "获取用户列表失败 (code: $CODE)"
    fi
else
    fail "获取用户列表失败 (HTTP $HTTP_CODE)"
fi

# 4. 角色列表测试
echo ""
echo "--- 角色列表测试 ---"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/roles" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

if [ "$HTTP_CODE" = "200" ]; then
    pass "获取角色列表成功"
else
    fail "获取角色列表失败 (HTTP $HTTP_CODE)"
fi

# 5. 菜单列表测试
echo ""
echo "--- 菜单列表测试 ---"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/menus" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

if [ "$HTTP_CODE" = "200" ]; then
    pass "获取菜单列表成功"
else
    fail "获取菜单列表失败 (HTTP $HTTP_CODE)"
fi

# 6. 权限列表测试
echo ""
echo "--- 权限列表测试 ---"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/permissions" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

if [ "$HTTP_CODE" = "200" ]; then
    pass "获取权限列表成功"
else
    fail "获取权限列表失败 (HTTP $HTTP_CODE)"
fi

# 7. 仪表盘统计测试
echo ""
echo "--- 仪表盘统计测试 ---"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/dashboard/stats" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

if [ "$HTTP_CODE" = "200" ]; then
    pass "获取仪表盘统计成功"
else
    fail "获取仪表盘统计失败 (HTTP $HTTP_CODE)"
fi

# 8. 系统设置测试
echo ""
echo "--- 系统设置测试 ---"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/settings" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

if [ "$HTTP_CODE" = "200" ]; then
    pass "获取系统设置成功"
else
    fail "获取系统设置失败 (HTTP $HTTP_CODE)"
fi

# 9. 通知列表测试
echo ""
echo "--- 通知列表测试 ---"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/notifications" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

if [ "$HTTP_CODE" = "200" ]; then
    pass "获取通知列表成功"
else
    fail "获取通知列表失败 (HTTP $HTTP_CODE)"
fi

# 10. 定时任务测试
echo ""
echo "--- 定时任务测试 ---"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/v1/scheduled-tasks" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

if [ "$HTTP_CODE" = "200" ]; then
    pass "获取定时任务列表成功"
else
    fail "获取定时任务列表失败 (HTTP $HTTP_CODE)"
fi

echo ""
echo "=== 测试完成 ==="
