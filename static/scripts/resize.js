function getAvailableHeight() {
    const navbar = document.querySelector('nav')
    const postBar = document.querySelector('.single-post-bar')
    const navbarHeight = navbar ? navbar.offsetHeight : 0
    const postBarHeight = postBar ? postBar.offsetHeight : 0
    const padding = 40
    return window.innerHeight - navbarHeight - postBarHeight - padding
}

function updateSizeSelection(selectedSize) {
    document.querySelectorAll('.single-post-bar-size-actions a').forEach((link) => {
        link.classList.remove('size-selected')
        if (link.onclick && link.onclick.toString().includes(selectedSize)) {
            link.classList.add('size-selected')
        }
    })
}

let currentSize = 'fitBoth'

function handleResize() {
    if (currentSize === 'fitHeight' || currentSize === 'fitBoth') {
        switchSize(currentSize)
    }
}

function switchSize(size) {
    currentSize = size
    const img = document.getElementById('post-image')
    const sizeData = sizes[size]

    img.className = ''
    img.style.width = ''
    img.style.height = ''
    img.style.maxWidth = ''
    img.style.maxHeight = ''

    if (size === 'fitHeight') {
        const availableHeight = getAvailableHeight()
        img.style.maxHeight = availableHeight + 'px'
        img.style.width = 'auto'
        img.style.maxWidth = 'none'
        img.classList.add('fit-height')
    } else if (size === 'fitWidth') {
        img.style.width = '100%'
        img.style.height = 'auto'
        img.style.maxWidth = '100%'
        img.classList.add('fit-width')
    } else if (size === 'fitBoth') {
        const availableHeight = getAvailableHeight()
        const imageWidth = parseInt(sizeData.width || img.naturalWidth)
        const imageHeight = parseInt(sizeData.height || img.naturalHeight)

        if (imageHeight > imageWidth) {
            img.style.maxHeight = availableHeight + 'px'
            img.style.width = 'auto'
            img.style.maxWidth = 'none'
        } else {
            img.style.width = '100%'
            img.style.height = 'auto'
            img.style.maxWidth = '100%'
        }
        img.classList.add('fit-both')
    } else {
        img.style.maxWidth = 'none'
        if (sizeData.width) img.style.width = sizeData.width + 'px'
        if (sizeData.height) img.style.height = sizeData.height + 'px'
        img.classList.add('fixed-size')
    }

    if (img.src !== sizeData.src) {
        img.src = sizeData.src
    }
    updateSizeSelection(size)
}

window.addEventListener('resize', handleResize)