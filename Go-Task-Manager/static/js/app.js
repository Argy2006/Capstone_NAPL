/**
 * TaskFlow – Personal Task Manager
 * Frontend Controller & API Client
 */

document.addEventListener('DOMContentLoaded', () => {
    console.log('TaskFlow initialized. Connected to Go backend.');

    const taskListEl = document.getElementById('taskList');
    const searchInput = document.getElementById('searchInput');
    const newTaskBtn = document.getElementById('newTaskBtn');
    const refreshBtn = document.getElementById('refreshBtn');

    // Fetch initial task state from Go backend API
    async function loadTasks() {
        try {
            const response = await fetch('/api/tasks');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const tasks = await response.json();
            console.log('Loaded tasks from server:', tasks);
            renderTasks(tasks);
        } catch (error) {
            console.warn('Could not retrieve tasks from /api/tasks:', error);
        }
    }

    // Render task cards or empty state
    function renderTasks(tasks) {
        if (!taskListEl) return;

        if (!tasks || tasks.length === 0) {
            taskListEl.innerHTML = `
                <div class="empty-state">
                    <h4>No Tasks Found</h4>
                    <p>Your task list is currently empty. Direct Google Stitch UI import or task creation will populate this section.</p>
                </div>
            `;
            return;
        }

        taskListEl.innerHTML = tasks.map(task => `
            <div class="task-card" data-id="${escapeHtml(task.id || '')}">
                <div class="task-info">
                    <h4>${escapeHtml(task.title || 'Untitled Task')}</h4>
                    <p>${escapeHtml(task.description || 'No description provided.')}</p>
                    <div class="task-meta">
                        <span class="tag tag-priority-${(task.priority || 'medium').toLowerCase()}">${escapeHtml(task.priority || 'Medium')}</span>
                        <span class="tag tag-status">${escapeHtml(task.status || 'Pending')}</span>
                        ${task.category ? `<span class="tag tag-category">${escapeHtml(task.category)}</span>` : ''}
                        ${task.dueDate ? `<span class="tag tag-due">📅 ${escapeHtml(task.dueDate)}</span>` : ''}
                    </div>
                </div>
            </div>
        `).join('');
    }

    // Utility to prevent XSS
    function escapeHtml(str) {
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }

    // Event Listeners
    if (searchInput) {
        searchInput.addEventListener('input', (e) => {
            const query = e.target.value.toLowerCase();
            const cards = document.querySelectorAll('.task-card');
            cards.forEach(card => {
                const text = card.textContent.toLowerCase();
                card.style.display = text.includes(query) ? '' : 'none';
            });
        });
    }

    if (newTaskBtn) {
        newTaskBtn.addEventListener('click', () => {
            alert('Task creation endpoint prepared. Full CRUD functionality will be implemented following Phase II confirmation.');
        });
    }

    if (refreshBtn) {
        refreshBtn.addEventListener('click', () => {
            loadTasks();
        });
    }

    // Initial load
    loadTasks();
});

