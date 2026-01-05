/**
 * App.js - Alpine.js Main Application Component
 * Personal Finance Dashboard
 */

/**
 * Global app state and methods
 * This component is initialized on the body element
 */
function app() {
    return {
        // Toast notification state
        toast: {
            show: false,
            message: '',
            type: 'info' // 'success', 'error', 'info'
        },

        /**
         * Show toast notification
         * @param {string} message - Message to display
         * @param {string} type - Toast type (success, error, info)
         * @param {number} duration - Duration in milliseconds
         */
        showToast(message, type = 'info', duration = 3000) {
            this.toast = { show: true, message, type };

            setTimeout(() => {
                this.toast.show = false;
            }, duration);
        },

        /**
         * Close any open modal
         * Called when pressing Escape key
         */
        closeModal() {
            // This will be handled by child components
            // The event bubbles up from the modal
        },

        /**
         * Format currency helper (available globally)
         */
        formatCurrency(amount) {
            return new Intl.NumberFormat('id-ID', {
                style: 'currency',
                currency: 'IDR',
                minimumFractionDigits: 0
            }).format(amount);
        },

        /**
         * Format date helper
         */
        formatDate(date, options = {}) {
            const defaultOptions = {
                day: 'numeric',
                month: 'long',
                year: 'numeric'
            };
            return new Date(date).toLocaleDateString('id-ID', { ...defaultOptions, ...options });
        }
    };
}

// Make Alpine store for global toast notifications
document.addEventListener('alpine:init', () => {
    Alpine.store('toast', {
        show: false,
        message: '',
        type: 'info'
    });

    // Watch for store changes to update toast
    Alpine.effect(() => {
        const toast = Alpine.store('toast');
        if (toast.show) {
            setTimeout(() => {
                Alpine.store('toast', { ...toast, show: false });
            }, 3000);
        }
    });
});

// Keyboard shortcuts handler
document.addEventListener('DOMContentLoaded', () => {
    document.addEventListener('keydown', (e) => {
        // Don't trigger shortcuts when typing in inputs
        const activeEl = document.activeElement;
        if (activeEl.tagName === 'INPUT' || activeEl.tagName === 'TEXTAREA' || activeEl.tagName === 'SELECT') {
            return;
        }

        // Arrow keys for month navigation
        if (e.key === 'ArrowLeft' || e.key === 'ArrowRight') {
            // Month navigation will be handled by component
        }

        // Slash for search focus
        if (e.key === '/') {
            const searchInput = document.querySelector('[data-search]');
            if (searchInput) {
                e.preventDefault();
                searchInput.focus();
            }
        }
    });
});

// Service Worker registration (for PWA support - optional)
if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
        // Uncomment to enable PWA
        // navigator.serviceWorker.register('/sw.js');
    });
}

// Prevent form resubmission on page refresh
if (window.history.replaceState) {
    window.history.replaceState(null, null, window.location.href);
}
