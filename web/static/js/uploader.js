const folderList = document.getElementById('folder-list');
const addFolderBtn = document.getElementById('add-folder-btn');
const addFolderBtnNode = document.getElementById('add-folder-btn-node');
const fileInput = document.getElementById('file-input');
const folderTemplate = document.getElementById('folder-template').content;
const uploadBtn = document.getElementById('upload-btn');
const logoutBtn = document.getElementById('logout-btn');

let folderCount = 0;
const maxFolders = 5;

addFolderBtn.addEventListener('click', () => {
    if (folderCount < maxFolders) {
        fileInput.click();
    }
});

logoutBtn.addEventListener('click', async () => {
    const response = await fetch('/logout', {
        method: "POST",
    }); 

    if (response.ok) {
        window.location.reload();
    }
})

fileInput.addEventListener('change', handleFolderSelection);

function handleFolderSelection(event) {
    const files = Array.from(event.target.files);

    console.log(files)

    const folderName = files[0].webkitRelativePath.split('/')[0];

    const folderClone = document.importNode(folderTemplate, true);
    const folderNameElem = folderClone.querySelector('.folder-name');
    folderNameElem.textContent = folderName;

    const fileListElem = folderClone.querySelector('.file-list');
    const toggleBtn = folderClone.querySelector('.toggle-files-btn');
    const progressBar = folderClone.querySelector('.progress-bar');
    const trashBtn = folderClone.querySelector('.delete-folder-btn');
    
    toggleBtn.addEventListener('click', () => {
        const isVisible = getComputedStyle(fileListElem).display !== 'none';
        fileListElem.style.display = isVisible ? 'none' : 'block';
        toggleBtn.classList.toggle('open', !isVisible);
    });

    folderList.insertBefore(folderClone, addFolderBtnNode);

    trashBtn.addEventListener('click', async () => {
        const response = await fetch(`/photos?foldername=${folderName}`, {
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
        let formData = new FormData();
        let listItems = [];
        
        let count = 0;  
        let length = folder.files.length;

        for (let i = 0; i < folder.files.length; i++) {
            if (!folder.files[i].type.startsWith('image/')) {
                length--
                continue
            }
            count++;

            formData.append(`file${count}`, folder.files[i]);

            const listItem = document.createElement('li');
            listItem.textContent = folder.files[i].name;
            listItems.push(listItem)
            
            if ((i + 1) % 5 === 0 || i === folder.files.length - 1) {
                const response = await fetch (`/photos?folder=${folder.name}&count=${count}`, {
                    method: 'POST',
                    body: formData,
                });

                if (response.ok) {
                    listItems.forEach((listItem) => {
                        fileListElem.appendChild(listItem);
                        folder.uploadedFiles++;
                        folder.progressBar.value = (folder.uploadedFiles / length) * 100;
                    })
                   
                } else {
                    console.error(`Ошибка загрузки файлов`);
                }

                formData = new FormData();
                listItems = [];
                count = 0;
            }
        }
    })(folder);
}