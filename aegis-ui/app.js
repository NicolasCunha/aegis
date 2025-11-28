// API Base URL - adjust based on environment
const API_BASE = '/api/aegis';

// State management
const state = {
    users: [],
    roles: [],
    permissions: [],
    healthStatus: 'unknown'
};

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    initializeTabs();
    initializeRefreshButtons();
    checkHealth();
    loadAllData();
    
    // Auto-refresh health every 30 seconds
    setInterval(checkHealth, 30000);
});

// Tab navigation
function initializeTabs() {
    const tabButtons = document.querySelectorAll('.tab-button');
    
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const targetTab = button.dataset.tab;
            
            // Update active button
            tabButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');
            
            // Update active content
            document.querySelectorAll('.tab-content').forEach(content => {
                content.classList.remove('active');
            });
            document.getElementById(`${targetTab}Tab`).classList.add('active');
        });
    });
}

// Refresh buttons
function initializeRefreshButtons() {
    document.getElementById('refreshUsers').addEventListener('click', loadUsers);
    document.getElementById('refreshRoles').addEventListener('click', loadRoles);
    document.getElementById('refreshPermissions').addEventListener('click', loadPermissions);
}

// API calls
async function apiCall(endpoint, options = {}) {
    try {
        const response = await fetch(`${API_BASE}${endpoint}`, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            }
        });
        
        if (!response.ok) {
            throw new Error(`API Error: ${response.status}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('API call failed:', error);
        throw error;
    }
}

// Health check
async function checkHealth() {
    try {
        const response = await fetch(`${API_BASE}/health`);
        const data = await response.json();
        
        state.healthStatus = data.status === 'healthy' ? 'healthy' : 'unhealthy';
        updateHealthStatus();
    } catch (error) {
        state.healthStatus = 'unhealthy';
        updateHealthStatus();
    }
}

function updateHealthStatus() {
    const statusEl = document.getElementById('healthStatus');
    const indicator = statusEl.querySelector('.status-indicator');
    const text = statusEl.querySelector('.status-text');
    
    indicator.className = `status-indicator ${state.healthStatus}`;
    text.textContent = state.healthStatus === 'healthy' ? 'Service Online' : 'Service Offline';
}

// Load all data
async function loadAllData() {
    await Promise.all([
        loadUsers(),
        loadRoles(),
        loadPermissions()
    ]);
}

// Load users
async function loadUsers() {
    const container = document.getElementById('usersContainer');
    container.innerHTML = '<div class="loading">Loading users...</div>';
    
    try {
        const data = await apiCall('/users');
        state.users = data;
        renderUsers();
    } catch (error) {
        container.innerHTML = '<div class="error">Failed to load users. Please try again.</div>';
    }
}

function renderUsers() {
    const container = document.getElementById('usersContainer');
    document.getElementById('usersCount').textContent = "(" + state.users.length + ")";
    
    if (state.users.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">👤</div>
                <p>No users found</p>
            </div>
        `;
        return;
    }
    
    container.innerHTML = `
        <table class="data-table">
            <thead>
                <tr>
                    <th>Subject</th>
                    <th>Roles</th>
                    <th>Permissions</th>
                </tr>
            </thead>
            <tbody>
                ${state.users.map(user => `
                    <tr data-user-id="${user.id}" style="cursor: pointer;">
                        <td>
                            <div class="table-subject">
                                <span class="table-subject-email">${escapeHtml(user.subject)}</span>
                                <span class="table-subject-id">${user.id}</span>
                            </div>
                        </td>
                        <td>
                            <div class="table-badges">
                                ${(user.roles || []).map(role => 
                                    `<span class="badge badge-role">🎭 ${escapeHtml(role)}</span>`
                                ).join('')}
                                ${(user.roles || []).length === 0 ? '<span style="color: var(--text-muted)">None</span>' : ''}
                            </div>
                        </td>
                        <td>
                            <div class="table-badges">
                                ${(user.permissions || []).map(permission => 
                                    `<span class="badge badge-permission">🔐 ${escapeHtml(permission)}</span>`
                                ).join('')}
                                ${(user.permissions || []).length === 0 ? '<span style="color: var(--text-muted)">None</span>' : ''}
                            </div>
                        </td>
                    </tr>
                `).join('')}
            </tbody>
        </table>
    `;
}

// Load roles
async function loadRoles() {
    const container = document.getElementById('rolesContainer');
    container.innerHTML = '<div class="loading">Loading roles...</div>';
    
    try {
        const data = await apiCall('/roles');
        state.roles = data;
        renderRoles();
    } catch (error) {
        container.innerHTML = '<div class="error">Failed to load roles. Please try again.</div>';
    }
}

function renderRoles() {
    const container = document.getElementById('rolesContainer');
    document.getElementById('rolesCount').textContent = "(" + state.roles.length + ")";
    
    if (state.roles.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">🎭</div>
                <p>No roles found</p>
            </div>
        `;
        return;
    }
    
    container.innerHTML = `
        <table class="data-table">
            <thead>
                <tr>
                    <th>Role Name</th>
                    <th>Description</th>
                    <th>Created</th>
                    <th>Updated</th>
                </tr>
            </thead>
            <tbody>
                ${state.roles.map(role => `
                    <tr>
                        <td>
                            <div class="table-subject">
                                <span class="table-subject-email">🎭 ${escapeHtml(role.name)}</span>
                            </div>
                        </td>
                        <td>${escapeHtml(role.description || 'No description')}</td>
                        <td>${formatDate(role.created_at)}</td>
                        <td>${formatDate(role.updated_at)}</td>
                    </tr>
                `).join('')}
            </tbody>
        </table>
    `;
}

// Load permissions
async function loadPermissions() {
    const container = document.getElementById('permissionsContainer');
    container.innerHTML = '<div class="loading">Loading permissions...</div>';
    
    try {
        const data = await apiCall('/permissions');
        state.permissions = data;
        renderPermissions();
    } catch (error) {
        container.innerHTML = '<div class="error">Failed to load permissions. Please try again.</div>';
    }
}

function renderPermissions() {
    const container = document.getElementById('permissionsContainer');
    document.getElementById('permissionsCount').textContent = "(" + state.permissions.length + ")";
    if (state.permissions.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">🔐</div>
                <p>No permissions found</p>
            </div>
        `;
        return;
    }
    
    container.innerHTML = `
        <table class="data-table">
            <thead>
                <tr>
                    <th>Permission Name</th>
                    <th>Description</th>
                    <th>Created</th>
                    <th>Updated</th>
                </tr>
            </thead>
            <tbody>
                ${state.permissions.map(permission => `
                    <tr>
                        <td>
                            <div class="table-subject">
                                <span class="table-subject-email">🔐 ${escapeHtml(permission.name)}</span>
                            </div>
                        </td>
                        <td>${escapeHtml(permission.description || 'No description')}</td>
                        <td>${formatDate(permission.created_at)}</td>
                        <td>${formatDate(permission.updated_at)}</td>
                    </tr>
                `).join('')}
            </tbody>
        </table>
    `;
}

// Utility functions
function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    return date.toLocaleString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}
