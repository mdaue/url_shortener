// static/js/app.js
document.addEventListener('DOMContentLoaded', function() {
    // Listen for the custom event triggered by HTMX
    document.body.addEventListener('urlAdded', function() {
        // Show a success message
        const message = document.createElement('div');
        message.textContent = 'URL shortened successfully!';
        message.style.cssText = 'position: fixed; top: 20px; right: 20px; background: #4CAF50; color: white; padding: 10px; border-radius: 4px; box-shadow: 0 2px 4px rgba(0,0,0,0.2);';
        document.body.appendChild(message);
        
        // Remove the message after 3 seconds
        setTimeout(function() {
            document.body.removeChild(message);
        }, 3000);
    });
});