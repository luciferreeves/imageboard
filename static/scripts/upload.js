/**
 * Tracks all images to be uploaded (local files and link-fetched images).
 * Key is either the file name (for local) or URL (for link).
 * @type {Map<string, { blob: Blob, rating: string, previewElement: HTMLElement, type: 'local' | 'link', nameOrUrl: string }>}
 */
const imageBlobMapping = new Map();

/**
 * Tracks whether an upload is currently in progress
 * @type {boolean}
 */
let isUploading = false;

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
        uploadAllBtn.onclick = async function () {
            await uploadAllImages();
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
        if (isUploading) return; // Prevent removal during upload
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

/**
 * Uploads all images in the imageBlobMapping to the server.
 * Shows progress and handles success/error states.
 */
async function uploadAllImages() {
    if (imageBlobMapping.size === 0) return;

    const totalImages = imageBlobMapping.size;
    let uploadedCount = 0;
    let hasErrors = false;

    // Set upload state and disable remove buttons
    isUploading = true;
    disableAllRemoveButtons();

    // Disable the upload button and show progress
    uploadAllBtn.disabled = true;
    updateUploadButtonText(uploadedCount, totalImages, false);

    try {
        for (const [key, imageData] of imageBlobMapping) {
            try {
                // Get the selected rating for this image
                const selectedRating = getSelectedRating(key);
                imageData.rating = selectedRating;

                // Create and submit form for this image
                await uploadSingleImage(imageData, selectedRating);

                // Mark this image as successfully uploaded
                markImageAsUploaded(imageData.previewElement);
                uploadedCount++;
                updateUploadButtonText(uploadedCount, totalImages, false);

            } catch (error) {
                console.error(`Failed to upload image ${key}:`, error);
                markImageAsError(imageData.previewElement, error.message);
                hasErrors = true;
            }
        }

        // Update button text based on results
        if (hasErrors) {
            updateUploadButtonText(uploadedCount, totalImages, true);
        } else {
            uploadAllBtn.textContent = `✓ All ${totalImages} images uploaded successfully!`;
            uploadAllBtn.className = 'upload-all-btn success';

            // Clear all uploaded images after a delay
            setTimeout(() => {
                clearAllUploadedImages();
            }, 2000);
        }

    } catch (error) {
        console.error('Upload process failed:', error);
        uploadAllBtn.textContent = 'Upload failed - Try again';
        uploadAllBtn.className = 'upload-all-btn error';
        uploadAllBtn.disabled = false;
    } finally {
        // Reset upload state
        isUploading = false;
        enableAllRemoveButtons();
    }
}

/**
 * Disables all remove buttons to prevent removal during upload.
 */
function disableAllRemoveButtons() {
    const removeButtons = document.querySelectorAll('.preview-remove-btn:not(.uploaded):not(.error)');
    removeButtons.forEach(btn => {
        btn.disabled = true;
        btn.classList.add('uploading');
    });
}

/**
 * Re-enables all remove buttons after upload is complete.
 */
function enableAllRemoveButtons() {
    const removeButtons = document.querySelectorAll('.preview-remove-btn:not(.uploaded):not(.error)');
    removeButtons.forEach(btn => {
        btn.disabled = false;
        btn.classList.remove('uploading');
    });
}

/**
 * Gets the selected rating for an image from its radio buttons.
 * @param {string} key - The key for the image
 * @returns {string} The selected rating
 */
function getSelectedRating(key) {
    const checkedInput = document.querySelector(`input[name="rating-${key}"]:checked`);
    return checkedInput ? checkedInput.value : 'safe';
}

/**
 * Uploads a single image to the server by creating a FormData and submitting it.
 * @param {Object} imageData - The image data object
 * @param {string} rating - The selected rating
 * @returns {Promise<void>}
 */
async function uploadSingleImage(imageData, rating) {
    const formData = new FormData();

    // Add the image blob to the form
    if (imageData.type === 'local') {
        formData.append('image', imageData.blob, imageData.nameOrUrl);
    } else {
        // For link images, create a File object from the blob
        const file = new File([imageData.blob], 'image.jpg', { type: imageData.blob.type });
        formData.append('image', file);
        formData.append('source_url', imageData.nameOrUrl);
    }

    // Add the rating
    formData.append('rating', rating);

    // Submit to the backend
    const response = await fetch('/posts/new', {
        method: 'POST',
        body: formData
    });

    if (!response.ok) {
        const errorText = await response.text().catch(() => 'Unknown error');
        throw new Error(errorText || `Upload failed with status ${response.status}`);
    }

    return response;
}

/**
 * Updates the upload button text to show progress.
 * @param {number} current - Current number of uploaded images
 * @param {number} total - Total number of images to upload
 * @param {boolean} hasErrors - Whether there were errors during upload
 */
function updateUploadButtonText(current, total, hasErrors) {
    if (hasErrors) {
        uploadAllBtn.textContent = `⚠ Uploaded ${current}/${total} (some failed)`;
        uploadAllBtn.className = 'upload-all-btn warning';
        uploadAllBtn.disabled = false;
    } else if (current === total) {
        uploadAllBtn.textContent = `✓ All ${total} images uploaded!`;
        uploadAllBtn.className = 'upload-all-btn success';
    } else {
        uploadAllBtn.textContent = `⏳ Uploading (${current}/${total})`;
        uploadAllBtn.className = 'upload-all-btn uploading';
    }
}

/**
 * Marks an image preview as successfully uploaded.
 * @param {HTMLElement} previewElement - The preview element to mark
 */
function markImageAsUploaded(previewElement) {
    previewElement.classList.add('uploaded');
    const removeBtn = previewElement.querySelector('.preview-remove-btn');
    if (removeBtn) {
        removeBtn.textContent = '✓ Uploaded';
        removeBtn.disabled = true;
        removeBtn.className = 'preview-remove-btn uploaded';
    }
}

/**
 * Marks an image preview as failed to upload.
 * @param {HTMLElement} previewElement - The preview element to mark
 * @param {string} errorMessage - The error message to display
 */
function markImageAsError(previewElement, errorMessage) {
    previewElement.classList.add('upload-error');
    const removeBtn = previewElement.querySelector('.preview-remove-btn');
    if (removeBtn) {
        removeBtn.textContent = '✗ Failed';
        removeBtn.className = 'preview-remove-btn error';
        removeBtn.title = errorMessage;
    }
}

/**
 * Clears all successfully uploaded images from the preview area.
 */
function clearAllUploadedImages() {
    const uploadedElements = document.querySelectorAll('.preview-area.uploaded');
    uploadedElements.forEach(element => {
        const key = getImageKeyFromElement(element);
        if (key) {
            imageBlobMapping.delete(key);
        }
        element.remove();
    });
    updateUploadAllBtn();

    // Reset button state if all images are cleared
    if (imageBlobMapping.size === 0) {
        if (uploadAllBtn) {
            uploadAllBtn.textContent = 'Upload All';
            uploadAllBtn.className = 'upload-all-btn';
            uploadAllBtn.disabled = false;
        }
    }
}

/**
 * Gets the image key from a preview element by looking at its radio button names.
 * @param {HTMLElement} previewElement - The preview element
 * @returns {string|null} The image key or null if not found
 */
function getImageKeyFromElement(previewElement) {
    const radioInput = previewElement.querySelector('input[type="radio"]');
    if (radioInput && radioInput.name) {
        // Extract key from "rating-{key}" format
        const match = radioInput.name.match(/^rating-(.+)$/);
        return match ? match[1] : null;
    }
    return null;
}

setupLocalImageUpload();
