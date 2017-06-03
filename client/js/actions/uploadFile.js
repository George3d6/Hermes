const uploadFile = (form_id) => {
    let publicFile = "false";
    document.getElementById('public_switch').addEventListener('click', () => {
        if(publicFile === "false") {
            publicFile = "true";
        } else {
            publicFile = "false";
        }
        console.log(publicFile);
    })
    const uploadForm = document.getElementById(form_id);
    uploadForm.onsubmit = () => {
        const formData = new FormData(uploadForm);
        formData.append('file', document.getElementById('file').files[0]);
        formData.append('compression', document.getElementById('compression').value);
        formData.append('public', publicFile);
        const xhr = new XMLHttpRequest();
        xhr.open("POST", "/upload/");
        xhr.send(formData);
        xhr.onreadystatechange = function() {
            console.log(xhr.responseText);
        }
    }
}

export default uploadFile
