document.getElementById('tokenForm').addEventListener('submit', function(event) {
    event.preventDefault();

    const vkToken = document.getElementById('vkToken').value;
    const folderPath = document.getElementById('folderPath').files[0].webkitRelativePath;

    const formData = new FormData();
    formData.append('vkToken', vkToken);
    formData.append('folderPath', folderPath);

    fetch('http://localhost:8080/upload', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        alert('Token and Folder Path sent successfully!');
        console.log(data);
    })
    .catch(error => console.error('Error:', error));
});