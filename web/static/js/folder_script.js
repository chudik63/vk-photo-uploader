const folderList = document.getElementById('folder-list');
const addFolderBtn = document.getElementById('add-folder-btn');
const fileInput = document.getElementById('file-input');
const folderTemplate = document.getElementById('folder-template').content;
const uploadBtn = document.getElementById('upload-btn');

const vkSendBtn = document.getElementById('vk-send-btn');
vkSendBtn.disabled = true

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
    
    files.forEach(file => {
        const listItem = document.createElement('li');
        listItem.textContent = file.name;
        fileListElem.appendChild(listItem);
    });

    toggleBtn.addEventListener('click', () => {
        const isVisible = getComputedStyle(fileListElem).display !== 'none';
        fileListElem.style.display = isVisible ? 'none' : 'block';
        toggleBtn.classList.toggle('open', !isVisible);
    });

    folderList.appendChild(folderClone);

    folderCount++;
    if (folderCount >= maxFolders) {
        addFolderBtn.disabled = true;
    }

    const folder = {
        files: files,
        progressBar: progressBar,
        uploadedFiles: 0,
        name: folderName,
    }

    folders.push(folder);
    
    (async () => {
        await Promise.all(folder.files.map(async (file) => {
            const formData = new FormData();
            formData.append('file', file);
            formData.append('path', file.webkitRelativePath);
            formData.append('lastModified', file.lastModified); 
    
            const response = await fetch('/upload', {
                method: 'POST',
                body: formData
            });
    
            if (response.ok) {
                folder.uploadedFiles++;
            }
        }));

        if (folder.uploadedFiles == folder.files.length) {
            alert(`Папка ${folder.name} загружена`);
            vkSendBtn.disabled = false;
        } else {
            alert(`Ошибка загрузки папки ${folder.name}`);
        }
        
    })();

}

vkSendBtn.addEventListener('click', async () => {
    folders.forEach(async (folder) => {
        if (folder.uploadedFiles == folder.files.length) {
            await sendFolder(folder)
        }
    })
});

async function sendFolder(folder) {
    let uploadedFiles = 0;

    folder.files.forEach(async (file) => {
        const formData = new FormData();
        formData.append('path', file.webkitRelativePath);

        const response = await fetch('/send', {
            method: 'POST',
            body: formData
        });
        const totalFiles = folder.files.length;

        if (response.ok) {
            uploadedFiles++;
            folder.progressBar.value = (uploadedFiles / totalFiles) * 100;
        }
    })
}