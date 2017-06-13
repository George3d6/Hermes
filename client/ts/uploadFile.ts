//Uploads a file and some metadata based on a pre-defined form
const uploadFile = (form_id: string) => {
    if (window.navigator.userAgent.toLowerCase().indexOf('firefox') > -1) {
        const uploadForm: HTMLFormElement = <HTMLFormElement>document.getElementById(form_id);
        document.getElementById("submit_form").addEventListener("click",(e) => {
            e.preventDefault();

            let reader: FileReader = new FileReader();
            reader.readAsArrayBuffer((<HTMLInputElement>document.getElementById('file')).files[0]);

            reader.onload = function(evt) {

                const formData = new FormData(uploadForm);
                const isPublic: string = String((<HTMLInputElement>document.getElementById('public_switch')).checked);
                formData.append('file', (<any>evt.target).result);
                formData.append('compression', (<HTMLInputElement>document.getElementById('compression')).value);
                formData.append('ispublis', isPublic);
                alert(' the form value is:  ' + formData.get('ispublis'));
                const xhr = new XMLHttpRequest();
                xhr.open("POST", "/upload/");
                xhr.send(formData);
                xhr.onreadystatechange = function() {
                    console.log(xhr.responseText + '  \n status is: ' + xhr.statusText);
                }

            };
        });
    } else {
        console.log("Warnning, your browser does not support asynchronous upload of large files");
    }
}
export default uploadFile
