document.getElementById('file-input').addEventListener('change', async function(event) {
    const files = event.target.files;
    const folderPathElement = document.getElementById('folder-path');
    const fileListElement = document.getElementById('file-list');
    const progressBarElement = document.getElementById('progress-bar');
    let uploadedFiles = 0;

    if (files.length > 0) {
        const folderPath = files[0].webkitRelativePath.split('/')[0];
        folderPathElement.textContent = folderPath;

        for (let i = 0; i < files.length; i++) {
            const file = files[i];
            
            const formData = new FormData();
            formData.append('file', file);
            formData.append('path', file.path);
            formData.append('lastModified', file.lastModified); 

            try {
                const response = await fetch('/upload', {
                    method: 'POST',
                    body: formData
                });
                
                if (response.ok) {
                    const jsonResponse = await response.json();

                    const fileItem = document.createElement('div');
                    fileItem.textContent = `${formData.name} - Uploaded`;
                    fileListElement.appendChild(fileItem);

                    uploadedFiles++;
                    progressBarElement.style.width = `${(uploadedFiles / files.length) * 100}%`;
                } else {
                    console.error(`Failed to upload ${formData.name}`);
                }
            } catch (error) {
                console.error(`Error uploading file: ${formData.name}`, error);
            }
        }
    }
});


