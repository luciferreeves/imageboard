class TagEditor {
    constructor() {
        this.debounceTimer = null;
        this.selectedSuggestionIndex = -1;
        this.currentInput = null;
        this.init();
    }

    init() {
        document.addEventListener('DOMContentLoaded', () => {
            this.bindEvents();
        });
    }

    bindEvents() {
        const tagInputs = document.querySelectorAll('.tag-input');
        const tagRemoveBtns = document.querySelectorAll('.tag-remove-btn');

        tagInputs.forEach(input => {
            input.addEventListener('input', (e) => this.handleInput(e));
            input.addEventListener('keydown', (e) => this.handleKeydown(e));
            input.addEventListener('blur', (e) => this.handleBlur(e));
            input.addEventListener('focus', (e) => this.handleFocus(e));
        });

        tagRemoveBtns.forEach(btn => {
            btn.addEventListener('click', (e) => this.removeTag(e));
        });

        // Close suggestions when clicking outside
        document.addEventListener('click', (e) => {
            if (!e.target.closest('.tag-input-wrapper')) {
                this.hideSuggestions();
            }
        });
    }

    handleInput(e) {
        const input = e.target;
        const query = input.value.trim();
        const tagType = input.dataset.type;

        clearTimeout(this.debounceTimer);

        if (query.length < 2) {
            this.hideSuggestions();
            return;
        }

        input.classList.add('loading');

        this.debounceTimer = setTimeout(() => {
            this.searchTags(query, tagType, input);
        }, 300);
    }

    handleKeydown(e) {
        const suggestionsContainer = this.getSuggestionsContainer(e.target);
        const suggestions = suggestionsContainer.querySelectorAll('.tag-suggestion');

        switch (e.key) {
            case 'ArrowDown':
                e.preventDefault();
                this.selectedSuggestionIndex = Math.min(this.selectedSuggestionIndex + 1, suggestions.length - 1);
                this.updateSuggestionSelection(suggestions);
                break;

            case 'ArrowUp':
                e.preventDefault();
                this.selectedSuggestionIndex = Math.max(this.selectedSuggestionIndex - 1, -1);
                this.updateSuggestionSelection(suggestions);
                break;

            case 'Enter':
                e.preventDefault();
                if (this.selectedSuggestionIndex >= 0 && suggestions[this.selectedSuggestionIndex]) {
                    this.selectSuggestion(suggestions[this.selectedSuggestionIndex]);
                } else {
                    this.addTagFromInput(e.target);
                }
                break;

            case 'Escape':
                this.hideSuggestions();
                e.target.blur();
                break;
        }
    }

    handleBlur(e) {
        // Delay hiding suggestions to allow for clicks
        setTimeout(() => {
            this.hideSuggestions();
        }, 200);
    }

    handleFocus(e) {
        this.currentInput = e.target;
        const query = e.target.value.trim();
        if (query.length >= 2) {
            this.searchTags(query, e.target.dataset.type, e.target);
        }
    }

    async searchTags(query, tagType, input) {
        try {
            const postId = this.getPostIdFromUrl();
            const response = await fetch(`/tags/search_for_image.json?q=${encodeURIComponent(query)}&type=${tagType}&image_id=${postId}`);
            const data = await response.json();

            input.classList.remove('loading');
            this.showSuggestions(data, input, query);
        } catch (error) {
            console.error('Search failed:', error);
            input.classList.remove('loading');
            this.hideSuggestions();
        }
    }

    showSuggestions(data, input, query) {
        const suggestionsContainer = this.getSuggestionsContainer(input);
        const { tags, can_create } = data;
        const expectedTagType = input.dataset.type;

        let html = '';

        // Existing tags - only show tags that match the expected type
        tags.forEach(tag => {
            if (tag.type === expectedTagType) {
                html += `
                    <div class="tag-suggestion" data-tag-id="${tag.ID}" data-tag-name="${tag.name}" data-tag-type="${tag.type}">
                        <span class="tag-suggestion-name" style="color: ${this.getTagTypeColor(tag.type)};">${tag.name}</span>
                        <span class="tag-suggestion-count">(${tag.count})</span>
                    </div>
                `;
            }
        });

        // Create new tag option - only show if no exact match found and user can create tags
        const exactMatch = tags.some(tag => tag.name.toLowerCase() === query.toLowerCase() && tag.type === expectedTagType);
        const hasMatchingTypeTags = tags.some(tag => tag.type === expectedTagType);

        if (can_create && !exactMatch && !hasMatchingTypeTags) {
            html += `
                <div class="tag-suggestion tag-suggestion-create" data-create="true" data-tag-name="${query}" data-tag-type="${expectedTagType}">
                    <span class="tag-suggestion-name">Create "${query}" as ${expectedTagType}</span>
                    <span class="tag-suggestion-count">New tag</span>
                </div>
            `;
        }

        suggestionsContainer.innerHTML = html;

        if (html) {
            suggestionsContainer.classList.add('show');
            this.bindSuggestionEvents(suggestionsContainer);
            this.selectedSuggestionIndex = -1;
        } else {
            suggestionsContainer.classList.remove('show');
        }
    }

    bindSuggestionEvents(container) {
        const suggestions = container.querySelectorAll('.tag-suggestion');
        suggestions.forEach((suggestion, index) => {
            suggestion.addEventListener('click', () => this.selectSuggestion(suggestion));
            suggestion.addEventListener('mouseenter', () => {
                this.selectedSuggestionIndex = index;
                this.updateSuggestionSelection(suggestions);
            });
        });
    }

    updateSuggestionSelection(suggestions) {
        suggestions.forEach((suggestion, index) => {
            if (index === this.selectedSuggestionIndex) {
                suggestion.classList.add('selected');
            } else {
                suggestion.classList.remove('selected');
            }
        });
    }

    async selectSuggestion(suggestion) {
        const tagName = suggestion.dataset.tagName;
        const tagType = suggestion.dataset.tagType;
        const isCreate = suggestion.dataset.create === 'true';

        try {
            await this.addTag(tagName, tagType);
            this.hideSuggestions();

            if (this.currentInput) {
                this.currentInput.value = '';
            }
        } catch (error) {
            console.error('Failed to add tag:', error);
            let errorMessage = error.message;

            // Handle specific cross-category errors
            if (errorMessage.includes('already exists as') || errorMessage.includes('previously existed as')) {
                errorMessage = `Tag "${tagName}" ${errorMessage}`;
            }

            this.showError(errorMessage);
        }
    }

    async addTagFromInput(input) {
        const tagName = input.value.trim();
        const tagType = input.dataset.type;

        if (!tagName) return;

        try {
            await this.addTag(tagName, tagType);
            input.value = '';
            this.hideSuggestions();
        } catch (error) {
            console.error('Failed to add tag:', error);
            let errorMessage = error.message;

            // Handle specific cross-category errors
            if (errorMessage.includes('already exists as') || errorMessage.includes('previously existed as')) {
                errorMessage = `Tag "${tagName}" ${errorMessage}`;
            }

            this.showError(errorMessage);
        }
    }

    async addTag(tagName, tagType) {
        // Validate that we're adding to the correct category
        if (!tagType || !['general', 'artist', 'character', 'copyright', 'meta'].includes(tagType)) {
            throw new Error('Invalid tag type');
        }

        const postId = this.getPostIdFromUrl();
        const response = await fetch(`/tags/add_to_image.json?image_id=${postId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                tag_name: tagName,
                tag_type: tagType
            })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to add tag');
        }

        const result = await response.json();
        if (result.success) {
            // Validate that the returned tag matches the requested type
            if (result.tag.type !== tagType) {
                throw new Error(`Tag type mismatch: expected ${tagType}, got ${result.tag.type}`);
            }
            this.addTagToUI(result.tag, tagType);
        }
    } async removeTag(e) {
        e.preventDefault();
        e.stopPropagation();

        const tagId = e.target.dataset.tagId;
        const tagItem = e.target.closest('.tag-item');

        try {
            const postId = this.getPostIdFromUrl();

            const response = await fetch(`/tags/remove_from_image.json?image_id=${postId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    tag_id: parseInt(tagId)
                })
            });

            if (!response.ok) {
                throw new Error('Failed to remove tag');
            }

            const result = await response.json();
            if (result.success) {
                const category = tagItem.closest('.tag-category');
                const tagList = category.querySelector('.tag-list');

                tagItem.style.animation = 'fadeOut 0.3s ease';
                setTimeout(() => {
                    tagItem.remove();
                    this.updateTagCount(category);

                    // Check if no tags remain and add "No tags" message
                    const remainingTags = tagList.querySelectorAll('.tag-item');
                    if (remainingTags.length === 0) {
                        const tagType = category.dataset.type;
                        const noTagsDiv = document.createElement('div');
                        noTagsDiv.className = 'no-tags';
                        noTagsDiv.textContent = `No ${tagType} tags`;
                        tagList.appendChild(noTagsDiv);
                    }
                }, 300);
            }
        } catch (error) {
            console.error('Failed to remove tag:', error);
            this.showError('Failed to remove tag');
        }
    } addTagToUI(tag, expectedTagType) {
        // Validate that the tag type matches what we expect
        if (tag.type !== expectedTagType) {
            console.error(`Tag type mismatch: expected ${expectedTagType}, got ${tag.type}`);
            this.showError(`Cannot add ${tag.type} tag "${tag.name}" to ${expectedTagType} section`);
            return;
        }

        const tagList = document.getElementById(`tag-list-${expectedTagType}`);

        // Check if tag already exists in this category
        const existingTag = tagList.querySelector(`[data-tag-id="${tag.ID}"]`);
        if (existingTag) {
            return; // Tag already exists, don't add duplicate
        }

        // Also check if tag exists in any other category (cross-category validation)
        const allTagLists = document.querySelectorAll('.tag-list');
        for (let otherList of allTagLists) {
            if (otherList.id !== `tag-list-${expectedTagType}`) {
                const existingInOther = otherList.querySelector(`[data-tag-id="${tag.ID}"]`);
                if (existingInOther) {
                    console.warn(`Tag "${tag.name}" already exists in another category`);
                    return;
                }
            }
        }

        const noTagsElement = tagList.querySelector('.no-tags');

        if (noTagsElement) {
            noTagsElement.remove();
        }

        const tagItem = document.createElement('div');
        tagItem.className = 'tag-item';
        tagItem.dataset.tagId = tag.ID;
        tagItem.style.animation = 'fadeIn 0.3s ease';

        tagItem.innerHTML = `
            <a href="/tags/${tag.name}" class="tag-link" style="color: ${this.getTagTypeColor(tag.type)};">
                ${tag.name}
            </a>
            ${tag.parent ? `<span class="tag-parent-indicator" title="Child of ${tag.parent.name}">⬆</span>` : ''}
            ${tag.children && tag.children.length > 0 ? `<span class="tag-children-indicator" title="Has ${tag.children.length} children">⬇</span>` : ''}
            <button type="button" class="tag-remove-btn" data-tag-id="${tag.ID}" title="Remove tag">×</button>
        `;

        tagList.appendChild(tagItem);

        // Bind remove event
        tagItem.querySelector('.tag-remove-btn').addEventListener('click', (e) => this.removeTag(e));

        // Update count
        this.updateTagCount(tagList.closest('.tag-category'));
    }

    updateTagCount(category) {
        const tagList = category.querySelector('.tag-list');
        const countElement = category.querySelector('.tag-count');
        const tagItems = tagList.querySelectorAll('.tag-item');

        countElement.textContent = `(${tagItems.length})`;
    }

    getSuggestionsContainer(input) {
        return input.parentElement.querySelector('.tag-suggestions');
    }

    hideSuggestions() {
        const allSuggestions = document.querySelectorAll('.tag-suggestions');
        allSuggestions.forEach(container => {
            container.classList.remove('show');
        });
        this.selectedSuggestionIndex = -1;
    }

    getPostIdFromUrl() {
        const match = window.location.pathname.match(/\/posts\/(\d+)/);
        return match ? match[1] : null;
    }

    getTagTypeColor(tagType) {
        const colors = {
            general: '#4ECDC4',
            artist: '#FF6B9D',
            character: '#FFB347',
            copyright: '#A8E6CF',
            meta: '#DDA0DD'
        };
        return colors[tagType] || '#E6E6FA';
    }

    showError(message) {
        // Find the tag editor container and add error message
        const tagEditor = document.querySelector('.tag-editor');

        // Remove any existing error
        const existingError = tagEditor.querySelector('.error');
        if (existingError) {
            existingError.remove();
        }

        // Create standard error div
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error';
        errorDiv.textContent = message;

        // Insert at the top of tag editor
        tagEditor.insertBefore(errorDiv, tagEditor.firstChild);

        // Remove after 3 seconds
        setTimeout(() => {
            errorDiv.remove();
        }, 3000);
    }
}

// CSS animations for tag items
const style = document.createElement('style');
style.textContent = `
    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(-10px); }
        to { opacity: 1; transform: translateY(0); }
    }
    
    @keyframes fadeOut {
        from { opacity: 1; transform: translateY(0); }
        to { opacity: 0; transform: translateY(-10px); }
    }
`;
document.head.appendChild(style);

// Initialize the tag editor
new TagEditor();
