function showError(message) {
    const errorMessage = document.getElementById('error-message');
    errorMessage.textContent = message;
    errorMessage.style.display = 'block';
}

function hideError() {
    const errorMessage = document.getElementById('error-message');
    errorMessage.textContent = '';
    errorMessage.style.display = 'none';
}