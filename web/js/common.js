// 公共工具函数
const API_BASE = '/api/v1';

// ====== 深色模式 ======
function initTheme() {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
        document.documentElement.classList.add('dark-mode');
    }
}

function toggleTheme() {
    const isDark = document.documentElement.classList.toggle('dark-mode');
    localStorage.setItem('theme', isDark ? 'dark' : 'light');

    // 更新按钮图标
    const btn = document.getElementById('theme-toggle');
    if (btn) {
        btn.innerHTML = isDark ? '<i class="layui-icon layui-icon-light"></i>' : '<i class="layui-icon layui-icon-moon"></i>';
    }
}

// 页面加载时初始化主题
initTheme();

// API 请求封装
async function request(url, options = {}) {
    const token = localStorage.getItem('token');
    const headers = {
        'Content-Type': 'application/json',
        ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
        ...options.headers
    };

    try {
        const response = await fetch(API_BASE + url, {
            ...options,
            headers
        });

        if (response.status === 401) {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = '/login';
            return;
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Request error:', error);
        return { code: -1, message: '网络错误' };
    }
}

// 加载当前用户信息
async function loadCurrentUser() {
    try {
        const res = await request('/auth/me');
        if (res.code === 0 && res.data) {
            localStorage.setItem('user', JSON.stringify(res.data));
            const usernameEl = document.getElementById('username');
            if (usernameEl) {
                usernameEl.textContent = res.data.username || 'Admin';
            }
            return res.data;
        }
    } catch (e) {
        console.error('Load user error:', e);
    }
    return null;
}

// 检查 Token 是否过期
function isTokenExpired() {
    const token = localStorage.getItem('token');
    if (!token) return true;
    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        if (payload.exp) {
            return payload.exp * 1000 < Date.now();
        }
        return false;
    } catch (e) {
        return true;
    }
}

// 检查登录状态
function checkAuth() {
    if (!localStorage.getItem('token') || isTokenExpired()) {
        window.location.href = '/login';
        return false;
    }
    return true;
}

// HTML 转义
function esc(str) {
    if (str === null || str === undefined) return '';
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

// 格式化 JSON
function formatJson(value) {
    if (!value) return '';
    if (typeof value === 'string') {
        try {
            const parsed = JSON.parse(value);
            return JSON.stringify(parsed, null, 2);
        } catch (e) {
            return value;
        }
    }
    return JSON.stringify(value, null, 2);
}

// 格式化文件大小
function formatFileSize(size) {
    if (!size) return '0 B';
    if (size < 1024) return size + ' B';
    if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB';
    if (size < 1024 * 1024 * 1024) return (size / 1024 / 1024).toFixed(1) + ' MB';
    return (size / 1024 / 1024 / 1024).toFixed(1) + ' GB';
}

// 格式化日期
function formatDate(dateStr) {
    if (!dateStr) return '-';
    const date = new Date(dateStr);
    return date.toLocaleString('zh-CN');
}

// 获取 URL 参数
function getQueryParam(name) {
    const params = new URLSearchParams(window.location.search);
    return params.get(name);
}

// 显示加载中
function showLoading() {
    return layer.load(1, { shade: [0.1, '#fff'] });
}

// 关闭加载
function hideLoading(index) {
    layer.close(index);
}

// 成功提示
function showSuccess(msg) {
    layer.msg(msg, { icon: 1 });
}

// 错误提示
function showError(msg) {
    layer.msg(msg, { icon: 2 });
}

// 侧边栏菜单模板
const sidebarTemplates = {
    dashboard: `
        <ul class="sidebar-menu">
            <li class="sidebar-menu-item" data-page="index">
                <a href="/pages/dashboard/index.html"><i class="layui-icon layui-icon-chart"></i> 首页概览</a>
            </li>
            <li class="sidebar-menu-item" data-page="notifications">
                <a href="/pages/dashboard/notifications.html"><i class="layui-icon layui-icon-notice"></i> 我的通知</a>
            </li>
        </ul>
    `,
    system: `
        <ul class="sidebar-menu">
            <li class="sidebar-menu-item" data-page="users">
                <a href="/pages/system/users.html"><i class="layui-icon layui-icon-username"></i> 用户管理</a>
            </li>
            <li class="sidebar-menu-item" data-page="roles">
                <a href="/pages/system/roles.html"><i class="layui-icon layui-icon-group"></i> 角色管理</a>
            </li>
            <li class="sidebar-menu-item" data-page="menus">
                <a href="/pages/system/menus.html"><i class="layui-icon layui-icon-menu-fill"></i> 菜单管理</a>
            </li>
            <li class="sidebar-menu-item" data-page="permissions">
                <a href="/pages/system/permissions.html"><i class="layui-icon layui-icon-auz"></i> 权限管理</a>
            </li>
            <li class="sidebar-menu-item" data-page="uploads">
                <a href="/pages/system/uploads.html"><i class="layui-icon layui-icon-picture"></i> 文件管理</a>
            </li>
            <li class="sidebar-menu-item" data-page="tasks">
                <a href="/pages/system/tasks.html"><i class="layui-icon layui-icon-time"></i> 定时任务</a>
            </li>
            <li class="sidebar-menu-item" data-page="settings">
                <a href="/pages/system/settings.html"><i class="layui-icon layui-icon-set"></i> 系统设置</a>
            </li>
        </ul>
    `,
    data: `
        <ul class="sidebar-menu">
            <li class="sidebar-menu-item" data-page="export">
                <a href="/pages/data/export.html"><i class="layui-icon layui-icon-export"></i> 数据导出</a>
            </li>
            <li class="sidebar-menu-item" data-page="import">
                <a href="/pages/data/import.html"><i class="layui-icon layui-icon-upload-drag"></i> 数据导入</a>
            </li>
        </ul>
    `,
    logs: `
        <ul class="sidebar-menu">
            <li class="sidebar-menu-item" data-page="operation">
                <a href="/pages/logs/operation.html"><i class="layui-icon layui-icon-log"></i> 操作日志</a>
            </li>
            <li class="sidebar-menu-item" data-page="audit">
                <a href="/pages/logs/audit.html"><i class="layui-icon layui-icon-read"></i> 审计日志</a>
            </li>
        </ul>
    `
};

// 渲染侧边栏
function renderSidebar(module, currentPage) {
    const container = document.querySelector('.layui-side-scroll');
    if (container && sidebarTemplates[module]) {
        container.innerHTML = sidebarTemplates[module];
        // 设置当前页高亮
        const items = container.querySelectorAll('.sidebar-menu-item');
        items.forEach(item => {
            if (item.dataset.page === currentPage) {
                item.classList.add('active');
            }
        });
    }
}
