const folderList = document.getElementById('folder-list');
const addFolderBtn = document.getElementById('add-folder-btn');
const addFolderBtnNode = document.getElementById('add-folder-btn-node');
const fileInput = document.getElementById('file-input');
const folderTemplate = document.getElementById('folder-template').content;
const uploadBtn = document.getElementById('upload-btn');

// const vkSendBtn = document.getElementById('vk-send-btn');
// vkSendBtn.disabled = true

let folderCount = 0;
const maxFolders = 5;

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
    const trashBtn = folderClone.querySelector('.delete-folder-btn');
    
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

    folderList.insertBefore(folderClone, addFolderBtnNode);

    trashBtn.addEventListener('click', async () => {
        const response = await fetch(`/uploader/delete?foldername=${folderName}`, {
            method: 'DELETE',
        }); 

        if (response.ok) {
            trashBtn.parentNode.parentNode.remove()
        }
        
        folderCount--;
        if (folderCount < maxFolders) {
            addFolderBtn.disabled = false; 
        }
    });

    

    folderCount++;
    if (folderCount >= maxFolders) {
        addFolderBtn.disabled = true;
    }

    const folder = {
        files: files,
        progressBar: progressBar,
        uploadedFiles: 0,
        name: folderName,
    };

    (async function(folder) {
        for (const file of folder.files) {
            const formData = new FormData();
            formData.append('file', file);
            formData.append('path', file.webkitRelativePath);
            formData.append('lastModified', file.lastModified);
    
            try {
                const response = await fetch('/uploader/upload', {
                    method: 'POST',
                    body: formData
                });
    
                if (response.ok) {
                    folder.uploadedFiles++;
                    folder.progressBar.value = (folder.uploadedFiles / folder.files.length) * 100;
                } else {
                    console.error(`Ошибка загрузки файла ${file.name}`);
                }
            } catch (error) {
                console.error(`Ошибка при загрузке файла ${file.name}:`, error);
            }
        }
    })(folder);
}

// vkSendBtn.addEventListener('click', async () => {
//     folders.forEach(async (folder) => {
//         if (folder.uploadedFiles == folder.files.length) {
//             await sendFolder(folder)
//         }
//     })
// });

// async function sendFolder(folder) {
//     let uploadedFiles = 0;

//     folder.files.forEach(async (file) => {
//         const formData = new FormData();
//         formData.append('path', file.webkitRelativePath);

//         const response = await fetch('/send', {
//             method: 'POST',
//             body: formData
//         });
//         const totalFiles = folder.files.length;

//         if (response.ok) {
//             uploadedFiles++;
//             folder.progressBar.value = (uploadedFiles / totalFiles) * 100;
//         }
//     })
// }