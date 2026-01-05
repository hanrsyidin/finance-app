/**
 * Utils.js - Utility Functions
 * Personal Finance Dashboard
 */

/**
 * Format number as Indonesian Rupiah currency
 * @param {number} amount - Amount to format
 * @returns {string} Formatted currency string
 */
function formatCurrency(amount) {
    return new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0
    }).format(amount);
}

/**
 * Format date to Indonesian locale
 * @param {string|Date} date - Date to format
 * @param {object} options - Intl.DateTimeFormat options
 * @returns {string} Formatted date string
 */
function formatDate(date, options = {}) {
    const defaultOptions = {
        day: 'numeric',
        month: 'long',
        year: 'numeric'
    };
    return new Date(date).toLocaleDateString('id-ID', { ...defaultOptions, ...options });
}

/**
 * Format date to relative time (e.g., "2 hours ago")
 * @param {string|Date} date - Date to format
 * @returns {string} Relative time string
 */
function formatRelativeTime(date) {
    const now = new Date();
    const past = new Date(date);
    const diffMs = now - past;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);

    if (diffMins < 1) return 'Baru saja';
    if (diffMins < 60) return `${diffMins} menit lalu`;
    if (diffHours < 24) return `${diffHours} jam lalu`;
    if (diffDays < 7) return `${diffDays} hari lalu`;
    
    return formatDate(date);
}

/**
 * Parse currency string to number
 * @param {string} str - Currency string (e.g., "Rp 25.000")
 * @returns {number} Parsed number
 */
function parseCurrency(str) {
    return parseInt(str.replace(/\D/g, '')) || 0;
}

/**
 * Format number with thousand separator
 * @param {number} num - Number to format
 * @returns {string} Formatted number
 */
function formatNumber(num) {
    return new Intl.NumberFormat('id-ID').format(num);
}

/**
 * Get month name in Indonesian
 * @param {number} month - Month index (0-11)
 * @returns {string} Month name
 */
function getMonthName(month) {
    const months = [
        'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
        'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
    ];
    return months[month];
}

/**
 * Get current month in YYYY-MM format
 * @returns {string} Current month
 */
function getCurrentMonth() {
    return new Date().toISOString().slice(0, 7);
}

/**
 * Get today's date in YYYY-MM-DD format
 * @returns {string} Today's date
 */
function getToday() {
    return new Date().toISOString().slice(0, 10);
}

/**
 * Debounce function
 * @param {Function} func - Function to debounce
 * @param {number} wait - Wait time in milliseconds
 * @returns {Function} Debounced function
 */
function debounce(func, wait = 300) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

/**
 * Simple hash function for generating unique IDs
 * @returns {string} Random ID
 */
function generateId() {
    return Math.random().toString(36).substring(2, 9);
}

/**
 * Download data as CSV file
 * @param {Array} data - Array of objects
 * @param {string} filename - Output filename
 */
function downloadCSV(data, filename = 'export.csv') {
    if (!data.length) return;
    
    const headers = Object.keys(data[0]);
    const csvContent = [
        headers.join(','),
        ...data.map(row => 
            headers.map(header => {
                let value = row[header];
                if (typeof value === 'string' && value.includes(',')) {
                    value = `"${value}"`;
                }
                return value;
            }).join(',')
        )
    ].join('\n');
    
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = filename;
    link.click();
}

/**
 * Show browser notification
 * @param {string} title - Notification title
 * @param {string} body - Notification body
 */
async function showNotification(title, body) {
    if (!('Notification' in window)) return;
    
    if (Notification.permission === 'granted') {
        new Notification(title, { body });
    } else if (Notification.permission !== 'denied') {
        const permission = await Notification.requestPermission();
        if (permission === 'granted') {
            new Notification(title, { body });
        }
    }
}

/**
 * Copy text to clipboard
 * @param {string} text - Text to copy
 * @returns {Promise<boolean>} Success status
 */
async function copyToClipboard(text) {
    try {
        await navigator.clipboard.writeText(text);
        return true;
    } catch (err) {
        console.error('Failed to copy:', err);
        return false;
    }
}

// Export for module usage (if needed)
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        formatCurrency,
        formatDate,
        formatRelativeTime,
        parseCurrency,
        formatNumber,
        getMonthName,
        getCurrentMonth,
        getToday,
        debounce,
        generateId,
        downloadCSV,
        showNotification,
        copyToClipboard
    };
}
