const uploadFile = (form_id) => {
    const uploadForm = document.getElementById(form_id);
    uploadForm.onsubmit = () => {
        const formData = new FormData(uploadForm);
        formData.append('file', document.getElementById('file').files[0]);
        formData.append('compression', document.getElementById('compression').value);
        const xhr = new XMLHttpRequest();
        xhr.open("POST", "/upload/");
        xhr.send(formData);
        xhr.onreadystatechange = function() {
            console.log(xhr.responseText);
        }
        /*
        //Not using fetch for now because of various errors
        fetch('/upload/', {
            method: 'POST',
            headers: {
                'Content-Type': 'multipart/form-data'
            },
            body: formData
        }).then((response) => {
            console.log(response);
        }).catch((err) => {
            console.log(err);
        })
        */
    }
}

export default uploadFile
