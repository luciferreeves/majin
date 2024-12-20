function initializeEmailInterface() {
    const emailList = document.getElementById('email-list');
    const resizer = document.getElementById('resizer');
    const contextMenu = document.getElementById('context-menu');

    // Cleanup old event listeners if they exist
    if (resizer.resizeInitialized) {
        return; // Already initialized
    }

    // Resize functionality
    let isResizing = false;
    let startX = 0;
    let startWidth = 0;
    const minWidth = 384; // w-96
    const maxWidth = 800;

    function startResizing(e) {
        isResizing = true;
        startX = e.clientX;
        startWidth = emailList.offsetWidth;
        document.addEventListener('mousemove', handleResize);
        document.addEventListener('mouseup', stopResizing);
    }

    function handleResize(e) {
        if (!isResizing) return;

        const diff = e.clientX - startX;
        let newWidth = startWidth + diff;

        // Enforce min/max constraints
        newWidth = Math.max(minWidth, Math.min(maxWidth, newWidth));
        emailList.style.width = `${newWidth}px`;
    }

    function stopResizing() {
        isResizing = false;
        document.removeEventListener('mousemove', handleResize);
        document.removeEventListener('mouseup', stopResizing);
    }

    // Context menu functionality
    function showContextMenu(e) {
        e.preventDefault();
        contextMenu.style.display = 'block';

        // Adjust menu position to stay within viewport
        const x = Math.min(e.pageX, window.innerWidth - contextMenu.offsetWidth);
        const y = Math.min(e.pageY, window.innerHeight - contextMenu.offsetHeight);

        contextMenu.style.left = `${x}px`;
        contextMenu.style.top = `${y}px`;
    }

    function hideContextMenu() {
        if (contextMenu) {
            contextMenu.style.display = 'none';
        }
    }

    // Event listeners
    resizer.addEventListener('mousedown', startResizing);
    resizer.resizeInitialized = true; // Mark as initialized

    document.querySelectorAll('.email-item').forEach(item => {
        item.addEventListener('contextmenu', showContextMenu);
    });

    document.addEventListener('click', hideContextMenu);
    document.addEventListener('contextmenu', function (e) {
        if (!e.target.closest('.email-item')) {
            e.preventDefault();
            hideContextMenu();
        }
    });
}

// Initialize on regular page load
document.addEventListener('DOMContentLoaded', initializeEmailInterface);

// Initialize after HTMX content loads
document.addEventListener('htmx:afterSettle', initializeEmailInterface);