function validateCSSFieldValue(value) {
    const validCSSValuePattern = /^(auto|(\d+(\.\d+)?(px|em|rem|%)?)|calc\(.+\))$/;
    return validCSSValuePattern.test(value);
}

function validateCSSFontSize(value) {
    const validFontSizePattern = /^(small|medium|large|x-large|\d+(\.\d+)?(px|em|rem|%)?)$/;
    return validFontSizePattern.test(value);
}

function setPreferences() {
    const preferences = {
        sidebar_width: document.getElementById('sidebar-width').value,
        main_content_width: document.getElementById('main-width').value,
        h1_font_size: document.getElementById('h1-font-size').value,
        body_font_size: document.getElementById('body-font-size').value,
        small_font_size: document.getElementById('small-font-size').value,
        posts_per_page: parseInt(document.getElementById('posts-per-page').value, 10)
    }

    for (const key in preferences) {
        if (preferences[key] === '') {
            showError(`Please fill in the ${key.replace(/_/g, ' ')} field.`);
            return;
        }

        switch (key) {
            case 'sidebar_width':
            case 'main_content_width':
                if (!validateCSSFieldValue(preferences[key])) {
                    showError(`Invalid value for ${key.replace(/_/g, ' ')}: ${preferences[key]}. Please enter a valid CSS value (e.g., '300px', '50%', 'auto').`);
                    return;
                }
                break;
            case 'h1_font_size':
            case 'body_font_size':
            case 'small_font_size':
                if (!validateCSSFontSize(preferences[key])) {
                    showError(`Invalid font size for ${key.replace(/_/g, ' ')}: ${preferences[key]}. Please enter a valid font size (e.g., 'small', 'medium', 'large', '16px', '1.2em').`);
                    return;
                }
                break;
            case 'posts_per_page':
                if (isNaN(preferences[key]) || preferences[key] <= 0) {
                    showError('Posts per page must be a positive integer.');
                    return;
                }
                break;
        }
    }

    document.cookie = `preferences=${JSON.stringify(preferences)}; path=/; SameSite=Lax;`
    window.location.reload()
}

function resetPreferences() {
    hideError();
    document.cookie = 'preferences=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; SameSite=Lax;'
    window.location.reload()
}
