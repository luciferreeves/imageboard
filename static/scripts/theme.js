document.addEventListener('DOMContentLoaded', function () {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);

    const preferencesForm = document.getElementById('preferences-form');
    if (preferencesForm) {
        const themeRadios = document.querySelectorAll('input[name="theme"]');
        themeRadios.forEach(radio => {
            if (radio.value === savedTheme) {
                radio.checked = true;
            }
        });

        preferencesForm.addEventListener('submit', function (e) {
            e.preventDefault();
            const selectedTheme = document.querySelector('input[name="theme"]:checked').value;
            localStorage.setItem('theme', selectedTheme);
            document.documentElement.setAttribute('data-theme', selectedTheme);

            let successMsg = document.querySelector('.success-message');
            if (successMsg) {
                successMsg.remove();
            }

            const message = document.createElement('div');
            message.className = 'success-message';
            message.textContent = 'Preferences saved successfully!';
            preferencesForm.parentNode.insertBefore(message, preferencesForm);

            setTimeout(() => {
                if (message.parentNode) {
                    message.remove();
                }
            }, 3000);
        });
    }
});