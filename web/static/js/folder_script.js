const folderList = document.getElementById('folder-list');
const addFolderBtn = document.getElementById('add-folder-btn');
const fileInput = document.getElementById('file-input');
const folderTemplate = document.getElementById('folder-template').content;
const uploadBtn = document.getElementById('upload-btn');

let folderCount = 0;
const maxFolders = 5;
let folders = [];

addFolderBtn.addEventListener('click', () => {
    if (folderCount < maxFolders) {
        fileInput.click();
    }
});

fileInput.addEventListener('change', handleFolderSelection);

function handleFolderSelection(event) {
    const files = Array.from(event.target.files);

    const folderName = files[0].webkitRelativePath.split('/')[0];

    const folderClone = document.importNode(folderTemplate, true);
    const folderNameElem = folderClone.querySelector('.folder-name');
    folderNameElem.textContent = folderName;

    const fileListElem = folderClone.querySelector('.file-list');
    const toggleBtn = folderClone.querySelector('.toggle-files-btn');
    const progressBar = folderClone.querySelector('.progress-bar');
    const vkSendBtn = folderClone.querySelector('.vk-send-button');
    vkSendBtn.disabled = true

    files.forEach(file => {
        const listItem = document.createElement('li');
        listItem.textContent = file.name;
        fileListElem.appendChild(listItem);
    });

    toggleBtn.addEventListener('click', () => {
        const isVisible = fileListElem.style.display !== 'none';
        fileListElem.style.display = isVisible ? 'none' : 'block';
        toggleBtn.classList.toggle('open', !isVisible);
    });

    vkSendBtn.addEventListener('click', () => {

    });

    folderList.appendChild(folderClone);

    folderCount++;
    if (folderCount >= maxFolders) {
        addFolderBtn.disabled = true;
    }

    const folder = {
        files: files,
        progressBar: progressBar,
        vkSendBtn: vkSendBtn,
    }

    folders.push(folder)
}

uploadBtn.addEventListener('click', async () => {
    folders.forEach(async (folder) => {
        if (folder.progressBar.value == 0) {
            await uploadFolder(folder)
        }
    })
});

async function uploadFolder(folder) {
    let uploadedFiles = 0;

    folder.files.forEach(async (file) => {
        const formData = new FormData();
        formData.append('file', file);
        formData.append('path', file.webkitRelativePath);
        formData.append('lastModified', file.lastModified); 

        const response = await fetch('/upload', {
            method: 'POST',
            body: formData
        });

        const totalFiles = folder.files.length;

        if (response.ok) {
            uploadedFiles++;
            folder.progressBar.value = (uploadedFiles / totalFiles) * 100;
        }

        folder.vkSendBtn.disabled = false;
    })
}