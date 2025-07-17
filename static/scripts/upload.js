/**
 * Tracks all images to be uploaded (local files and link-fetched images).
 * Key is either the file name (for local) or URL (for link).
 * @type {Map<string, { blob: Blob, rating: string, previewElement: HTMLElement, type: 'local' | 'link', nameOrUrl: string }>}
 */
const imageBlobMapping = new Map();

/**
 * Shows or hides the Upload All button in the .upload-area container based on imageBlobMapping size.
 * The button is created once and appended/removed as needed.
 * @function
 */
let uploadAllBtn = null;
function updateUploadAllBtn() {
    const uploadArea = document.querySelector('.upload-area');
    if (!uploadArea) return;
    if (!uploadAllBtn) {
        uploadAllBtn = document.createElement('button');
        uploadAllBtn.type = 'button';
        uploadAllBtn.textContent = 'Upload All';
        uploadAllBtn.className = 'upload-all-btn';
        /**
         * TODO: Implement actual upload logic here
         */
        uploadAllBtn.onclick = function () {
            alert('Upload All clicked!');
        };
    }
    if (imageBlobMapping.size > 0) {
        if (!uploadArea.contains(uploadAllBtn)) {
            uploadArea.appendChild(uploadAllBtn);
        }
    } else {
        if (uploadArea.contains(uploadAllBtn)) {
            uploadArea.removeChild(uploadAllBtn);
        }
    }
}

/**
 * Creates a preview element for an image (local or link).
 * @param {string} key - The key for the image (filename or URL)
 * @param {Blob} blob - The image blob
 * @param {'local'|'link'} type - Type of image
 * @param {string} nameOrUrl - Filename (local) or URL (link)
 * @returns {HTMLDivElement}
 */
function createPreviewElement(key, blob, type, nameOrUrl) {
    const previewElement = document.createElement('div');
    previewElement.className = 'preview-area';

    const previewImage = document.createElement('img');
    previewImage.className = 'preview-image';
    previewImage.src = URL.createObjectURL(blob);
    previewElement.appendChild(previewImage);

    const previewDetailsArea = document.createElement('div');
    previewDetailsArea.className = 'preview-details';
    previewElement.appendChild(previewDetailsArea);

    if (type === 'link') {
        const previewLink = document.createElement('a');
        previewLink.className = 'preview-link';
        previewLink.target = '_blank';
        previewLink.href = nameOrUrl;
        previewLink.textContent = nameOrUrl;
        previewDetailsArea.appendChild(previewLink);
    } else {
        const previewFile = document.createElement('span');
        previewFile.className = 'preview-link';
        previewFile.textContent = nameOrUrl;
        previewDetailsArea.appendChild(previewFile);
    }

    const previewRatingForm = document.createElement('form');
    previewRatingForm.className = 'preview-rating-form';
    ['Safe', 'Questionable', 'Sensitive', 'Explicit'].forEach((rating, idx) => {
        const inputId = `rating-${rating.toLowerCase()}-${key}`;
        const input = document.createElement('input');
        input.type = 'radio';
        input.name = `rating-${key}`;
        input.value = rating.toLowerCase();
        input.checked = idx === 0;
        input.id = inputId;
        const label = document.createElement('label');
        label.textContent = rating;
        label.setAttribute('for', inputId);
        previewRatingForm.appendChild(input);
        previewRatingForm.appendChild(label);
    });
    previewDetailsArea.appendChild(previewRatingForm);

    const removeBtn = document.createElement('button');
    removeBtn.type = 'button';
    removeBtn.textContent = 'Remove';
    removeBtn.className = 'preview-remove-btn';
    removeBtn.onclick = () => {
        previewElement.remove();
        imageBlobMapping.delete(key);
        updateUploadAllBtn();
    };
    previewDetailsArea.appendChild(removeBtn);

    return previewElement;
}

window.onbeforeunload = function (_event) {
    if (imageBlobMapping.size > 0) {
        return 'Are you sure you want to leave? Data you have entered may not be saved.';
    }
};

/**
 * Shows or hides a retro loading indicator in the upload area.
 * @param {boolean} show
 */
function setUploadLoadingIndicator(show) {
    let indicator = document.getElementById('_uploadLoadingIndicator');
    const dragBox = document.querySelector('.upload-drag-box');
    if (!dragBox) return;
    dragBox.style.position = 'relative';
    if (!indicator) {
        indicator = document.createElement('div');
        indicator.id = '_uploadLoadingIndicator';
        indicator.className = 'upload-loading-indicator';
        // Segmented ring spinner markup, no text
        indicator.innerHTML = `
            <div class="upload-loading-icon">
                <div class="ib-loader-spinner">
                    <div class="ib-loader-seg"></div>
                    <div class="ib-loader-seg"></div>
                    <div class="ib-loader-seg"></div>
                    <div class="ib-loader-seg"></div>
                    <div class="ib-loader-seg"></div>
                    <div class="ib-loader-seg"></div>
                    <div class="ib-loader-seg"></div>
                    <div class="ib-loader-seg"></div>
                    <div class="ib-loader-seg"></div>
                    <div class="ib-loader-seg"></div>
                    <div class="ib-loader-seg"></div>
                    <div class="ib-loader-seg"></div>
                </div>
            </div>
        `;
        dragBox.appendChild(indicator);
    }
    // Animation logic for retro ring spinner
    const spinner = indicator.querySelector('.ib-loader-spinner');
    const segments = spinner ? spinner.querySelectorAll('.ib-loader-seg') : [];
    if (show && segments.length) {
        let frame = 0;
        if (!indicator._retroAnimInterval) {
            indicator._retroAnimInterval = setInterval(() => {
                segments.forEach((seg, idx) => seg.classList.toggle('active', idx === frame));
                frame = (frame + 1) % segments.length;
            }, 80);
        }
    } else if (!show && indicator._retroAnimInterval) {
        clearInterval(indicator._retroAnimInterval);
        indicator._retroAnimInterval = null;
        segments.forEach(seg => seg.classList.remove('active'));
    }
    indicator.style.display = show ? 'flex' : 'none';
}

/**
 * Handles uploading an image via link using the backend proxy.
 * Optimized and uses async/await.
 * @returns {Promise<void>}
 */
async function uploadViaLink() {
    const uploadViaLinkInputBox = document.getElementById('_uploadViaLink_InputBox');
    const uploadViaLinkUploadPreviewsArea = document.getElementById('_uploadViaLink_UploadPreviewsArea');
    const link = uploadViaLinkInputBox.value.trim();
    hideError();

    if (!link) {
        showError('Please enter a valid image URL.');
        return;
    }
    if (imageBlobMapping.has(link)) {
        showError('This image has already been added.');
        return;
    }
    setUploadLoadingIndicator(true);
    try {
        const proxyUrl = `/posts/new/ilinkfetch?url=${encodeURIComponent(link)}`;
        const response = await fetch(proxyUrl);
        if (!response.ok) {
            let errorMsg = 'Failed to fetch the image from the provided URL.';
            try {
                const text = await response.text();
                if (text && text !== 'Failed to fetch image') errorMsg = text;
            } catch {
                errorMsg = 'An error occurred while fetching the image.';
            }
            showError(errorMsg || 'An error occurred while fetching the image.');
            return;
        }
        const contentType = response.headers.get('content-type') || '';
        if (!contentType.startsWith('image/')) {
            showError('The URL does not point to a valid image.');
            return;
        }
        const blob = await response.blob();
        if (!blob) {
            showError('No image data received from the URL.');
            return;
        }
        const previewElement = createPreviewElement(link, blob, 'link', link);
        uploadViaLinkUploadPreviewsArea.appendChild(previewElement);
        imageBlobMapping.set(link, { blob, rating: 'safe', previewElement, type: 'link', nameOrUrl: link });
        updateUploadAllBtn();
    } catch (error) {
        console.error('Error fetching image:', error);
        showError('An error occurred while fetching the image.');
    } finally {
        setUploadLoadingIndicator(false);
        uploadViaLinkInputBox.value = '';
    }
}

/**
 * Handles drag-and-drop and click-to-select for local image files.
 */
function setupLocalImageUpload() {
    const dragBox = document.querySelector('.upload-drag-box');
    const previewsArea = document.getElementById('_uploadViaLink_UploadPreviewsArea');
    const dragHeading = dragBox ? dragBox.querySelector('h1') : null;
    if (!dragBox || !previewsArea || !dragHeading) return;

    dragBox.addEventListener('dragover', function (e) {
        e.preventDefault();
        dragBox.classList.add('dragover');
        dragHeading.textContent = 'Release to upload!';
    });
    dragBox.addEventListener('dragleave', function (e) {
        e.preventDefault();
        dragBox.classList.remove('dragover');
        dragHeading.textContent = 'Drop files here or just click this box!';
    });
    dragBox.addEventListener('drop', async function (e) {
        e.preventDefault();
        dragBox.classList.remove('dragover');
        dragHeading.textContent = 'Drop files here or just click this box!';
        const files = e.dataTransfer.files;
        if (files && files.length > 0) {
            handleFiles(files);
        } else {
            // Try to get a URL from the drop
            let url = '';
            if (e.dataTransfer.items) {
                for (let i = 0; i < e.dataTransfer.items.length; i++) {
                    const item = e.dataTransfer.items[i];
                    if (item.kind === 'string' && (item.type === 'text/uri-list' || item.type === 'text/plain')) {
                        item.getAsString(function (s) {
                            if (s && s.match(/^https?:\/\/.+/)) {
                                handleDroppedUrl(s);
                            }
                        });
                        return;
                    }
                }
            }
            // Fallback for some browsers
            url = e.dataTransfer.getData('text/uri-list') || e.dataTransfer.getData('text/plain');
            if (url && url.match(/^https?:\/\/.+/)) {
                handleDroppedUrl(url);
            }
        }
    });
    dragBox.addEventListener('click', function () {
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = 'image/*';
        input.multiple = true;
        input.onchange = function () {
            handleFiles(input.files);
        };
        input.click();
    });
}

/**
 * Handles dropped URLs, tries to upload if it's an image or fetch if not.
 * @param {string} url
 */
function handleDroppedUrl(url) {
    // Accept direct image URLs or blob/data URLs
    const imageExt = /\.(jpg|jpeg|png|gif|webp|bmp|svg)$/i;
    if (imageExt.test(url) || url.startsWith('blob:') || url.startsWith('data:image/')) {
        uploadViaLinkDirect(url);
    } else {
        // Try to fetch and see if it's an image
        fetch(url, { method: 'HEAD' })
            .then(resp => {
                const type = resp.headers.get('content-type') || '';
                if (type.startsWith('image/')) {
                    uploadViaLinkDirect(url);
                } else {
                    showError('Dropped URL is not a direct image. Try dragging the image itself, not the page.');
                }
            })
            .catch(() => {
                showError('Could not fetch dropped URL.');
            });
    }
}

/**
 * Directly uploads an image via link (used for drag-drop URLs)
 * @param {string} url
 */
async function uploadViaLinkDirect(url) {
    setUploadLoadingIndicator(true);
    const uploadViaLinkUploadPreviewsArea = document.getElementById('_uploadViaLink_UploadPreviewsArea');
    if (imageBlobMapping.has(url)) {
        showError('This image has already been added.');
        return;
    }
    try {
        const proxyUrl = `/posts/new/ilinkfetch?url=${encodeURIComponent(url)}`;
        const response = await fetch(proxyUrl);
        if (!response.ok) {
            let errorMsg = 'Failed to fetch the image from the provided URL.';
            try {
                const text = await response.text();
                if (text && text !== 'Failed to fetch image') errorMsg = text;
            } catch {
                errorMsg = 'An error occurred while fetching the image.';
            }
            showError(errorMsg || 'An error occurred while fetching the image.');
            return;
        }
        const contentType = response.headers.get('content-type') || '';
        if (!contentType.startsWith('image/')) {
            showError('The URL does not point to a valid image.');
            return;
        }
        const blob = await response.blob();
        if (!blob) {
            showError('No image data received from the URL.');
            return;
        }
        const previewElement = createPreviewElement(url, blob, 'link', url);
        uploadViaLinkUploadPreviewsArea.appendChild(previewElement);
        imageBlobMapping.set(url, { blob, rating: 'safe', previewElement, type: 'link', nameOrUrl: url });
        updateUploadAllBtn();
    } catch (error) {
        console.error('Error fetching image:', error);
        showError('An error occurred while fetching the image.');
    } finally {
        setUploadLoadingIndicator(false);
    }
}

/**
 * Handles adding local files to the preview and mapping.
 * @param {FileList} files
 */
function handleFiles(files) {
    const previewsArea = document.getElementById('_uploadViaLink_UploadPreviewsArea');
    for (const file of files) {
        if (!file.type.startsWith('image/')) continue;
        if (imageBlobMapping.has(file.name)) continue;
        const previewElement = createPreviewElement(file.name, file, 'local', file.name);
        previewsArea.appendChild(previewElement);
        imageBlobMapping.set(file.name, { blob: file, rating: 'safe', previewElement, type: 'local', nameOrUrl: file.name });
    }
    updateUploadAllBtn();
}

setupLocalImageUpload();
